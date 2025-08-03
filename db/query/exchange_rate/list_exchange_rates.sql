-- name: ListExchangeRates :many
SELECT id, from_currency, to_currency, rate, created_at, updated_at FROM exchange_rates
ORDER BY from_currency, to_currency; 