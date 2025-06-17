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

// BroadcastFunc is a function type for broadcasting messages to clients.
type BroadcastFunc func(msgType string, payload interface{})

// Game represents a single game instance with its state.
type Game struct {
	ID          string // This is the Game PIN
	PresenterID string
	HostID      string
	State       GameState

	quiz            *quiz.Quiz
	players         map[string]*Player
	currentQuestion int
	questionAnswers map[string]PlayerAnswer // Map[playerID]PlayerAnswer for the current question

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
		g.broadcastMessage("show_section", map[string]string{"title": "Get Ready!"})
		log.Printf("Game %s: Showing section. Waiting for host.", g.ID)

	case StateSection, StateScores:
		// If we are on a break screen or have just finished showing scores, start the next question.
		if g.currentQuestion >= len(g.quiz.Questions) {
			// No more questions, end the game.
			g.finishGame()
		} else {
			// Start the next question
			g.startQuestion()
		}

	default:
		return fmt.Errorf("cannot advance state from current state: %s", g.State)
	}
	return nil
}

// startQuestion prepares and broadcasts the current question and starts its timer.
func (g *Game) startQuestion() {
	q := g.quiz.Questions[g.currentQuestion]

	g.State = StateQuestion
	g.questionAnswers = make(map[string]PlayerAnswer) // Use the new struct
	g.questionStartTime = time.Now()

	log.Printf("Game %s: Starting question %d.", g.ID, g.currentQuestion+1)
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

	log.Printf("Game %s: Finishing question %d.", g.ID, g.currentQuestion+1)
	g.State = StateScores
	q := g.quiz.Questions[g.currentQuestion]

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
	g.currentQuestion++ // Move to the next question index for the next round
}

// handlePlayerAnswer is called from the service when a player submits an answer.
func (g *Game) handlePlayerAnswer(playerID string, answerIndex int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Ignore answers if not in the question phase or if player has already answered.
	if g.State != StateQuestion {
		return errors.New("not accepting answers right now")
	}
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
		"questionNumber": g.currentQuestion + 1,
		"totalQuestions": len(g.quiz.Questions),
	}
	g.broadcastMessage("new_question", payload)
}

// broadcastScores sends the results of the question and the current leaderboard.
func (g *Game) broadcastScores(q quiz.Question) {
	g.mu.RLock() // Use a read-lock as we are only reading the player scores
	defer g.mu.RUnlock()

	payload := map[string]interface{}{
		"correctOptionIndex": q.CorrectOptionIndex,
		"explanation":        q.Explanation,
		"leaderboard":        g.players, // Send the updated map of players
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
