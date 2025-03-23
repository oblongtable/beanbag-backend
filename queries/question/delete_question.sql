-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE ques_id = $1;
