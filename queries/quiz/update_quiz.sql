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
