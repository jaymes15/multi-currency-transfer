-- Add cross-currency support to transfers table
ALTER TABLE transfers ADD COLUMN exchange_rate DECIMAL(10,8);
ALTER TABLE transfers ADD COLUMN from_currency VARCHAR(3);
ALTER TABLE transfers ADD COLUMN to_currency VARCHAR(3);

-- Add index for currency pair lookups
CREATE INDEX idx_transfers_currency_pair ON transfers(from_currency, to_currency); 