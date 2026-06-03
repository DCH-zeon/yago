BEGIN;

DROP INDEX IF EXISTS idx_customer_roles_customer_id;
DROP INDEX IF EXISTS idx_customer_roles_role_id;
DROP TABLE IF EXISTS customer_roles CASCADE;

COMMIT;
