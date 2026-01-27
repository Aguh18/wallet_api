-- Change balance column type from BIGINT to NUMERIC for wallets table
ALTER TABLE wallets ALTER COLUMN balance TYPE NUMERIC(20, 2);

-- Add comment for clarity
COMMENT ON COLUMN wallets.balance IS 'Wallet balance in currency unit (e.g., 100.50 = 100 dollars and 50 cents)';

-- Change balance columns in transactions table to NUMERIC
ALTER TABLE transactions ALTER COLUMN amount TYPE NUMERIC(20, 2);
ALTER TABLE transactions ALTER COLUMN balance_before TYPE NUMERIC(20, 2);
ALTER TABLE transactions ALTER COLUMN balance_after TYPE NUMERIC(20, 2);

-- Add comments for transaction columns
COMMENT ON COLUMN transactions.amount IS 'Transaction amount in currency unit';
COMMENT ON COLUMN transactions.balance_before IS 'Balance before transaction in currency unit';
COMMENT ON COLUMN transactions.balance_after IS 'Balance after transaction in currency unit';
