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

-- name: CreateQuestionMinimal :one
INSERT INTO questions (quiz_id, description, timer_option, timer)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE ques_id = $1;


-- name: GetQuestion :one
SELECT * FROM questions
WHERE ques_id = $1 LIMIT 1;

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

-- name: ListQuestionsByQuiz :many
SELECT * FROM questions
WHERE quiz_id = $1
ORDER BY ques_id;