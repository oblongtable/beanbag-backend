// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: update_answer.sql

package db

import (
	"context"
	"database/sql"
)

const updateAnswer = `-- name: UpdateAnswer :one
UPDATE answers
SET
    ques_id = COALESCE($2, ques_id),
    description = COALESCE($3, description),
    is_correct = COALESCE($4, is_correct),
    updated_at = NOW()
WHERE ans_id = $1
RETURNING ans_id, ques_id, description, is_correct, created_at, updated_at
`

type UpdateAnswerParams struct {
	AnsID       int32
	QuesID      sql.NullInt32
	Description sql.NullString
	IsCorrect   sql.NullBool
}

func (q *Queries) UpdateAnswer(ctx context.Context, arg UpdateAnswerParams) (Answer, error) {
	row := q.db.QueryRowContext(ctx, updateAnswer,
		arg.AnsID,
		arg.QuesID,
		arg.Description,
		arg.IsCorrect,
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
