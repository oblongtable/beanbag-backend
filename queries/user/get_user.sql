-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;
