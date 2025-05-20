-- -- +migrate Down
-- DROP TABLE IF EXISTS product_stock_movements;
-- DROP TYPE IF EXISTS movement_type_enum;

-- -- Drop indexes
-- DROP INDEX IF EXISTS idx_message_read_status_user_id;
-- DROP INDEX IF EXISTS idx_messages_created_at;
-- DROP INDEX IF EXISTS idx_messages_sender_id;
-- DROP INDEX IF EXISTS idx_messages_conversation_id;
-- DROP INDEX IF EXISTS idx_conversation_participants_user_id;

-- -- Drop tables in reverse order of creation (to avoid foreign key constraints)
-- DROP TABLE IF EXISTS message_read_status;
-- DROP TABLE IF EXISTS messages;
-- DROP TABLE IF EXISTS conversation_participants;
-- DROP TABLE IF EXISTS conversations;