-- name: ListAccounts :many
SELECT id, owner, balance, currency, created_at FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3; 