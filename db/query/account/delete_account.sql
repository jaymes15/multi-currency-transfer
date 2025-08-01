-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1; 