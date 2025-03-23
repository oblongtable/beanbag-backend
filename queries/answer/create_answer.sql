-- name: CreateAnswer :one
INSERT INTO answers (
    ques_id,
    description,
    is_correct,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;
