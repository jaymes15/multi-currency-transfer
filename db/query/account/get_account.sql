-- name: GetAccount :one
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE id = $1 LIMIT 1; 