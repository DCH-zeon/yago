-- Migration: create_customers_table
-- Description: Creates customers table with indexes and constraints

BEGIN;

-- Create a customers table
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(36),
    first_name VARCHAR(36),
    last_name VARCHAR(36),
    email VARCHAR(255),
    phone VARCHAR(16) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_deleted_at ON customers(deleted_at);

COMMENT ON TABLE customers IS 'Application customers table';
COMMENT ON COLUMN customers.id IS 'Primary key';
COMMENT ON COLUMN customers.external_id IS 'Customer id from 1C (optional)';
COMMENT ON COLUMN customers.first_name IS 'Customer first name';
COMMENT ON COLUMN customers.last_name IS 'Customer last name';
COMMENT ON COLUMN customers.email IS 'Customer email address (optional)';
COMMENT ON COLUMN customers.phone IS 'Customer phone number (unique, required for auth)';
COMMENT ON COLUMN customers.password_hash IS 'Bcrypt hashed password';
COMMENT ON COLUMN customers.created_at IS 'Timestamp when customer was created';
COMMENT ON COLUMN customers.updated_at IS 'Timestamp when customer was last updated';
COMMENT ON COLUMN customers.deleted_at IS 'Timestamp when customer was soft deleted (NULL if active)';

COMMIT;

