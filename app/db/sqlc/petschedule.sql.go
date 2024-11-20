// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: petschedule.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const activeReminder = `-- name: ActiveReminder :exec
UPDATE pet_schedule
SET is_active = $2
WHERE id = $1
`

type ActiveReminderParams struct {
	ID       int64       `json:"id"`
	IsActive pgtype.Bool `json:"is_active"`
}

func (q *Queries) ActiveReminder(ctx context.Context, arg ActiveReminderParams) error {
	_, err := q.db.Exec(ctx, activeReminder, arg.ID, arg.IsActive)
	return err
}

const createPetSchedule = `-- name: CreatePetSchedule :exec
INSERT INTO pet_schedule (
   pet_id,
   title,
   reminder_datetime,
   event_repeat,
   end_type,
   end_date,
   notes,
   is_active
) VALUES ($1, $2, $3, $4, $5, $6, $7, false)
`

type CreatePetScheduleParams struct {
	PetID            pgtype.Int8      `json:"pet_id"`
	Title            pgtype.Text      `json:"title"`
	ReminderDatetime pgtype.Timestamp `json:"reminder_datetime"`
	EventRepeat      pgtype.Text      `json:"event_repeat"`
	EndType          pgtype.Bool      `json:"end_type"`
	EndDate          pgtype.Date      `json:"end_date"`
	Notes            pgtype.Text      `json:"notes"`
}

func (q *Queries) CreatePetSchedule(ctx context.Context, arg CreatePetScheduleParams) error {
	_, err := q.db.Exec(ctx, createPetSchedule,
		arg.PetID,
		arg.Title,
		arg.ReminderDatetime,
		arg.EventRepeat,
		arg.EndType,
		arg.EndDate,
		arg.Notes,
	)
	return err
}

const deletePetSchedule = `-- name: DeletePetSchedule :exec
Update pet_schedule
SET removedat = now()
WHERE id = $1
`

func (q *Queries) DeletePetSchedule(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deletePetSchedule, id)
	return err
}

const getAllSchedulesByPet = `-- name: GetAllSchedulesByPet :many
SELECT id, pet_id, title, reminder_datetime, event_repeat, end_type, end_date, notes, created_at, is_active, removedat FROM pet_schedule 
WHERE pet_id = $1 and removedat is null
ORDER BY reminder_datetime 
LIMIT $2 OFFSET $3
`

type GetAllSchedulesByPetParams struct {
	PetID  pgtype.Int8 `json:"pet_id"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

func (q *Queries) GetAllSchedulesByPet(ctx context.Context, arg GetAllSchedulesByPetParams) ([]PetSchedule, error) {
	rows, err := q.db.Query(ctx, getAllSchedulesByPet, arg.PetID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PetSchedule{}
	for rows.Next() {
		var i PetSchedule
		if err := rows.Scan(
			&i.ID,
			&i.PetID,
			&i.Title,
			&i.ReminderDatetime,
			&i.EventRepeat,
			&i.EndType,
			&i.EndDate,
			&i.Notes,
			&i.CreatedAt,
			&i.IsActive,
			&i.Removedat,
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

const getPetScheduleById = `-- name: GetPetScheduleById :one
SELECT id, pet_id, title, reminder_datetime, event_repeat, end_type, end_date, notes, created_at, is_active, removedat FROM pet_schedule
WHERE id = $1
`

func (q *Queries) GetPetScheduleById(ctx context.Context, id int64) (PetSchedule, error) {
	row := q.db.QueryRow(ctx, getPetScheduleById, id)
	var i PetSchedule
	err := row.Scan(
		&i.ID,
		&i.PetID,
		&i.Title,
		&i.ReminderDatetime,
		&i.EventRepeat,
		&i.EndType,
		&i.EndDate,
		&i.Notes,
		&i.CreatedAt,
		&i.IsActive,
		&i.Removedat,
	)
	return i, err
}

const listPetSchedulesByUsername = `-- name: ListPetSchedulesByUsername :many
SELECT pet_schedule.id, pet_schedule.pet_id, pet_schedule.title, pet_schedule.reminder_datetime, pet_schedule.event_repeat, pet_schedule.end_type, pet_schedule.end_date, pet_schedule.notes, pet_schedule.created_at, pet_schedule.is_active, pet_schedule.removedat, pet.name
FROM pet_schedule
LEFT JOIN pet ON pet_schedule.pet_id = pet.petid
LEFT JOIN users ON pet.username = users.username
WHERE users.username = $1 and pet_schedule.removedat is null
ORDER BY pet.petid, pet_schedule.reminder_datetime
`

type ListPetSchedulesByUsernameRow struct {
	ID               int64            `json:"id"`
	PetID            pgtype.Int8      `json:"pet_id"`
	Title            pgtype.Text      `json:"title"`
	ReminderDatetime pgtype.Timestamp `json:"reminder_datetime"`
	EventRepeat      pgtype.Text      `json:"event_repeat"`
	EndType          pgtype.Bool      `json:"end_type"`
	EndDate          pgtype.Date      `json:"end_date"`
	Notes            pgtype.Text      `json:"notes"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
	IsActive         pgtype.Bool      `json:"is_active"`
	Removedat        pgtype.Timestamp `json:"removedat"`
	Name             pgtype.Text      `json:"name"`
}

func (q *Queries) ListPetSchedulesByUsername(ctx context.Context, username string) ([]ListPetSchedulesByUsernameRow, error) {
	rows, err := q.db.Query(ctx, listPetSchedulesByUsername, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListPetSchedulesByUsernameRow{}
	for rows.Next() {
		var i ListPetSchedulesByUsernameRow
		if err := rows.Scan(
			&i.ID,
			&i.PetID,
			&i.Title,
			&i.ReminderDatetime,
			&i.EventRepeat,
			&i.EndType,
			&i.EndDate,
			&i.Notes,
			&i.CreatedAt,
			&i.IsActive,
			&i.Removedat,
			&i.Name,
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

const updatePetSchedule = `-- name: UpdatePetSchedule :exec
UPDATE pet_schedule
SET title = $2,
    reminder_datetime = $3,
    event_repeat = $4,
    end_type = $5,
    end_date = $6,
    notes = $7,
    is_active = $8
WHERE id = $1
`

type UpdatePetScheduleParams struct {
	ID               int64            `json:"id"`
	Title            pgtype.Text      `json:"title"`
	ReminderDatetime pgtype.Timestamp `json:"reminder_datetime"`
	EventRepeat      pgtype.Text      `json:"event_repeat"`
	EndType          pgtype.Bool      `json:"end_type"`
	EndDate          pgtype.Date      `json:"end_date"`
	Notes            pgtype.Text      `json:"notes"`
	IsActive         pgtype.Bool      `json:"is_active"`
}

func (q *Queries) UpdatePetSchedule(ctx context.Context, arg UpdatePetScheduleParams) error {
	_, err := q.db.Exec(ctx, updatePetSchedule,
		arg.ID,
		arg.Title,
		arg.ReminderDatetime,
		arg.EventRepeat,
		arg.EndType,
		arg.EndDate,
		arg.Notes,
		arg.IsActive,
	)
	return err
}
