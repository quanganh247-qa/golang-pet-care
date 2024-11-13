// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: pet.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPet = `-- name: CreatePet :one
INSERT INTO Pet (username, Name, Type, Breed, Age, Weight, Gender, HealthNotes, data_image, original_image, birth_date, microchip_number, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, true)
RETURNING petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image
`

type CreatePetParams struct {
	Username        string        `json:"username"`
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	Breed           pgtype.Text   `json:"breed"`
	Age             pgtype.Int4   `json:"age"`
	Weight          pgtype.Float8 `json:"weight"`
	Gender          pgtype.Text   `json:"gender"`
	Healthnotes     pgtype.Text   `json:"healthnotes"`
	DataImage       []byte        `json:"data_image"`
	OriginalImage   string        `json:"original_image"`
	BirthDate       pgtype.Date   `json:"birth_date"`
	MicrochipNumber pgtype.Text   `json:"microchip_number"`
}

func (q *Queries) CreatePet(ctx context.Context, arg CreatePetParams) (Pet, error) {
	row := q.db.QueryRow(ctx, createPet,
		arg.Username,
		arg.Name,
		arg.Type,
		arg.Breed,
		arg.Age,
		arg.Weight,
		arg.Gender,
		arg.Healthnotes,
		arg.DataImage,
		arg.OriginalImage,
		arg.BirthDate,
		arg.MicrochipNumber,
	)
	var i Pet
	err := row.Scan(
		&i.Petid,
		&i.Name,
		&i.Type,
		&i.Breed,
		&i.Age,
		&i.Gender,
		&i.Healthnotes,
		&i.Weight,
		&i.BirthDate,
		&i.Username,
		&i.MicrochipNumber,
		&i.LastCheckupDate,
		&i.IsActive,
		&i.DataImage,
		&i.OriginalImage,
	)
	return i, err
}

const deletePet = `-- name: DeletePet :exec
DELETE FROM Pet WHERE PetID = $1
`

func (q *Queries) DeletePet(ctx context.Context, petid int64) error {
	_, err := q.db.Exec(ctx, deletePet, petid)
	return err
}

const getPetByID = `-- name: GetPetByID :one
SELECT petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image FROM Pet WHERE PetID = $1
`

func (q *Queries) GetPetByID(ctx context.Context, petid int64) (Pet, error) {
	row := q.db.QueryRow(ctx, getPetByID, petid)
	var i Pet
	err := row.Scan(
		&i.Petid,
		&i.Name,
		&i.Type,
		&i.Breed,
		&i.Age,
		&i.Gender,
		&i.Healthnotes,
		&i.Weight,
		&i.BirthDate,
		&i.Username,
		&i.MicrochipNumber,
		&i.LastCheckupDate,
		&i.IsActive,
		&i.DataImage,
		&i.OriginalImage,
	)
	return i, err
}

const listPets = `-- name: ListPets :many
SELECT petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image FROM Pet ORDER BY PetID LIMIT $1 OFFSET $2
`

type ListPetsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPets(ctx context.Context, arg ListPetsParams) ([]Pet, error) {
	rows, err := q.db.Query(ctx, listPets, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Pet{}
	for rows.Next() {
		var i Pet
		if err := rows.Scan(
			&i.Petid,
			&i.Name,
			&i.Type,
			&i.Breed,
			&i.Age,
			&i.Gender,
			&i.Healthnotes,
			&i.Weight,
			&i.BirthDate,
			&i.Username,
			&i.MicrochipNumber,
			&i.LastCheckupDate,
			&i.IsActive,
			&i.DataImage,
			&i.OriginalImage,
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

const listPetsByUsername = `-- name: ListPetsByUsername :many
SELECT petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image FROM Pet WHERE username = $1 ORDER BY PetID LIMIT $2 OFFSET $3
`

type ListPetsByUsernameParams struct {
	Username string `json:"username"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) ListPetsByUsername(ctx context.Context, arg ListPetsByUsernameParams) ([]Pet, error) {
	rows, err := q.db.Query(ctx, listPetsByUsername, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Pet{}
	for rows.Next() {
		var i Pet
		if err := rows.Scan(
			&i.Petid,
			&i.Name,
			&i.Type,
			&i.Breed,
			&i.Age,
			&i.Gender,
			&i.Healthnotes,
			&i.Weight,
			&i.BirthDate,
			&i.Username,
			&i.MicrochipNumber,
			&i.LastCheckupDate,
			&i.IsActive,
			&i.DataImage,
			&i.OriginalImage,
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

const setPetInactive = `-- name: SetPetInactive :exec
UPDATE Pet SET is_active = $2 WHERE PetID = $1
`

type SetPetInactiveParams struct {
	Petid    int64       `json:"petid"`
	IsActive pgtype.Bool `json:"is_active"`
}

func (q *Queries) SetPetInactive(ctx context.Context, arg SetPetInactiveParams) error {
	_, err := q.db.Exec(ctx, setPetInactive, arg.Petid, arg.IsActive)
	return err
}

const updatePet = `-- name: UpdatePet :exec
UPDATE Pet
SET Name = $2, Type = $3, Breed = $4, Age = $5, Weight = $6, Gender = $7, HealthNotes = $8, data_image = $9, is_active = $10
WHERE PetID = $1
`

type UpdatePetParams struct {
	Petid       int64         `json:"petid"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Breed       pgtype.Text   `json:"breed"`
	Age         pgtype.Int4   `json:"age"`
	Weight      pgtype.Float8 `json:"weight"`
	Gender      pgtype.Text   `json:"gender"`
	Healthnotes pgtype.Text   `json:"healthnotes"`
	DataImage   []byte        `json:"data_image"`
	IsActive    pgtype.Bool   `json:"is_active"`
}

func (q *Queries) UpdatePet(ctx context.Context, arg UpdatePetParams) error {
	_, err := q.db.Exec(ctx, updatePet,
		arg.Petid,
		arg.Name,
		arg.Type,
		arg.Breed,
		arg.Age,
		arg.Weight,
		arg.Gender,
		arg.Healthnotes,
		arg.DataImage,
		arg.IsActive,
	)
	return err
}
