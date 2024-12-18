// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: pet_logs.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const deletePetLog = `-- name: DeletePetLog :exec
DELETE FROM pet_logs
WHERE  log_id = $1
`

func (q *Queries) DeletePetLog(ctx context.Context, logID int64) error {
	_, err := q.db.Exec(ctx, deletePetLog, logID)
	return err
}

const getPetLogByID = `-- name: GetPetLogByID :one
SELECT pet_logs.petid, pet_logs.datetime, pet_logs.title, pet_logs.notes
FROM pet_logs
LEFT JOIN pet ON pet_logs.petid = pet.petid
WHERE pet_logs.petid = $1 AND pet_logs.log_id = $2 AND pet.is_active = true
`

type GetPetLogByIDParams struct {
	Petid int64 `json:"petid"`
	LogID int64 `json:"log_id"`
}

type GetPetLogByIDRow struct {
	Petid    int64            `json:"petid"`
	Datetime pgtype.Timestamp `json:"datetime"`
	Title    pgtype.Text      `json:"title"`
	Notes    pgtype.Text      `json:"notes"`
}

func (q *Queries) GetPetLogByID(ctx context.Context, arg GetPetLogByIDParams) (GetPetLogByIDRow, error) {
	row := q.db.QueryRow(ctx, getPetLogByID, arg.Petid, arg.LogID)
	var i GetPetLogByIDRow
	err := row.Scan(
		&i.Petid,
		&i.Datetime,
		&i.Title,
		&i.Notes,
	)
	return i, err
}

const getPetLogsByPetID = `-- name: GetPetLogsByPetID :many
SELECT pet_logs.petid, pet_logs.datetime, pet_logs.title, pet_logs.notes, pet_logs.log_id
FROM pet_logs
LEFT JOIN pet ON pet_logs.petid = pet.petid
WHERE pet_logs.petid = $1 AND pet.is_active = true
ORDER BY pet_logs.datetime DESC
LIMIT $2 OFFSET $3
`

type GetPetLogsByPetIDParams struct {
	Petid  int64 `json:"petid"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetPetLogsByPetIDRow struct {
	Petid    int64            `json:"petid"`
	Datetime pgtype.Timestamp `json:"datetime"`
	Title    pgtype.Text      `json:"title"`
	Notes    pgtype.Text      `json:"notes"`
	LogID    int64            `json:"log_id"`
}

func (q *Queries) GetPetLogsByPetID(ctx context.Context, arg GetPetLogsByPetIDParams) ([]GetPetLogsByPetIDRow, error) {
	rows, err := q.db.Query(ctx, getPetLogsByPetID, arg.Petid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPetLogsByPetIDRow{}
	for rows.Next() {
		var i GetPetLogsByPetIDRow
		if err := rows.Scan(
			&i.Petid,
			&i.Datetime,
			&i.Title,
			&i.Notes,
			&i.LogID,
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

const insertPetLog = `-- name: InsertPetLog :one
INSERT INTO pet_logs (petid, datetime, title, notes)
VALUES ($1, $2, $3, $4) RETURNING log_id, petid, datetime, title, notes
`

type InsertPetLogParams struct {
	Petid    int64            `json:"petid"`
	Datetime pgtype.Timestamp `json:"datetime"`
	Title    pgtype.Text      `json:"title"`
	Notes    pgtype.Text      `json:"notes"`
}

func (q *Queries) InsertPetLog(ctx context.Context, arg InsertPetLogParams) (PetLog, error) {
	row := q.db.QueryRow(ctx, insertPetLog,
		arg.Petid,
		arg.Datetime,
		arg.Title,
		arg.Notes,
	)
	var i PetLog
	err := row.Scan(
		&i.LogID,
		&i.Petid,
		&i.Datetime,
		&i.Title,
		&i.Notes,
	)
	return i, err
}

const updatePetLog = `-- name: UpdatePetLog :exec
UPDATE pet_logs
SET title = $2, notes = $3
WHERE log_id = $1
`

type UpdatePetLogParams struct {
	LogID int64       `json:"log_id"`
	Title pgtype.Text `json:"title"`
	Notes pgtype.Text `json:"notes"`
}

func (q *Queries) UpdatePetLog(ctx context.Context, arg UpdatePetLogParams) error {
	_, err := q.db.Exec(ctx, updatePetLog, arg.LogID, arg.Title, arg.Notes)
	return err
}
