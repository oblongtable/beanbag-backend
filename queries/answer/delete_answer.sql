-- name: DeleteAnswer :exec
DELETE FROM answers
WHERE ans_id = $1;
