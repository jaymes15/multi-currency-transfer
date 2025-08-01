-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3; 