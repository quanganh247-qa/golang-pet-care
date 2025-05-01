-- -- name: CreateConversation :one
-- INSERT INTO conversations (
--   type, 
--   name
-- ) VALUES (
--   $1, $2
-- ) RETURNING *;

-- -- name: GetConversation :one
-- SELECT * FROM conversations
-- WHERE id = $1;

-- -- name: GetUserConversations :many
-- SELECT c.* 
-- FROM conversations c
-- JOIN conversation_participants cp ON c.id = cp.conversation_id
-- WHERE cp.user_id = $1 AND cp.left_at IS NULL
-- ORDER BY c.updated_at DESC;

-- -- name: AddParticipantToConversation :one
-- INSERT INTO conversation_participants (
--   conversation_id,
--   user_id,
--   is_admin
-- ) VALUES (
--   $1, $2, $3
-- ) RETURNING *;

-- -- name: RemoveParticipantFromConversation :exec
-- UPDATE conversation_participants
-- SET left_at = NOW()
-- WHERE conversation_id = $1 AND user_id = $2 AND left_at IS NULL;

-- -- name: GetConversationParticipants :many
-- SELECT u.id, u.username, u.email, u.full_name, cp.is_admin, cp.joined_at
-- FROM conversation_participants cp
-- JOIN users u ON cp.user_id = u.id
-- WHERE cp.conversation_id = $1 AND cp.left_at IS NULL;

-- -- name: CreateMessage :one
-- INSERT INTO messages (
--   conversation_id,
--   sender_id,
--   content,
--   message_type,
--   metadata
-- ) VALUES (
--   $1, $2, $3, $4, $5
-- ) RETURNING *;

-- -- name: GetMessagesByConversation :many
-- SELECT m.*, u.username as sender_username, u.full_name as sender_name
-- FROM messages m
-- JOIN users u ON m.sender_id = u.id
-- WHERE m.conversation_id = $1
-- ORDER BY m.created_at
-- LIMIT $2
-- OFFSET $3;

-- -- name: GetConversationLastMessage :one
-- SELECT m.*, u.username as sender_username, u.full_name as sender_name
-- FROM messages m
-- JOIN users u ON m.sender_id = u.id
-- WHERE m.conversation_id = $1
-- ORDER BY m.created_at DESC
-- LIMIT 1;

-- -- name: MarkMessageAsRead :one
-- INSERT INTO message_read_status (
--   message_id,
--   user_id
-- ) VALUES (
--   $1, $2
-- ) ON CONFLICT DO NOTHING
-- RETURNING *;

-- -- name: GetUnreadMessageCount :one
-- SELECT COUNT(*) 
-- FROM messages m
-- LEFT JOIN message_read_status mrs ON m.id = mrs.message_id AND mrs.user_id = $2
-- WHERE m.conversation_id = $1 
-- AND mrs.message_id IS NULL
-- AND m.sender_id != $2;

-- -- name: IsUserInConversation :one
-- SELECT EXISTS(
--   SELECT 1 
--   FROM conversation_participants 
--   WHERE conversation_id = $1 AND user_id = $2 AND left_at IS NULL
-- ) AS is_member;

-- -- name: UpdateConversationTimestamp :exec
-- UPDATE conversations
-- SET updated_at = NOW()
-- WHERE id = $1;