package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/oblongtable/beanbag-backend/db"
)

type QuizService struct {
	queries *db.Queries
}

func NewQuizService(queries *db.Queries) *QuizService {
	return &QuizService{queries: queries}
}

func (s *QuizService) CreateQuiz(ctx context.Context, title string, creatorID int32) (*db.Quiz, error) {
	params := db.CreateQuizParams{
		QuizTitle:   title,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatorID:   sql.NullInt32{Int32: creatorID, Valid: true},
		IsPriv:      false,
		Timer:       0,
	}

	quiz, err := s.queries.CreateQuiz(ctx, params)
	if err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (s *QuizService) GetQuiz(ctx context.Context, quizID int32) (*db.Quiz, error) {
	quiz, err := s.queries.GetQuiz(ctx, quizID)
	if err != nil {
		return nil, fmt.Errorf("error getting quiz: %w", err)
	}
	return &quiz, nil
}
