-- Rename table accounts to wallets
ALTER TABLE accounts RENAME TO wallets;

-- Rename columns in wallets table
ALTER TABLE wallets RENAME COLUMN account_name TO wallet_name;

-- Rename columns in transactions table
ALTER TABLE transactions RENAME COLUMN account_id TO wallet_id;

-- Rename indexes
ALTER INDEX idx_accounts_user_id RENAME TO idx_wallets_user_id;
ALTER INDEX idx_accounts_status RENAME TO idx_wallets_status;
ALTER INDEX idx_transactions_account_id RENAME TO idx_transactions_wallet_id;
