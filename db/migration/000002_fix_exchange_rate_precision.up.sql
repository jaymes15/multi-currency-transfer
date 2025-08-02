-- Fix exchange_rate precision in transfers table
-- The current DECIMAL(10,8) is too small for exchange rates like 2037.43 (GBP to NGN)
-- Change to DECIMAL(20,8) to match the exchange_rates table

ALTER TABLE "transfers" 
ALTER COLUMN "exchange_rate" TYPE DECIMAL(20,8); 