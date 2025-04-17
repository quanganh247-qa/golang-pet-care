-- Drop indexes
DROP INDEX IF EXISTS idx_smtp_configs_email;
DROP INDEX IF EXISTS idx_smtp_configs_is_default;

-- Drop table
DROP TABLE IF EXISTS smtp_configs; 