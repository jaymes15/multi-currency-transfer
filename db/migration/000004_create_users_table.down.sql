-- Drop foreign key constraint from accounts table
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

-- Drop unique constraint on accounts (owner, currency)
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

-- Drop users table
DROP TABLE IF EXISTS "users"; 