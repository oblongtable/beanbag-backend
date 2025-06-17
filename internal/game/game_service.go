// internal/game/game_service.go
package game

import (
	"errors" 
	"fmt"
	"sync"

	"github.com/oblongtable/beanbag-backend/internal/quiz"
)

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
func (s *GameService) CreateGame(roomID, presenterID, hostID string, quizFilePath string, broadcastFunc BroadcastFunc) (*Game, error) {
	q, err := quiz.LoadQuizFromFile(quizFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load quiz: %w", err)
	}

	game := &Game{
		ID:              roomID,
		PresenterID:     presenterID,
		HostID:          hostID,
		State:           StateLobby,
		quiz:            q,
		players:         make(map[string]*Player),
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
