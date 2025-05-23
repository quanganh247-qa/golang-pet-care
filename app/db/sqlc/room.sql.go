// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: room.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const assignRoomToAppointment = `-- name: AssignRoomToAppointment :exec
UPDATE rooms 
SET current_appointment_id = $2
WHERE id = $1
`

type AssignRoomToAppointmentParams struct {
	ID                   int64       `json:"id"`
	CurrentAppointmentID pgtype.Int8 `json:"current_appointment_id"`
}

func (q *Queries) AssignRoomToAppointment(ctx context.Context, arg AssignRoomToAppointmentParams) error {
	_, err := q.db.Exec(ctx, assignRoomToAppointment, arg.ID, arg.CurrentAppointmentID)
	return err
}

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (name, type, status, current_appointment_id, available_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, type, status, current_appointment_id, available_at
`

type CreateRoomParams struct {
	Name                 string           `json:"name"`
	Type                 string           `json:"type"`
	Status               pgtype.Text      `json:"status"`
	CurrentAppointmentID pgtype.Int8      `json:"current_appointment_id"`
	AvailableAt          pgtype.Timestamp `json:"available_at"`
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error) {
	row := q.db.QueryRow(ctx, createRoom,
		arg.Name,
		arg.Type,
		arg.Status,
		arg.CurrentAppointmentID,
		arg.AvailableAt,
	)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Status,
		&i.CurrentAppointmentID,
		&i.AvailableAt,
	)
	return i, err
}

const deleteRoom = `-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = $1
`

func (q *Queries) DeleteRoom(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteRoom, id)
	return err
}

const getAvailableRooms = `-- name: GetAvailableRooms :many
SELECT id, name, type, status, current_appointment_id, available_at
FROM rooms
WHERE status = 'available' 
LIMIT $1 OFFSET $2
`

type GetAvailableRoomsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAvailableRooms(ctx context.Context, arg GetAvailableRoomsParams) ([]Room, error) {
	rows, err := q.db.Query(ctx, getAvailableRooms, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Room{}
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Status,
			&i.CurrentAppointmentID,
			&i.AvailableAt,
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

const getRoomByID = `-- name: GetRoomByID :one
SELECT id, name, type, status, current_appointment_id, available_at FROM rooms WHERE id = $1
`

func (q *Queries) GetRoomByID(ctx context.Context, id int64) (Room, error) {
	row := q.db.QueryRow(ctx, getRoomByID, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Type,
		&i.Status,
		&i.CurrentAppointmentID,
		&i.AvailableAt,
	)
	return i, err
}

const releaseRoom = `-- name: ReleaseRoom :exec
UPDATE rooms
SET status = 'available',
    current_appointment_id = NULL,
    available_at = $1
WHERE id = $2
`

type ReleaseRoomParams struct {
	AvailableAt pgtype.Timestamp `json:"available_at"`
	ID          int64            `json:"id"`
}

func (q *Queries) ReleaseRoom(ctx context.Context, arg ReleaseRoomParams) error {
	_, err := q.db.Exec(ctx, releaseRoom, arg.AvailableAt, arg.ID)
	return err
}

const updateRoom = `-- name: UpdateRoom :exec
UPDATE rooms
SET name = $2,
    type = $3,
    status = $4,
    current_appointment_id = $5,
    available_at = $6
WHERE id = $1
`

type UpdateRoomParams struct {
	ID                   int64            `json:"id"`
	Name                 string           `json:"name"`
	Type                 string           `json:"type"`
	Status               pgtype.Text      `json:"status"`
	CurrentAppointmentID pgtype.Int8      `json:"current_appointment_id"`
	AvailableAt          pgtype.Timestamp `json:"available_at"`
}

func (q *Queries) UpdateRoom(ctx context.Context, arg UpdateRoomParams) error {
	_, err := q.db.Exec(ctx, updateRoom,
		arg.ID,
		arg.Name,
		arg.Type,
		arg.Status,
		arg.CurrentAppointmentID,
		arg.AvailableAt,
	)
	return err
}
