// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: quiz.sql

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

const createQuizMinimal = `-- name: CreateQuizMinimal :one
INSERT INTO quizzes (quiz_title, creator_id)
VALUES ($1, $2)
RETURNING quiz_id, creator_id, quiz_title, description, is_priv, timer, created_at, updated_at
`

type CreateQuizMinimalParams struct {
	QuizTitle string
	CreatorID sql.NullInt32
}

func (q *Queries) CreateQuizMinimal(ctx context.Context, arg CreateQuizMinimalParams) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, createQuizMinimal, arg.QuizTitle, arg.CreatorID)
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

const deleteQuiz = `-- name: DeleteQuiz :exec
DELETE FROM quizzes
WHERE quiz_id = $1
`

func (q *Queries) DeleteQuiz(ctx context.Context, quizID int32) error {
	_, err := q.db.ExecContext(ctx, deleteQuiz, quizID)
	return err
}

const getQuiz = `-- name: GetQuiz :one
SELECT quiz_id, creator_id, quiz_title, description, is_priv, timer, created_at, updated_at FROM quizzes
WHERE quiz_id = $1 LIMIT 1
`

func (q *Queries) GetQuiz(ctx context.Context, quizID int32) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, getQuiz, quizID)
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

const updateQuiz = `-- name: UpdateQuiz :one
UPDATE quizzes
SET
    creator_id = COALESCE($2, creator_id),
    quiz_title = COALESCE($3, quiz_title),
    description = COALESCE($4, description),
    is_priv = COALESCE($5, is_priv),
    timer = COALESCE($6, timer),
    updated_at = NOW()
WHERE quiz_id = $1
RETURNING quiz_id, creator_id, quiz_title, description, is_priv, timer, created_at, updated_at
`

type UpdateQuizParams struct {
	QuizID      int32
	CreatorID   sql.NullInt32
	QuizTitle   sql.NullString
	Description sql.NullString
	IsPriv      sql.NullBool
	Timer       sql.NullInt32
}

func (q *Queries) UpdateQuiz(ctx context.Context, arg UpdateQuizParams) (Quiz, error) {
	row := q.db.QueryRowContext(ctx, updateQuiz,
		arg.QuizID,
		arg.CreatorID,
		arg.QuizTitle,
		arg.Description,
		arg.IsPriv,
		arg.Timer,
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
