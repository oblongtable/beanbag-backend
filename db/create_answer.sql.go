// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: create_answer.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createAnswer = `-- name: CreateAnswer :one
INSERT INTO answers (
    ques_id,
    description,
    is_correct,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING ans_id, ques_id, description, is_correct, created_at, updated_at
`

type CreateAnswerParams struct {
	QuesID      sql.NullInt32
	Description string
	IsCorrect   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (q *Queries) CreateAnswer(ctx context.Context, arg CreateAnswerParams) (Answer, error) {
	row := q.db.QueryRowContext(ctx, createAnswer,
		arg.QuesID,
		arg.Description,
		arg.IsCorrect,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Answer
	err := row.Scan(
		&i.AnsID,
		&i.QuesID,
		&i.Description,
		&i.IsCorrect,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
