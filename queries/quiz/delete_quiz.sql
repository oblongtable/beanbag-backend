-- name: DeleteQuiz :exec
DELETE FROM quizzes
WHERE quiz_id = $1;
