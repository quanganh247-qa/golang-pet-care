// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: vaccination.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createVaccination = `-- name: CreateVaccination :one
INSERT INTO Vaccination (petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING vaccinationID, petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes
`

type CreateVaccinationParams struct {
	Petid            pgtype.Int8      `json:"petid"`
	Vaccinename      string           `json:"vaccinename"`
	Dateadministered pgtype.Timestamp `json:"dateadministered"`
	Nextduedate      pgtype.Timestamp `json:"nextduedate"`
	Vaccineprovider  pgtype.Text      `json:"vaccineprovider"`
	Batchnumber      pgtype.Text      `json:"batchnumber"`
	Notes            pgtype.Text      `json:"notes"`
}

func (q *Queries) CreateVaccination(ctx context.Context, arg CreateVaccinationParams) (Vaccination, error) {
	row := q.db.QueryRow(ctx, createVaccination,
		arg.Petid,
		arg.Vaccinename,
		arg.Dateadministered,
		arg.Nextduedate,
		arg.Vaccineprovider,
		arg.Batchnumber,
		arg.Notes,
	)
	var i Vaccination
	err := row.Scan(
		&i.Vaccinationid,
		&i.Petid,
		&i.Vaccinename,
		&i.Dateadministered,
		&i.Nextduedate,
		&i.Vaccineprovider,
		&i.Batchnumber,
		&i.Notes,
	)
	return i, err
}

const deleteVaccination = `-- name: DeleteVaccination :exec
DELETE FROM Vaccination
WHERE vaccinationID = $1
`

func (q *Queries) DeleteVaccination(ctx context.Context, vaccinationid int64) error {
	_, err := q.db.Exec(ctx, deleteVaccination, vaccinationid)
	return err
}

const getVaccinationByID = `-- name: GetVaccinationByID :one
SELECT vaccinationID, petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes
FROM Vaccination
WHERE vaccinationID = $1
`

func (q *Queries) GetVaccinationByID(ctx context.Context, vaccinationid int64) (Vaccination, error) {
	row := q.db.QueryRow(ctx, getVaccinationByID, vaccinationid)
	var i Vaccination
	err := row.Scan(
		&i.Vaccinationid,
		&i.Petid,
		&i.Vaccinename,
		&i.Dateadministered,
		&i.Nextduedate,
		&i.Vaccineprovider,
		&i.Batchnumber,
		&i.Notes,
	)
	return i, err
}

const listVaccinationsByPetID = `-- name: ListVaccinationsByPetID :many
SELECT vaccinationID, petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes
FROM Vaccination
WHERE petID = $1
ORDER BY dateAdministered DESC LIMIT $2 OFFSET $3
`

type ListVaccinationsByPetIDParams struct {
	Petid  pgtype.Int8 `json:"petid"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

func (q *Queries) ListVaccinationsByPetID(ctx context.Context, arg ListVaccinationsByPetIDParams) ([]Vaccination, error) {
	rows, err := q.db.Query(ctx, listVaccinationsByPetID, arg.Petid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Vaccination{}
	for rows.Next() {
		var i Vaccination
		if err := rows.Scan(
			&i.Vaccinationid,
			&i.Petid,
			&i.Vaccinename,
			&i.Dateadministered,
			&i.Nextduedate,
			&i.Vaccineprovider,
			&i.Batchnumber,
			&i.Notes,
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

const updateVaccination = `-- name: UpdateVaccination :exec
UPDATE Vaccination
SET vaccineName = $2,
    dateAdministered = $3,
    nextDueDate = $4,
    vaccineProvider = $5,
    batchNumber = $6,
    notes = $7
WHERE vaccinationID = $1
`

type UpdateVaccinationParams struct {
	Vaccinationid    int64            `json:"vaccinationid"`
	Vaccinename      string           `json:"vaccinename"`
	Dateadministered pgtype.Timestamp `json:"dateadministered"`
	Nextduedate      pgtype.Timestamp `json:"nextduedate"`
	Vaccineprovider  pgtype.Text      `json:"vaccineprovider"`
	Batchnumber      pgtype.Text      `json:"batchnumber"`
	Notes            pgtype.Text      `json:"notes"`
}

func (q *Queries) UpdateVaccination(ctx context.Context, arg UpdateVaccinationParams) error {
	_, err := q.db.Exec(ctx, updateVaccination,
		arg.Vaccinationid,
		arg.Vaccinename,
		arg.Dateadministered,
		arg.Nextduedate,
		arg.Vaccineprovider,
		arg.Batchnumber,
		arg.Notes,
	)
	return err
}
