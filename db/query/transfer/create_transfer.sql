-- name: CreateTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount,
  converted_amount,
  exchange_rate,
  from_currency,
  to_currency,
  fee
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id, from_account_id, to_account_id, amount, converted_amount, exchange_rate, from_currency, to_currency, fee, created_at; 