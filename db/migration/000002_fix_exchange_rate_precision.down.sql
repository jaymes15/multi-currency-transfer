-- Revert exchange_rate precision in transfers table
ALTER TABLE "transfers" 
ALTER COLUMN "exchange_rate" TYPE DECIMAL(10,8); 