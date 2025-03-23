// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: get_user.sql

package db

import (
	"context"
)

const getUser = `-- name: GetUser :one
SELECT user_id, name, email, created_at, updated_at, deleted_at FROM users
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, userID int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Name,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
