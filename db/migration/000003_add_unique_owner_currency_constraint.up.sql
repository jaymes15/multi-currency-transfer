ALTER TABLE accounts
ADD CONSTRAINT unique_owner_currency UNIQUE (owner, currency);