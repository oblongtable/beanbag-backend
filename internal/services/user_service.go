package services

import (
	"context"
	"fmt"
	"time"
	"errors"
	"database/sql"

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

func (s *UserService) GetUserById(ctx context.Context, userID int32) (*db.User, error) {
	user, err := s.queries.GetUserById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	return &user, nil
}

// SyncUser finds a user by Email, updates name if found, creates if not.
// Returns the user, a boolean indicating if created, and an error.
func (s *UserService) SyncUser(ctx context.Context, name string, email string) (*db.User, bool, error) {
	// Try to find the user by email
	existingUser, err := s.queries.GetUserByEmail(ctx, email)

	if err == nil {
		// User found by Email
		fmt.Printf("User found with email %s, checking name.\n", email)
		// Check if the name needs updating
		if existingUser.Name != name {
			fmt.Printf("Updating name for user %d from '%s' to '%s'.\n", existingUser.UserID, existingUser.Name, name)
			// Use the existing user's primary key (UserID) to update
			updatedUser, updateErr := s.queries.UpdateUser(ctx, db.UpdateUserParams{
				UserID: existingUser.UserID,
				Name:   sql.NullString{String: name, Valid: true}, // Update name
			})
			if updateErr != nil {
				return nil, false, fmt.Errorf("failed to update existing user %s: %w", email, updateErr)
			}
			return &updatedUser, false, nil // Return updated user, created = false
		}
		// Name is the same, no update needed
		fmt.Printf("Name for user %s is already up-to-date.\n", email)
		return &existingUser, false, nil // Return existing user, created = false

	} else if !errors.Is(err, sql.ErrNoRows) {
		// An actual error occurred during lookup (not just 'not found')
		return nil, false, fmt.Errorf("error checking user by email %s: %w", email, err)
	}

	// User not found by email (err == sql.ErrNoRows), create them
	fmt.Printf("User with email %s not found, creating new user.\n", email)
	newUser, createErr := s.queries.CreateUserMinimal(ctx, db.CreateUserMinimalParams{
		Name:      name,
		Email:     email,
	})
	if createErr != nil {
		// Handle potential constraint errors (e.g., duplicate email if UNIQUE constraint exists and lookup raced)
		return nil, false, fmt.Errorf("failed to create new user for %s: %w", email, createErr)
	}

	return &newUser, true, nil // Return newly created user, created = true
}