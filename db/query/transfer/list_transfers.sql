-- name: ListTransfers :many
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE 
    from_account_id = $1 OR
    to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4; 