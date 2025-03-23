-- name: GetQuiz :one
SELECT * FROM quizzes
WHERE quiz_id = $1 LIMIT 1;
