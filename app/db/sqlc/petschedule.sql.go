// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: petschedule.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPetSchedule = `-- name: CreatePetSchedule :exec
INSERT INTO pet_schedule (
   pet_id,
   title,
   reminder_datetime,
   event_repeat,
   end_type,
   end_date,
   notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreatePetScheduleParams struct {
	PetID            pgtype.Int8      `json:"pet_id"`
	Title            pgtype.Text      `json:"title"`
	ReminderDatetime pgtype.Timestamp `json:"reminder_datetime"`
	EventRepeat      pgtype.Text      `json:"event_repeat"`
	EndType          pgtype.Text      `json:"end_type"`
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

const getAllSchedulesByPet = `-- name: GetAllSchedulesByPet :many
SELECT id, pet_id, title, reminder_datetime, event_repeat, end_type, end_date, notes, created_at FROM pet_schedule 
WHERE pet_id = $1
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

const listPetSchedulesByUsername = `-- name: ListPetSchedulesByUsername :many
SELECT pet_schedule.id, pet_schedule.pet_id, pet_schedule.title, pet_schedule.reminder_datetime, pet_schedule.event_repeat, pet_schedule.end_type, pet_schedule.end_date, pet_schedule.notes, pet_schedule.created_at, pet.name
FROM pet_schedule
LEFT JOIN pet ON pet_schedule.pet_id = pet.petid
LEFT JOIN users ON pet.username = users.username
WHERE users.username = $1
ORDER BY pet.petid, pet_schedule.reminder_datetime
`

type ListPetSchedulesByUsernameRow struct {
	ID               int64            `json:"id"`
	PetID            pgtype.Int8      `json:"pet_id"`
	Title            pgtype.Text      `json:"title"`
	ReminderDatetime pgtype.Timestamp `json:"reminder_datetime"`
	EventRepeat      pgtype.Text      `json:"event_repeat"`
	EndType          pgtype.Text      `json:"end_type"`
	EndDate          pgtype.Date      `json:"end_date"`
	Notes            pgtype.Text      `json:"notes"`
	CreatedAt        pgtype.Timestamp `json:"created_at"`
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
