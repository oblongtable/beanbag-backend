package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/oblongtable/beanbag-backend/db"
)

type AnswerService struct {
	queries *db.Queries
}

func NewAnswerService(queries *db.Queries) *AnswerService {
	return &AnswerService{queries: queries}
}

func (s *AnswerService) CreateAnswer(ctx context.Context, questionID int32, description string, isCorrect bool) (*db.Answer, error) {
	params := db.CreateAnswerParams{
		QuesID:      sql.NullInt32{Int32: questionID, Valid: true},
		Description: description,
		IsCorrect:   isCorrect,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	answer, err := s.queries.CreateAnswer(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating answer: %w", err)
	}

	return &answer, nil
}

func (s *AnswerService) GetAnswer(ctx context.Context, answerID int32) (*db.Answer, error) {
	answer, err := s.queries.GetAnswer(ctx, answerID)
	if err != nil {
		return nil, fmt.Errorf("error getting answer: %w", err)
	}
	return &answer, nil
}
