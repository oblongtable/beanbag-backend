-- name: CreateQuiz :one
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
) RETURNING *;
