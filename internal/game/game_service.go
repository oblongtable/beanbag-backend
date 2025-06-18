// internal/game/game_service.go
package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/oblongtable/beanbag-backend/internal/quiz"
)

const hardcodedQuiz = `{
  "title": "General Knowledge Challenge",
  "description": "A fun quiz to test your knowledge on a variety of topics!",
  "sections": [
    {
      "section": "General Knowledge",
      "type": "multiple-choice",
      "questions": [
        {
          "questionText": "What is the capital city of Australia?",
          "options": [
            "Sydney",
            "Melbourne",
            "Canberra",
            "Perth"
          ],
          "correctOptionIndex": 2,
          "timeLimit": 20,
          "points": 100,
          "explanation": "While Sydney is the largest city, Canberra has been the capital of Australia since 1927."
        },
        {
          "questionText": "Which planet is known as the 'Red Planet'?",
          "options": [
            "Venus",
            "Mars",
            "Jupiter",
            "Saturn"
          ],
          "correctOptionIndex": 1,
          "timeLimit": 20,
          "points": 100,
          "explanation": "Mars is often called the 'Red Planet' because of the iron oxide prevalent on its surface, which gives it a reddish appearance."
        },
        {
          "questionText": "What is the largest ocean on Earth?",
          "options": [
            "Atlantic Ocean",
            "Indian Ocean",
            "Arctic Ocean",
            "Pacific Ocean"
          ],
          "correctOptionIndex": 3,
          "timeLimit": 20,
          "points": 100,
          "explanation": "The Pacific Ocean is the largest and deepest of the world's five oceans, covering more than 30% of the Earth's surface."
        },
        {
          "questionText": "Who wrote the famous play 'Romeo and Juliet'?",
          "options": [
            "William Shakespeare",
            "Charles Dickens",
            "Jane Austen",
            "Mark Twain"
          ],
          "correctOptionIndex": 0,
          "timeLimit": 25,
          "points": 120,
          "explanation": "William Shakespeare, the famous English playwright, wrote the tragic romance of 'Romeo and Juliet'."
        },
        {
          "questionText": "What is the chemical symbol for Gold?",
          "options": [
            "Ag",
            "Go",
            "Au",
            "Gd"
          ],
          "correctOptionIndex": 2,
          "timeLimit": 15,
          "points": 100,
          "explanation": "The symbol 'Au' for gold comes from its Latin name, 'aurum'."
        },
        {
          "questionText": "In which country would you find the ancient city of Petra?",
          "options": [
            "Egypt",
            "Greece",
            "Peru",
            "Jordan"
          ],
          "correctOptionIndex": 3,
          "timeLimit": 30,
          "points": 150,
          "explanation": "The famous archaeological site of Petra, known for its rock-cut architecture, is located in southern Jordan."
        }
      ]
    }
  ]
}`

type GameService struct {
	games map[string]*Game
	mu    sync.RWMutex
}

func NewService() *GameService {
	return &GameService{
		games: make(map[string]*Game),
	}
}

// CreateGame loads a quiz, creates a Game, and links it to the room.
// It now also initializes the players map from the room participants.
func (s *GameService) CreateGame(roomID, presenterID, hostID string, broadcastFunc BroadcastFunc, initialPlayers []InitialPlayerInfo) (*Game, error) {
	var q quiz.Quiz
	err := json.NewDecoder(bytes.NewReader([]byte(hardcodedQuiz))).Decode(&q)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hardcoded quiz: %w", err)
	}

	playersMap := make(map[string]*Player)
	for _, pInfo := range initialPlayers {
		playersMap[pInfo.ID] = &Player{
			ID:    pInfo.ID,
			Name:  pInfo.Username,
			Score: 0, // Initialize score to 0
		}
	}

	game := &Game{
		ID:              roomID,
		PresenterID:     presenterID,
		HostID:          hostID,
		State:           StateLobby,
		quiz:            &q,
		players:         playersMap, // Use the populated players map
		questionAnswers: make(map[string]PlayerAnswer),
		broadcastFunc:   broadcastFunc, // Pass the broadcast function
	}

	s.mu.Lock()
	s.games[roomID] = game
	s.mu.Unlock()

	return game, nil
}

// Start will now just begin the title screen, not the whole loop.
func (s *GameService) StartGame(gameID string, hostID string) error {
	game, found := s.GetGame(gameID)
	if !found {
		return errors.New("game not found")
	}
	if game.HostID != hostID {
		return errors.New("only the host can start the game")
	}

	game.startTitleScreen()
	return nil
}

// NextAction is the new method called by the host to advance the game.
func (s *GameService) NextAction(gameID string, hostID string) error {
	game, found := s.GetGame(gameID)
	if !found {
		return errors.New("game not found")
	}
	if game.HostID != hostID {
		return errors.New("only the host can advance the game")
	}

	return game.nextState() // Delegate the action to the game instance
}

// HandleAnswer now needs more complex logic.
func (s *GameService) HandleAnswer(gameID, playerID string, answerIndex int) error {
	game, found := s.GetGame(gameID)
	if !found {
		return errors.New("game not found")
	}
	return game.handlePlayerAnswer(playerID, answerIndex)
}

// GetGame retrieves a game instance (needed by the ws_handler).
func (s *GameService) GetGame(gameID string) (*Game, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	game, found := s.games[gameID]
	return game, found
}
