-- Migration: create_refresh_tokens_table (rollback)
-- Description: Drops refresh_tokens table

BEGIN;

DROP INDEX IF EXISTS idx_refresh_tokens_customer_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash;
DROP INDEX IF EXISTS idx_refresh_tokens_family_id;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP TABLE IF EXISTS refresh_tokens;

COMMIT;
