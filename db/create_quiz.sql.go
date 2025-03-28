// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: create_quiz.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createQuiz = `-- name: CreateQuiz :one
INSERT INTO quizzes (
    creator_id,
    quiz_title,
    description,
    is_priv,
    timer,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING quiz_id, creator_id, quiz_title, description, is_priv, timer, created_at, updated_at
`

type CreateQuizParams struct {
	CreatorID   sql.NullInt32
	QuizTitle   string
	Description sql.NullString
	IsPriv      bool
	Timer       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreateQuiz(ctx context.Context, arg CreateQuizParams) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, createQuiz,
		arg.CreatorID,
		arg.QuizTitle,
		arg.Description,
		arg.IsPriv,
		arg.Timer,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Quiz
	err := row.Scan(
		&i.QuizID,
		&i.CreatorID,
		&i.QuizTitle,
		&i.Description,
		&i.IsPriv,
		&i.Timer,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
