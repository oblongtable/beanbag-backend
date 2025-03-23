package services

import (
	"context"
	"fmt"
	"time"

	"github.com/oblongtable/beanbag-backend/db"
)

type UserService struct {
	queries *db.Queries
}

func NewUserService(queries *db.Queries) *UserService {
	return &UserService{queries: queries}
}

func (s *UserService) CreateUser(ctx context.Context, name string, email string) (*db.User, error) {
	params := db.CreateUserParams{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &user, nil
}

func (s *UserService) GetUser(ctx context.Context, userID int32) (*db.User, error) {
	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}
