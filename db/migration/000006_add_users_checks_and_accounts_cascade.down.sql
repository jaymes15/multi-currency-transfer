-- Remove non-empty value constraints from users table
ALTER TABLE "users"
  DROP CONSTRAINT IF EXISTS users_username_nonempty,
  DROP CONSTRAINT IF EXISTS users_hashed_password_nonempty,
  DROP CONSTRAINT IF EXISTS users_full_name_nonempty,
  DROP CONSTRAINT IF EXISTS users_email_nonempty;

-- Recreate accounts.owner foreign key without ON DELETE CASCADE (default NO ACTION)
DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.table_constraints tc
    WHERE tc.constraint_type = 'FOREIGN KEY'
      AND tc.table_name = 'accounts'
      AND tc.constraint_name = 'accounts_owner_fkey'
  ) THEN
    ALTER TABLE "accounts" DROP CONSTRAINT accounts_owner_fkey;
  END IF;
END $$;

ALTER TABLE "accounts"
  ADD CONSTRAINT accounts_owner_fkey
  FOREIGN KEY (owner) REFERENCES "users" ("username");


