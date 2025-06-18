package game

import (
	"errors"
	"fmt"
	"log" // For logging errors
	"sync"
	"time"

	"github.com/oblongtable/beanbag-backend/internal/quiz"
)

type GameState string

const (
	StateLobby    GameState = "lobby"    // Waiting for players
	StateTitle    GameState = "title"    // Showing the initial quiz title screen
	StateSection  GameState = "section"  // A placeholder state for showing a section title/break
	StateQuestion GameState = "question" // Actively accepting answers
	StateScores   GameState = "scores"   // Showing results of a question
	StateFinished GameState = "finished" // Game over
)

type Player struct {
	ID    string
	Name  string
	Score int
}

type PlayerAnswer struct {
	AnswerIndex int
	TimeTaken   time.Duration
}

type InitialPlayerInfo struct {
	ID       string
	Username string
}

// BroadcastFunc is a function type for broadcasting messages to clients.
type BroadcastFunc func(msgType string, payload interface{})

// Game represents a single game instance with its state.
type Game struct {
	ID          string // This is the Game PIN
	PresenterID string
	HostID      string
	State       GameState

	quiz                     *quiz.Quiz
	players                  map[string]*Player
	currentSection           int                     // New: Index of the current section
	currentQuestionInSection int                     // New: Index of the current question within the current section
	questionAnswers          map[string]PlayerAnswer // Map[playerID]PlayerAnswer for the current question

	// --- New fields for advanced logic ---
	questionStartTime time.Time   // Records when the current question was sent, for scoring
	questionTimer     *time.Timer // The timer for the current question

	// To protect concurrent access to players, state, etc.
	mu sync.RWMutex

	// Function to broadcast messages to clients in the associated room
	broadcastFunc BroadcastFunc
}

func (g *Game) startTitleScreen() {
	g.mu.Lock()
	g.State = StateTitle
	g.mu.Unlock()

	g.broadcastMessage("show_title", map[string]string{
		"title":       g.quiz.Title,
		"description": g.quiz.Description,
	})
	log.Printf("Game %s: Showing title screen. Waiting for host.", g.ID)
}

// nextState is the core state machine driver.
func (g *Game) nextState() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch g.State {
	case StateTitle:
		// Logic to show a "section" screen before the questions begin.
		g.State = StateSection
		g.broadcastMessage("show_section", map[string]interface{}{
			"id":    1, // Assuming the first section has ID 1
			"title": g.quiz.Sections[0].Section,
		})
		log.Printf("Game %s: Showing section. Waiting for host.", g.ID)

	case StateSection, StateScores:
		// Check if there are more questions in the current section
		if g.currentQuestionInSection < len(g.quiz.Sections[g.currentSection].Questions) {
			g.startQuestion()
		} else {
			// Move to the next section
			g.currentSection++
			g.currentQuestionInSection = 0 // Reset question index for the new section

			// Check if there are more sections
			if g.currentSection < len(g.quiz.Sections) {
				// Show the next section title screen
				g.State = StateSection
				g.broadcastMessage("show_section", map[string]interface{}{ // Changed value type to interface{}
					"id":    g.currentSection + 1,
					"title": g.quiz.Sections[g.currentSection].Section,
				})
				log.Printf("Game %s: Showing section '%s'. Waiting for host.", g.ID, g.quiz.Sections[g.currentSection].Section)
			} else {
				// No more sections, end the game.
				g.finishGame()
			}
		}

	default:
		return fmt.Errorf("cannot advance state from current state: %s", g.State)
	}
	return nil
}

// startQuestion prepares and broadcasts the current question and starts its timer.
func (g *Game) startQuestion() {
	currentSection := &g.quiz.Sections[g.currentSection]
	q := currentSection.Questions[g.currentQuestionInSection]

	g.State = StateQuestion
	g.questionAnswers = make(map[string]PlayerAnswer) // Use the new struct
	g.questionStartTime = time.Now()

	log.Printf("Game %s: Starting question %d in section %d.", g.ID, g.currentQuestionInSection+1, g.currentSection+1)
	g.broadcastQuestion(q)

	// This timer will automatically call finishQuestion when the time is up.
	g.questionTimer = time.AfterFunc(time.Duration(q.TimeLimit)*time.Second, g.finishQuestion)
}

// finishQuestion is called when the timer runs out OR all players have answered.
func (g *Game) finishQuestion() {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Prevent this from running twice (e.g., if timer fires right after last player answers)
	if g.State != StateQuestion {
		return
	}

	log.Printf("Game %s: Finishing question %d in section %d.", g.ID, g.currentQuestionInSection+1, g.currentSection+1)
	g.State = StateScores
	currentSection := &g.quiz.Sections[g.currentSection]
	q := currentSection.Questions[g.currentQuestionInSection]

	// --- New Scoring Logic ---
	for playerID, answer := range g.questionAnswers {
		if answer.AnswerIndex == q.CorrectOptionIndex {
			if player, ok := g.players[playerID]; ok {
				// Base points + time bonus
				timeBonus := float64(q.Points) * 0.5 * (1 - (answer.TimeTaken.Seconds() / float64(q.TimeLimit)))
				player.Score += q.Points + int(timeBonus)
			}
		}
	}

	g.broadcastScores(q)
	g.currentQuestionInSection++ // Move to the next question index for the next round
}

// handlePlayerAnswer is called from the service when a player submits an answer.
func (g *Game) handlePlayerAnswer(playerID string, answerIndex int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Ignore answers if not in the question phase or if player has already answered.
	if g.State != StateQuestion {
		log.Printf("Game %s: Player %s tried to answer %d, but game state is %s (not StateQuestion).", g.ID, playerID, answerIndex, g.State)
		return errors.New("not accepting answers right now")
	}
	currentQuestion := g.quiz.Sections[g.currentSection].Questions[g.currentQuestionInSection]
	log.Printf("Game %s: Player %s submitted answer %d. Correct answer is %d. Current question: '%s'", g.ID, playerID, answerIndex, currentQuestion.CorrectOptionIndex, currentQuestion.QuestionText)
	if _, alreadyAnswered := g.questionAnswers[playerID]; alreadyAnswered {
		return errors.New("player has already answered")
	}

	// Record the answer and the time it took.
	g.questionAnswers[playerID] = PlayerAnswer{
		AnswerIndex: answerIndex,
		TimeTaken:   time.Since(g.questionStartTime),
	}

	// Check if all players have answered.
	if len(g.questionAnswers) == len(g.players) {
		// All players have answered, stop the timer and finish the question immediately.
		if g.questionTimer.Stop() {
			go g.finishQuestion()
		}
	}
	return nil
}

// finishGame replaces broadcastFinalResults and is called when the game ends.
func (g *Game) finishGame() {
	g.State = StateFinished
	// ... (The logic from your broadcastFinalResults with sorting) ...
	// g.broadcastMessage("game_over", payload)
	// log.Printf("Game %s finished.", g.ID)
}

// --- Updated Helper Methods for Broadcasting ---

// broadcastQuestion sends the question to all players, hiding the answer.
func (g *Game) broadcastQuestion(q quiz.Question) {
	// We create a new struct for the payload to control what data is sent.
	// We don't want to send the correctOptionIndex or explanation yet.
	payload := map[string]interface{}{
		"questionText":   q.QuestionText,
		"options":        q.Options,
		"timeLimit":      q.TimeLimit,
		"points":         q.Points,
		"questionNumber": g.currentQuestionInSection + 1,
		"totalQuestions": len(g.quiz.Sections[g.currentSection].Questions),
	}
	g.broadcastMessage("new_question", payload)
}

// broadcastScores sends the results of the question and the current leaderboard.
func (g *Game) broadcastScores(q quiz.Question) {
	// Define a struct that matches the frontend's LeaderboardEntry for this specific payload
	type QuestionLeaderboardEntry struct {
		ID    string `json:"ID"`
		Name  string `json:"Name"`
		Score int    `json:"Score"` // Score for this specific question
	}

	questionLeaderboard := make(map[string]QuestionLeaderboardEntry)

	// Iterate over all players in the game
	for playerID, player := range g.players {
		entry := QuestionLeaderboardEntry{
			ID:    player.ID,
			Name:  player.Name,
			Score: 0, // Default to 0 points for this question
		}

		// Check if the player submitted an answer for the current question
		if playerAnswer, ok := g.questionAnswers[playerID]; ok {
			if playerAnswer.AnswerIndex == q.CorrectOptionIndex {
				// Calculate points for this question
				timeBonus := float64(q.Points) * 0.5 * (1 - (playerAnswer.TimeTaken.Seconds() / float64(q.TimeLimit)))
				entry.Score = q.Points + int(timeBonus)
			}
		}
		questionLeaderboard[playerID] = entry
	}

	payload := map[string]interface{}{
		"correctOptionIndex": q.CorrectOptionIndex,
		"explanation":        q.Explanation,
		"leaderboard":        questionLeaderboard, // Send the map of player results for this question
	}
	g.broadcastMessage("question_result", payload)
}

// A generic helper to create and broadcast messages
func (g *Game) broadcastMessage(msgType string, payload interface{}) {
	if g.broadcastFunc != nil {
		g.broadcastFunc(msgType, payload)
	} else {
		log.Printf("Game %s: BroadcastFunc is not set. Cannot broadcast message of type %s.", g.ID, msgType)
	}
}
