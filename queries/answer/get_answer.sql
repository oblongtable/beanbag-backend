-- name: GetAnswer :one
SELECT * FROM answers
WHERE ans_id = $1 LIMIT 1;
