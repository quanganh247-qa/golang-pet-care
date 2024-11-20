<<<<<<< HEAD
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
=======
>>>>>>> 9fd7fc8 (feat: validate notification schema and APIs)
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
<<<<<<< HEAD
ORDER BY notifications.dueDate DESC
LIMIT $2 OFFSET $3;
>>>>>>> eb8d761 (updated pet schedule)
=======
ORDER BY notifications.datetime DESC
LIMIT $2 OFFSET $3;

-- name: IsReadNotification :exec
UPDATE notifications
SET is_read = true
WHERE notificationID = $1 ;
>>>>>>> 9fd7fc8 (feat: validate notification schema and APIs)
