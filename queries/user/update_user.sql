-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE(sqlc.narg(name), name),
    email = COALESCE(sqlc.narg(email), email),
    updated_at = NOW()
WHERE user_id = $1
RETURNING *;
