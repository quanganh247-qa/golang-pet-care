-- Drop indexes first
DROP INDEX IF EXISTS idx_offline_messages_client_id;
DROP INDEX IF EXISTS idx_offline_messages_username;
DROP INDEX IF EXISTS idx_offline_messages_status;

-- Drop the table
DROP TABLE IF EXISTS offline_messages; 