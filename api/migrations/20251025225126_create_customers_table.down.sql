-- Migration: create_customers_table (rollback)
-- Description: Drops customers table and associated indexes

BEGIN;

DROP INDEX IF EXISTS idx_customers_phone;
DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_customers_deleted_at;
DROP TABLE IF EXISTS customers;

COMMIT;

