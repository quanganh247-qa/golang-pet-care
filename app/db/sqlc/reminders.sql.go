// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: reminders.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createReminder = `-- name: CreateReminder :one
INSERT INTO Reminders (petid, title, description, due_date, repeat_interval, is_completed, notification_sent)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING reminder_id, petid, title, description, due_date, repeat_interval, is_completed, notification_sent
`

type CreateReminderParams struct {
	Petid            pgtype.Int8      `json:"petid"`
	Title            string           `json:"title"`
	Description      pgtype.Text      `json:"description"`
	DueDate          pgtype.Timestamp `json:"due_date"`
	RepeatInterval   pgtype.Text      `json:"repeat_interval"`
	IsCompleted      pgtype.Bool      `json:"is_completed"`
	NotificationSent pgtype.Bool      `json:"notification_sent"`
}

func (q *Queries) CreateReminder(ctx context.Context, arg CreateReminderParams) (Reminder, error) {
	row := q.db.QueryRow(ctx, createReminder,
		arg.Petid,
		arg.Title,
		arg.Description,
		arg.DueDate,
		arg.RepeatInterval,
		arg.IsCompleted,
		arg.NotificationSent,
	)
	var i Reminder
	err := row.Scan(
		&i.ReminderID,
		&i.Petid,
		&i.Title,
		&i.Description,
		&i.DueDate,
		&i.RepeatInterval,
		&i.IsCompleted,
		&i.NotificationSent,
	)
	return i, err
}

const deleteReminder = `-- name: DeleteReminder :exec
DELETE FROM Reminders
WHERE reminder_id = $1
`

func (q *Queries) DeleteReminder(ctx context.Context, reminderID int64) error {
	_, err := q.db.Exec(ctx, deleteReminder, reminderID)
	return err
}

const getRemindersByPetID = `-- name: GetRemindersByPetID :many
SELECT reminder_id, petid, title, description, due_date, repeat_interval, is_completed, notification_sent
FROM Reminders
WHERE petid = $1
ORDER BY due_date
`

func (q *Queries) GetRemindersByPetID(ctx context.Context, petid pgtype.Int8) ([]Reminder, error) {
	rows, err := q.db.Query(ctx, getRemindersByPetID, petid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Reminder{}
	for rows.Next() {
		var i Reminder
		if err := rows.Scan(
			&i.ReminderID,
			&i.Petid,
			&i.Title,
			&i.Description,
			&i.DueDate,
			&i.RepeatInterval,
			&i.IsCompleted,
			&i.NotificationSent,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listCompletedReminders = `-- name: ListCompletedReminders :many
SELECT reminder_id, petid, title, description, due_date, repeat_interval, is_completed, notification_sent
FROM Reminders
WHERE is_completed = true
ORDER BY due_date
`

func (q *Queries) ListCompletedReminders(ctx context.Context) ([]Reminder, error) {
	rows, err := q.db.Query(ctx, listCompletedReminders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Reminder{}
	for rows.Next() {
		var i Reminder
		if err := rows.Scan(
			&i.ReminderID,
			&i.Petid,
			&i.Title,
			&i.Description,
			&i.DueDate,
			&i.RepeatInterval,
			&i.IsCompleted,
			&i.NotificationSent,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateReminder = `-- name: UpdateReminder :exec
UPDATE Reminders
SET title = $2, description = $3, due_date = $4, repeat_interval = $5, is_completed = $6, notification_sent = $7
WHERE reminder_id = $1
`

type UpdateReminderParams struct {
	ReminderID       int64            `json:"reminder_id"`
	Title            string           `json:"title"`
	Description      pgtype.Text      `json:"description"`
	DueDate          pgtype.Timestamp `json:"due_date"`
	RepeatInterval   pgtype.Text      `json:"repeat_interval"`
	IsCompleted      pgtype.Bool      `json:"is_completed"`
	NotificationSent pgtype.Bool      `json:"notification_sent"`
}

func (q *Queries) UpdateReminder(ctx context.Context, arg UpdateReminderParams) error {
	_, err := q.db.Exec(ctx, updateReminder,
		arg.ReminderID,
		arg.Title,
		arg.Description,
		arg.DueDate,
		arg.RepeatInterval,
		arg.IsCompleted,
		arg.NotificationSent,
	)
	return err
}