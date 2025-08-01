-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1 LIMIT 1; 