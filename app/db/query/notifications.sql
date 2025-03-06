-- name: CreatetNotification :one
INSERT INTO notifications (
    username,
    title,
    content,
    notify_type,
    related_id,
    related_type,
    is_read
) VALUES (
    $1, $2, $3, $4, $5, $6, false
) RETURNING *;

-- name: ListNotificationsByUsername :many
SELECT * FROM notifications
WHERE username = $1
LIMIT $2 OFFSET $3;

-- name: DeleteNotificationsByUsername :exec
DELETE FROM notifications
WHERE username = $1;

-- name: MarkNotificationAsRead :exec
UPDATE notifications
SET is_read = true
WHERE id = $1;

-- name: CreateNotificationPreference :exec
INSERT INTO notification_preferences (
    username,
    topic,
    enabled
) VALUES (
    $1, $2, $3
);

-- name: UpdateNotificationPreference :exec
UPDATE notification_preferences
SET enabled = $2
WHERE username = $1 AND topic = $3;

-- name: GetNotificationPreferencesByUsername :many
SELECT * FROM notification_preferences
WHERE username = $1;