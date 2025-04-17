-- name: CreateOfflineMessage :one
INSERT INTO offline_messages (
    client_id,
    username,
    message_type,
    data,
    status,
    retry_count
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetPendingMessagesForClient :many
SELECT * FROM offline_messages
WHERE  status = 'pending'
ORDER BY created_at ASC;

-- name: MarkMessageDelivered :exec
UPDATE offline_messages
SET 
    status = 'delivered',
    delivered_at = NOW()
WHERE id = $1;

-- name: IncrementMessageRetryCount :exec
UPDATE offline_messages
SET 
    retry_count = retry_count + 1,
    status = CASE 
        WHEN retry_count >= 5 THEN 'failed'
        ELSE 'pending'
    END
WHERE id = $1;

-- name: DeleteOldMessages :exec
DELETE FROM offline_messages
WHERE status = $1 AND created_at < $2;

-- name: GetMessagesForRetry :many
SELECT * FROM offline_messages
WHERE status = ANY($1::varchar[]) AND retry_count < $2
ORDER BY created_at ASC; 