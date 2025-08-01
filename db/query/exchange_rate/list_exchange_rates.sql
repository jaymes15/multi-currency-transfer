-- name: ListExchangeRates :many
SELECT id, from_currency, to_currency, rate::float8, created_at FROM exchange_rates
ORDER BY from_currency, to_currency; 