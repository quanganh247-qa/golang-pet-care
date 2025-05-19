-- name: CreatetNotification :one
INSERT INTO notifications (
    username,
    title,
    content,
    notify_type,
    related_id,
    related_type,
    is_read,
    datetime
) VALUES (
    $1, $2, $3, $4, $5, $6, false, NOW()
) RETURNING *;

-- name: ListNotification :many
SELECT * FROM notifications LIMIT $1 OFFSET $2;

-- name: DeleteNotificationsByUsername :exec
DELETE FROM notifications
WHERE username = $1;

-- name: MarkNotificationAsRead :exec
UPDATE notifications
SET is_read = true
WHERE id = $1;
