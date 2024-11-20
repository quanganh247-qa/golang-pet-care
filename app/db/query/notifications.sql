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