-- name: CreateReminder :one
INSERT INTO Reminders (petid, title, description, due_date, repeat_interval, is_completed, notification_sent)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING reminder_id, petid, title, description, due_date, repeat_interval, is_completed, notification_sent;

-- name: GetRemindersByPetID :many
SELECT reminder_id, petid, title, description, due_date, repeat_interval, is_completed, notification_sent
FROM Reminders
WHERE petid = $1
ORDER BY due_date;

-- name: ListCompletedReminders :many
SELECT reminder_id, petid, title, description, due_date, repeat_interval, is_completed, notification_sent
FROM Reminders
WHERE is_completed = true
ORDER BY due_date;

-- name: UpdateReminder :exec
UPDATE Reminders
SET title = $2, description = $3, due_date = $4, repeat_interval = $5, is_completed = $6, notification_sent = $7
WHERE reminder_id = $1;

-- name: DeleteReminder :exec
DELETE FROM Reminders
WHERE reminder_id = $1;
