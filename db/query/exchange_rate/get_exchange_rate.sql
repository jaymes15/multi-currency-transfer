-- name: GetExchangeRate :one
SELECT id, from_currency, to_currency, rate, created_at, updated_at FROM exchange_rates
WHERE from_currency = $1 AND to_currency = $2
LIMIT 1; 