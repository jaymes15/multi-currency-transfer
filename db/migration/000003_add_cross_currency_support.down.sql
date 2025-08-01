-- Remove cross-currency support from transfers table
DROP INDEX IF EXISTS idx_transfers_currency_pair;
ALTER TABLE transfers DROP COLUMN IF EXISTS to_currency;
ALTER TABLE transfers DROP COLUMN IF EXISTS from_currency;
ALTER TABLE transfers DROP COLUMN IF EXISTS exchange_rate; 