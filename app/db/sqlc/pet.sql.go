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
INSERT INTO Pet (username, Name, Type, Breed, Age, Weight, Gender, HealthNotes, ProfileImage)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING petid, name, type, breed, age, gender, healthnotes, profileimage, weight, username
`

type CreatePetParams struct {
	Username     string        `json:"username"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Breed        pgtype.Text   `json:"breed"`
	Age          pgtype.Int4   `json:"age"`
	Weight       pgtype.Float8 `json:"weight"`
	Gender       pgtype.Text   `json:"gender"`
	Healthnotes  pgtype.Text   `json:"healthnotes"`
	Profileimage pgtype.Text   `json:"profileimage"`
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
		arg.Profileimage,
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
		&i.Profileimage,
		&i.Weight,
		&i.Username,
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
SELECT petid, name, type, breed, age, gender, healthnotes, profileimage, weight, username FROM Pet WHERE PetID = $1
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
		&i.Profileimage,
		&i.Weight,
		&i.Username,
	)
	return i, err
}

const listPets = `-- name: ListPets :many
SELECT petid, name, type, breed, age, gender, healthnotes, profileimage, weight, username FROM Pet ORDER BY PetID LIMIT $1 OFFSET $2
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
			&i.Profileimage,
			&i.Weight,
			&i.Username,
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

const updatePet = `-- name: UpdatePet :exec
UPDATE Pet
SET Name = $2, Type = $3, Breed = $4, Age = $5, Weight = $6, Gender = $7, HealthNotes = $8, ProfileImage = $9
WHERE PetID = $1
`

type UpdatePetParams struct {
	Petid        int64         `json:"petid"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Breed        pgtype.Text   `json:"breed"`
	Age          pgtype.Int4   `json:"age"`
	Weight       pgtype.Float8 `json:"weight"`
	Gender       pgtype.Text   `json:"gender"`
	Healthnotes  pgtype.Text   `json:"healthnotes"`
	Profileimage pgtype.Text   `json:"profileimage"`
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
		arg.Profileimage,
	)
	return err
}
