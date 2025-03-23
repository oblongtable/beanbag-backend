-- name: CreateQuestion :one
INSERT INTO questions (
    quiz_id,
    description,
    timer_option,
    timer,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;
