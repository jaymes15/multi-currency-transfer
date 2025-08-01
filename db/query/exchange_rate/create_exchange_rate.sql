-- name: CreateExchangeRate :one
INSERT INTO exchange_rates (
  from_currency,
  to_currency,
  rate
) VALUES (
  $1, $2, $3
) RETURNING id, from_currency, to_currency, rate::float8, created_at; 