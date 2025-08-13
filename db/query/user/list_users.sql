-- name: ListUsers :many
SELECT username, full_name, email, created_at FROM users
ORDER BY username; 