-- name: ListAllAccounts :many
SELECT id, owner, balance, currency, created_at FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2; 