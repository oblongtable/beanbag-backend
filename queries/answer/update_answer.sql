-- name: UpdateAnswer :one
UPDATE answers
SET
    ques_id = COALESCE(sqlc.narg(ques_id), ques_id),
    description = COALESCE(sqlc.narg(description), description),
    is_correct = COALESCE(sqlc.narg(is_correct), is_correct),
    updated_at = NOW()
WHERE ans_id = $1
RETURNING *;
