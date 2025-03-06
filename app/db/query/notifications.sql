<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD

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
=======

-- name: ListNotificationsByUsername :many
SELECT * FROM notifications
WHERE username = $1
>>>>>>> e859654 (Elastic search)
LIMIT $2 OFFSET $3;

-- name: DeleteNotificationsByUsername :exec
DELETE FROM notifications
WHERE username = $1;

-- name: MarkNotificationAsRead :exec
UPDATE notifications
SET is_read = true
<<<<<<< HEAD
WHERE notificationID = $1 ;
>>>>>>> 9fd7fc8 (feat: validate notification schema and APIs)
=======
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
>>>>>>> e859654 (Elastic search)
