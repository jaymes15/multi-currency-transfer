-- Combined schema migration for SimpleBank with exchange rate precision fixes
-- This includes all tables with proper DECIMAL types for financial precision

-- Create accounts table with DECIMAL balance
CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" DECIMAL(20,2) NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create entries table with DECIMAL amount
CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" DECIMAL(20,2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create transfers table with DECIMAL amounts and cross-currency support
-- Using DECIMAL(20,8) for exchange_rate to handle large rates like GBP to NGN (2000+)
CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" DECIMAL(20,2) NOT NULL,
  "converted_amount" DECIMAL(20,2),
  "exchange_rate" DECIMAL(20,8),
  "from_currency" VARCHAR(3),
  "to_currency" VARCHAR(3),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Create exchange rates table
CREATE TABLE exchange_rates (
    id BIGSERIAL PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(from_currency, to_currency)
);

-- Add foreign key constraints
ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

-- Create indexes for performance
CREATE INDEX ON "accounts" ("owner");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");
CREATE INDEX idx_exchange_rates_currencies ON exchange_rates(from_currency, to_currency);
CREATE INDEX idx_transfers_currency_pair ON transfers(from_currency, to_currency);

-- Add comments for documentation
COMMENT ON COLUMN "entries"."amount" IS 'Entry amount (positive for credits, negative for debits)';
COMMENT ON COLUMN "transfers"."amount" IS 'Original transfer amount in source currency';
COMMENT ON COLUMN "transfers"."converted_amount" IS 'Amount after currency conversion in target currency';
COMMENT ON COLUMN "accounts"."balance" IS 'Account balance in account currency'; 