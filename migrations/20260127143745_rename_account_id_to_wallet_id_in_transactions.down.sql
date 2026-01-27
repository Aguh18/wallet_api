-- Rollback: Rename column wallet_id back to account_id in transactions table
ALTER TABLE transactions RENAME COLUMN wallet_id TO account_id;

-- Rollback: Rename index
ALTER INDEX idx_transactions_wallet_id RENAME TO idx_transactions_account_id;

-- Rollback: Update foreign key constraint name if needed
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'transactions_wallet_id_fkey'
    ) THEN
        ALTER TABLE transactions
        DROP CONSTRAINT transactions_wallet_id_fkey,
        ADD CONSTRAINT transactions_account_id_fkey
        FOREIGN KEY (account_id) REFERENCES wallets(id) ON DELETE CASCADE;
    END IF;
END $$;
