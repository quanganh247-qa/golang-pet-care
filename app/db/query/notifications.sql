<<<<<<< HEAD
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
=======
-- Insert a notification
-- name: InsertNotification :one
INSERT INTO notifications (petID, title, body, dueDate, repeatInterval, isCompleted, notificationSent)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- Delete all notifications
-- name: DeleteAllNotifications :exec
DELETE FROM notifications;

-- name: GetNotificationsByUsername :many
SELECT notifications.*
FROM notifications
JOIN pet ON notifications.petID = pet.petid
JOIN users ON pet.username = users.username
WHERE users.username = $1
ORDER BY notifications.dueDate DESC
LIMIT $2 OFFSET $3;
>>>>>>> eb8d761 (updated pet schedule)
