-- Down migration for combined schema
-- Drop all tables in reverse order

-- Drop indexes first
DROP INDEX IF EXISTS idx_transfers_currency_pair;
DROP INDEX IF EXISTS idx_exchange_rates_currencies;
DROP INDEX IF EXISTS "transfers_from_account_id_to_account_id_idx";
DROP INDEX IF EXISTS "transfers_to_account_id_idx";
DROP INDEX IF EXISTS "transfers_from_account_id_idx";
DROP INDEX IF EXISTS "entries_account_id_idx";
DROP INDEX IF EXISTS "accounts_owner_idx";

-- Drop foreign key constraints
ALTER TABLE "transfers" DROP CONSTRAINT IF EXISTS "transfers_to_account_id_fkey";
ALTER TABLE "transfers" DROP CONSTRAINT IF EXISTS "transfers_from_account_id_fkey";
ALTER TABLE "entries" DROP CONSTRAINT IF EXISTS "entries_account_id_fkey";

-- Drop tables
DROP TABLE IF EXISTS "exchange_rates";
DROP TABLE IF EXISTS "transfers";
DROP TABLE IF EXISTS "entries";
DROP TABLE IF EXISTS "accounts"; 