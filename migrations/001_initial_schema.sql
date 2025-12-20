-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Reward events table
CREATE TABLE IF NOT EXISTS reward_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID UNIQUE NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_symbol VARCHAR(20) NOT NULL,
    quantity NUMERIC(18,6) NOT NULL CHECK (quantity > 0),
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Ledger entries table (double-entry accounting)
CREATE TABLE IF NOT EXISTS ledger_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES reward_events(event_id) ON DELETE CASCADE,
    entry_type VARCHAR(10) NOT NULL CHECK (entry_type IN ('STOCK', 'CASH', 'FEE')),
    symbol VARCHAR(20),
    debit NUMERIC(18,4) NOT NULL DEFAULT 0,
    credit NUMERIC(18,4) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (debit >= 0 AND credit >= 0)
);

-- Stock prices table
CREATE TABLE IF NOT EXISTS stock_prices (
    symbol VARCHAR(20) PRIMARY KEY,
    price NUMERIC(18,4) NOT NULL CHECK (price > 0),
    fetched_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_reward_events_user_id ON reward_events(user_id);
CREATE INDEX IF NOT EXISTS idx_reward_events_timestamp ON reward_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_reward_events_event_id ON reward_events(event_id);
CREATE INDEX IF NOT EXISTS idx_ledger_entries_event_id ON ledger_entries(event_id);
CREATE INDEX IF NOT EXISTS idx_stock_prices_fetched_at ON stock_prices(fetched_at);

-- Function to verify ledger balance (for auditing)
CREATE OR REPLACE FUNCTION verify_ledger_balance()
RETURNS TABLE(total_debit NUMERIC, total_credit NUMERIC, is_balanced BOOLEAN) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        SUM(debit) as total_debit,
        SUM(credit) as total_credit,
        SUM(debit) = SUM(credit) as is_balanced
    FROM ledger_entries;
END;
$$ LANGUAGE plpgsql;

