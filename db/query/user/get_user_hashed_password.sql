-- name: GetUserHashedPassword :one
SELECT hashed_password FROM users
WHERE username = $1 LIMIT 1;