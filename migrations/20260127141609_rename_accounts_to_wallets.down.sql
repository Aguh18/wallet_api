-- Rollback: Rename table wallets back to accounts
ALTER TABLE wallets RENAME TO accounts;

-- Rollback: Rename columns in accounts table
ALTER TABLE accounts RENAME COLUMN wallet_name TO account_name;

-- Rollback: Rename columns in transactions table
ALTER TABLE transactions RENAME COLUMN wallet_id TO account_id;

-- Rollback: Rename indexes
ALTER INDEX idx_wallets_user_id RENAME TO idx_accounts_user_id;
ALTER INDEX idx_wallets_status RENAME TO idx_accounts_status;
ALTER INDEX idx_transactions_wallet_id RENAME TO idx_transactions_account_id;
