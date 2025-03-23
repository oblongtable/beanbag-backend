-- name: GetQuestion :one
SELECT * FROM questions
WHERE ques_id = $1 LIMIT 1;
