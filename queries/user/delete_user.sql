-- name: DeleteUser :exec
DELETE FROM users
WHERE user_id = $1;
