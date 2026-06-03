-- Migration: create_refresh_tokens_table
-- Description: Creates refresh_tokens table for JWT refresh token rotation with reuse detection

BEGIN;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    customer_id UUID NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    family_id UUID NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_customer_id ON refresh_tokens(customer_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_family_id ON refresh_tokens(family_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

COMMENT ON TABLE refresh_tokens IS 'JWT refresh tokens for token rotation and reuse detection';
COMMENT ON COLUMN refresh_tokens.id IS 'Primary key (UUID)';
COMMENT ON COLUMN refresh_tokens.customer_id IS 'Foreign key to customers table';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA256 hash of the refresh token';
COMMENT ON COLUMN refresh_tokens.family_id IS 'UUID tracking token lineage for reuse detection';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'Expiration timestamp';
COMMENT ON COLUMN refresh_tokens.revoked IS 'Boolean value if token was revoked';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'Timestamp when token was revoked (NULL if active)';
COMMENT ON COLUMN refresh_tokens.created_at IS 'Timestamp when token was created';

COMMIT;
