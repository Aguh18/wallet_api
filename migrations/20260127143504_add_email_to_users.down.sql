-- Rollback: Remove email column from users table
ALTER TABLE users DROP COLUMN IF EXISTS email;

-- Rollback: Remove index
DROP INDEX IF EXISTS idx_users_email;
