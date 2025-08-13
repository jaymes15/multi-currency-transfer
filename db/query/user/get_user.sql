-- name: GetUser :one
SELECT username, full_name, email, created_at FROM users
WHERE username = $1 LIMIT 1; 