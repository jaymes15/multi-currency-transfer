-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount,
  exchange_rate,
  from_currency,
  to_currency
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING id, from_account_id, to_account_id, amount, exchange_rate, from_currency, to_currency, created_at; 