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

-- name: CreateQuizMinimal :one
INSERT INTO quizzes (quiz_title, creator_id)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteQuiz :exec
DELETE FROM quizzes
WHERE quiz_id = $1;

-- name: GetQuiz :one
SELECT * FROM quizzes
WHERE quiz_id = $1 LIMIT 1;

-- name: UpdateQuiz :one
UPDATE quizzes
SET
    creator_id = COALESCE(sqlc.narg(creator_id), creator_id),
    quiz_title = COALESCE(sqlc.narg(quiz_title), quiz_title),
    description = COALESCE(sqlc.narg(description), description),
    is_priv = COALESCE(sqlc.narg(is_priv), is_priv),
    timer = COALESCE(sqlc.narg(timer), timer),
    updated_at = NOW()
WHERE quiz_id = $1
RETURNING *;
