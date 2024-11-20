-- name: InsertNotification :one
INSERT INTO notifications (username, title, description,datetime)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteAllNotificationsByUser :exec
DELETE FROM notifications WHERE username = $1;

-- name: DeleteNotificationByID :exec
DELETE FROM notifications
WHERE notificationID = $1;

-- name: GetNotificationsByUsername :many
SELECT notifications.*
FROM notifications
JOIN users ON notifications.username = users.username
WHERE users.username = $1
ORDER BY notifications.datetime DESC
LIMIT $2 OFFSET $3;

-- name: IsReadNotification :exec
UPDATE notifications
SET is_read = true
WHERE notificationID = $1 ;