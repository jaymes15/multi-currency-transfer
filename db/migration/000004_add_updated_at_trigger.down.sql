-- Down migration to remove updated_at trigger

-- Drop trigger
DROP TRIGGER IF EXISTS update_exchange_rates_updated_at ON exchange_rates;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column(); 