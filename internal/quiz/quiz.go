// internal/quiz/quiz.go
package quiz

import (
	"encoding/json"
	"os"
)

// Question represents a single question in the quiz, matching the JSON structure.
type Question struct {
	ID                 int      `json:"id"`
	QuestionText       string   `json:"questionText"`
	Options            []string `json:"options"`
	CorrectOptionIndex int      `json:"correctOptionIndex"`
	TimeLimit          int      `json:"timeLimit"` // Time in seconds
	Points             int      `json:"points"`
	Explanation        string   `json:"explanation"`
}

// Quiz represents the entire quiz structure.
type Quiz struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

// LoadQuizFromFile reads a quiz from a JSON file on disk.
// This function does not need to change.
func LoadQuizFromFile(filePath string) (*Quiz, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var quiz Quiz
	if err := json.Unmarshal(data, &quiz); err != nil {
		return nil, err
	}
	return &quiz, nil
}
