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

-- name: CreateAnswerMinimal :one
INSERT INTO answers (
    ques_id, description, is_correct
) VALUES (
    $1, $2, $3
) returning *;

-- name: DeleteAnswer :exec
DELETE FROM answers
WHERE ans_id = $1;

-- name: GetAnswer :one
SELECT * FROM answers
WHERE ans_id = $1 LIMIT 1;


-- name: UpdateAnswer :one
UPDATE answers
SET
    ques_id = COALESCE(sqlc.narg(ques_id), ques_id),
    description = COALESCE(sqlc.narg(description), description),
    is_correct = COALESCE(sqlc.narg(is_correct), is_correct),
    updated_at = NOW()
WHERE ans_id = $1
RETURNING *;

-- name: ListAnswersByQuestionIDs :many
SELECT * FROM answers
WHERE ques_id = ANY($1::int[])
ORDER BY ques_id, ans_id;