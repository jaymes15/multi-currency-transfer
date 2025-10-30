-- Add non-empty value constraints to users table
ALTER TABLE "users"
  ADD CONSTRAINT users_username_nonempty CHECK (char_length(trim(username)) > 0),
  ADD CONSTRAINT users_hashed_password_nonempty CHECK (char_length(trim(hashed_password)) > 0),
  ADD CONSTRAINT users_full_name_nonempty CHECK (char_length(trim(full_name)) > 0),
  ADD CONSTRAINT users_email_nonempty CHECK (char_length(trim(email)) > 0);

-- Ensure accounts.owner references users.username with ON DELETE CASCADE
-- Drop existing FK if present, then recreate with cascade
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
  FOREIGN KEY (owner) REFERENCES "users" ("username") ON DELETE CASCADE;


