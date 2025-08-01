-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at FROM transfers
WHERE id = $1 LIMIT 1; 