-- name: GetAccountByOwnerAndCurrency :one
SELECT * FROM accounts
WHERE owner = $1 AND currency = $2
LIMIT 1;