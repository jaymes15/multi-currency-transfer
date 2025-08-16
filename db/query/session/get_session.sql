-- name: GetSession :one
SELECT id, username, client_ip, user_agent, is_blocked, expires_at, created_at
FROM sessions
WHERE id = $1 LIMIT 1; 