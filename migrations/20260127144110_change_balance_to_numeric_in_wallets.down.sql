-- Rollback: Change balance column type back from NUMERIC to BIGINT
-- Note: This will lose decimal precision! Only rollback if you're okay with that.

-- First, convert to NUMERIC without decimals, then to BIGINT
ALTER TABLE wallets ALTER COLUMN balance TYPE NUMERIC(20, 0) USING balance::numeric::bigint;
ALTER TABLE wallets ALTER COLUMN balance TYPE BIGINT;

-- Remove comments
COMMENT ON COLUMN wallets.balance IS NULL;

-- Rollback transaction columns to BIGINT
ALTER TABLE transactions ALTER COLUMN amount TYPE NUMERIC(20, 0) USING amount::numeric::bigint;
ALTER TABLE transactions ALTER COLUMN amount TYPE BIGINT;

ALTER TABLE transactions ALTER COLUMN balance_before TYPE NUMERIC(20, 0) USING balance_before::numeric::bigint;
ALTER TABLE transactions ALTER COLUMN balance_before TYPE BIGINT;

ALTER TABLE transactions ALTER COLUMN balance_after TYPE NUMERIC(20, 0) USING balance_after::numeric::bigint;
ALTER TABLE transactions ALTER COLUMN balance_after TYPE BIGINT;

-- Remove comments
COMMENT ON COLUMN transactions.amount IS NULL;
COMMENT ON COLUMN transactions.balance_before IS NULL;
COMMENT ON COLUMN transactions.balance_after IS NULL;
