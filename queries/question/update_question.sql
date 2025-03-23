-- name: UpdateQuestion :one
UPDATE questions
SET
    quiz_id = COALESCE(sqlc.narg(quiz_id), quiz_id),
    description = COALESCE(sqlc.narg(description), description),
    timer_option = COALESCE(sqlc.narg(timer_option), timer_option),
    timer = COALESCE(sqlc.narg(timer), timer),
    updated_at = NOW()
WHERE ques_id = $1
RETURNING *;
