CREATE TABLE exchange_rates (
    id BIGSERIAL PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(from_currency, to_currency)
);

-- Currency column already exists in accounts table

-- Create index for faster lookups
CREATE INDEX idx_exchange_rates_currencies ON exchange_rates(from_currency, to_currency); 