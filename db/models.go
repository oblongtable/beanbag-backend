// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql"
	"time"
)

type Answer struct {
	AnsID       int32
	QuesID      sql.NullInt32
	Description string
	IsCorrect   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Question struct {
	QuesID      int32
	QuizID      sql.NullInt32
	Description string
	TimerOption bool
	Timer       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Quiz struct {
	QuizID      int32
	CreatorID   sql.NullInt32
	QuizTitle   string
	Description sql.NullString
	IsPriv      bool
	Timer       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	UserID    int32
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}
