-- Rename column account_id to wallet_id in transactions table
ALTER TABLE transactions RENAME COLUMN account_id TO wallet_id;

-- Rename index
ALTER INDEX idx_transactions_account_id RENAME TO idx_transactions_wallet_id;

-- Update foreign key constraint name if needed
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'transactions_account_id_fkey'
    ) THEN
        ALTER TABLE transactions
        DROP CONSTRAINT transactions_account_id_fkey,
        ADD CONSTRAINT transactions_wallet_id_fkey
        FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE;
    END IF;
END $$;
