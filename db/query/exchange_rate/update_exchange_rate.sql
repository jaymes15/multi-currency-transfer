-- name: UpdateExchangeRate :one
UPDATE exchange_rates
SET rate = $3, created_at = NOW()
WHERE from_currency = $1 AND to_currency = $2
RETURNING id, from_currency, to_currency, rate, created_at, updated_at; 