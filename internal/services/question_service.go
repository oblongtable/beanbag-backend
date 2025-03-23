package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/oblongtable/beanbag-backend/db"
)

type QuestionService struct {
	queries *db.Queries
}

func NewQuestionService(queries *db.Queries) *QuestionService {
	return &QuestionService{queries: queries}
}

func (s *QuestionService) CreateQuestion(ctx context.Context, quizID int32, description string, timerOption bool, timer int32) (*db.Question, error) {
	params := db.CreateQuestionParams{
		QuizID:      sql.NullInt32{Int32: quizID, Valid: true},
		Description: description,
		TimerOption: timerOption,
		Timer:       timer,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	question, err := s.queries.CreateQuestion(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating question: %w", err)
	}

	return &question, nil
}

func (s *QuestionService) GetQuestion(ctx context.Context, questionID int32) (*db.Question, error) {
	question, err := s.queries.GetQuestion(ctx, questionID)
	if err != nil {
		return nil, fmt.Errorf("error getting question: %w", err)
	}
	return &question, nil
}
