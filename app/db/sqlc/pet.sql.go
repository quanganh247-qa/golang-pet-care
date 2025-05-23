// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: pet.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countPets = `-- name: CountPets :one
SELECT COUNT(*) FROM pets
WHERE is_active = true
`

func (q *Queries) CountPets(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countPets)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPet = `-- name: CreatePet :one
INSERT INTO pets (
    name,
    type,
    breed,
    age,
    gender,
    healthnotes,
    weight,
    birth_date,
    username,
    microchip_number,
    last_checkup_date,
    is_active,
    data_image,
    original_image
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, true, $12, $13
) RETURNING petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image
`

type CreatePetParams struct {
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	Breed           pgtype.Text   `json:"breed"`
	Age             pgtype.Int4   `json:"age"`
	Gender          pgtype.Text   `json:"gender"`
	Healthnotes     pgtype.Text   `json:"healthnotes"`
	Weight          pgtype.Float8 `json:"weight"`
	BirthDate       pgtype.Date   `json:"birth_date"`
	Username        string        `json:"username"`
	MicrochipNumber pgtype.Text   `json:"microchip_number"`
	LastCheckupDate pgtype.Date   `json:"last_checkup_date"`
	DataImage       []byte        `json:"data_image"`
	OriginalImage   pgtype.Text   `json:"original_image"`
}

func (q *Queries) CreatePet(ctx context.Context, arg CreatePetParams) (Pet, error) {
	row := q.db.QueryRow(ctx, createPet,
		arg.Name,
		arg.Type,
		arg.Breed,
		arg.Age,
		arg.Gender,
		arg.Healthnotes,
		arg.Weight,
		arg.BirthDate,
		arg.Username,
		arg.MicrochipNumber,
		arg.LastCheckupDate,
		arg.DataImage,
		arg.OriginalImage,
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
DELETE FROM pets WHERE petid = $1
`

func (q *Queries) DeletePet(ctx context.Context, petid int64) error {
	_, err := q.db.Exec(ctx, deletePet, petid)
	return err
}

const getAllPets = `-- name: GetAllPets :many
SELECT petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image FROM pets WHERE is_active is true
`

func (q *Queries) GetAllPets(ctx context.Context) ([]Pet, error) {
	rows, err := q.db.Query(ctx, getAllPets)
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

const getPetByID = `-- name: GetPetByID :one
SELECT petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image FROM pets 
WHERE petid = $1
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

const getPetDetailByUserID = `-- name: GetPetDetailByUserID :one
SELECT p.petid, p.name, p.type, p.breed, p.age, p.gender, p.healthnotes, p.weight, p.birth_date, p.username, p.microchip_number, p.last_checkup_date, p.is_active, p.data_image, p.original_image, u.full_name
FROM pets AS p
LEFT JOIN users AS u ON p.username = u.username
WHERE p.is_active = TRUE AND p.username = $1 AND p.name = $2
`

type GetPetDetailByUserIDParams struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type GetPetDetailByUserIDRow struct {
	Petid           int64         `json:"petid"`
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	Breed           pgtype.Text   `json:"breed"`
	Age             pgtype.Int4   `json:"age"`
	Gender          pgtype.Text   `json:"gender"`
	Healthnotes     pgtype.Text   `json:"healthnotes"`
	Weight          pgtype.Float8 `json:"weight"`
	BirthDate       pgtype.Date   `json:"birth_date"`
	Username        string        `json:"username"`
	MicrochipNumber pgtype.Text   `json:"microchip_number"`
	LastCheckupDate pgtype.Date   `json:"last_checkup_date"`
	IsActive        pgtype.Bool   `json:"is_active"`
	DataImage       []byte        `json:"data_image"`
	OriginalImage   pgtype.Text   `json:"original_image"`
	FullName        pgtype.Text   `json:"full_name"`
}

func (q *Queries) GetPetDetailByUserID(ctx context.Context, arg GetPetDetailByUserIDParams) (GetPetDetailByUserIDRow, error) {
	row := q.db.QueryRow(ctx, getPetDetailByUserID, arg.Username, arg.Name)
	var i GetPetDetailByUserIDRow
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
		&i.FullName,
	)
	return i, err
}

const getPetProfileSummary = `-- name: GetPetProfileSummary :many
SELECT p.petid, p.name, p.type, p.breed, p.age, p.gender, p.healthnotes, p.weight, p.birth_date, p.username, p.microchip_number, p.last_checkup_date, p.is_active, p.data_image, p.original_image, pt.id, pt.pet_id, pt.start_date, pt.end_date, pt.status, pt.description, pt.created_at, pt.doctor_id, pt.name, pt.type, pt.diseases, v.vaccinationid, v.petid, v.vaccinename, v.dateadministered, v.nextduedate, v.vaccineprovider, v.batchnumber, v.notes 
FROM pets AS p
LEFT JOIN pet_treatments AS pt ON p.petid = pt.pet_id
LEFT JOIN vaccinations AS v ON p.petid = v.petid
WHERE p.is_active = TRUE AND p.petid = $1
`

type GetPetProfileSummaryRow struct {
	Petid            int64              `json:"petid"`
	Name             string             `json:"name"`
	Type             string             `json:"type"`
	Breed            pgtype.Text        `json:"breed"`
	Age              pgtype.Int4        `json:"age"`
	Gender           pgtype.Text        `json:"gender"`
	Healthnotes      pgtype.Text        `json:"healthnotes"`
	Weight           pgtype.Float8      `json:"weight"`
	BirthDate        pgtype.Date        `json:"birth_date"`
	Username         string             `json:"username"`
	MicrochipNumber  pgtype.Text        `json:"microchip_number"`
	LastCheckupDate  pgtype.Date        `json:"last_checkup_date"`
	IsActive         pgtype.Bool        `json:"is_active"`
	DataImage        []byte             `json:"data_image"`
	OriginalImage    pgtype.Text        `json:"original_image"`
	ID               pgtype.Int8        `json:"id"`
	PetID            pgtype.Int8        `json:"pet_id"`
	StartDate        pgtype.Date        `json:"start_date"`
	EndDate          pgtype.Date        `json:"end_date"`
	Status           pgtype.Text        `json:"status"`
	Description      pgtype.Text        `json:"description"`
	CreatedAt        pgtype.Timestamptz `json:"created_at"`
	DoctorID         pgtype.Int4        `json:"doctor_id"`
	Name_2           pgtype.Text        `json:"name_2"`
	Type_2           pgtype.Text        `json:"type_2"`
	Diseases         pgtype.Text        `json:"diseases"`
	Vaccinationid    pgtype.Int8        `json:"vaccinationid"`
	Petid_2          pgtype.Int8        `json:"petid_2"`
	Vaccinename      pgtype.Text        `json:"vaccinename"`
	Dateadministered pgtype.Timestamp   `json:"dateadministered"`
	Nextduedate      pgtype.Timestamp   `json:"nextduedate"`
	Vaccineprovider  pgtype.Text        `json:"vaccineprovider"`
	Batchnumber      pgtype.Text        `json:"batchnumber"`
	Notes            pgtype.Text        `json:"notes"`
}

func (q *Queries) GetPetProfileSummary(ctx context.Context, petid int64) ([]GetPetProfileSummaryRow, error) {
	rows, err := q.db.Query(ctx, getPetProfileSummary, petid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPetProfileSummaryRow{}
	for rows.Next() {
		var i GetPetProfileSummaryRow
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
			&i.ID,
			&i.PetID,
			&i.StartDate,
			&i.EndDate,
			&i.Status,
			&i.Description,
			&i.CreatedAt,
			&i.DoctorID,
			&i.Name_2,
			&i.Type_2,
			&i.Diseases,
			&i.Vaccinationid,
			&i.Petid_2,
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

const listPets = `-- name: ListPets :many
SELECT petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image FROM pets
WHERE is_active = true 
ORDER BY name LIMIT $1 OFFSET $2
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
SELECT petid, name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, is_active, data_image, original_image FROM pets
WHERE username = $1 AND is_active = true
ORDER BY name LIMIT $2 OFFSET $3
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
UPDATE pets
SET is_active = false
WHERE petid = $1
`

func (q *Queries) SetPetInactive(ctx context.Context, petid int64) error {
	_, err := q.db.Exec(ctx, setPetInactive, petid)
	return err
}

const updatePet = `-- name: UpdatePet :exec
UPDATE pets
SET 
    name = $2,
    type = $3,
    breed = $4,
    age = $5,
    gender = $6,
    healthnotes = $7,
    weight = $8,
    birth_date = $9,
    microchip_number = $10,
    last_checkup_date = $11
WHERE petid = $1
`

type UpdatePetParams struct {
	Petid           int64         `json:"petid"`
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	Breed           pgtype.Text   `json:"breed"`
	Age             pgtype.Int4   `json:"age"`
	Gender          pgtype.Text   `json:"gender"`
	Healthnotes     pgtype.Text   `json:"healthnotes"`
	Weight          pgtype.Float8 `json:"weight"`
	BirthDate       pgtype.Date   `json:"birth_date"`
	MicrochipNumber pgtype.Text   `json:"microchip_number"`
	LastCheckupDate pgtype.Date   `json:"last_checkup_date"`
}

func (q *Queries) UpdatePet(ctx context.Context, arg UpdatePetParams) error {
	_, err := q.db.Exec(ctx, updatePet,
		arg.Petid,
		arg.Name,
		arg.Type,
		arg.Breed,
		arg.Age,
		arg.Gender,
		arg.Healthnotes,
		arg.Weight,
		arg.BirthDate,
		arg.MicrochipNumber,
		arg.LastCheckupDate,
	)
	return err
}

const updatePetAvatar = `-- name: UpdatePetAvatar :exec
UPDATE pets
SET 
    data_image = $2,
    original_image = $3
WHERE petid = $1
`

type UpdatePetAvatarParams struct {
	Petid         int64       `json:"petid"`
	DataImage     []byte      `json:"data_image"`
	OriginalImage pgtype.Text `json:"original_image"`
}

func (q *Queries) UpdatePetAvatar(ctx context.Context, arg UpdatePetAvatarParams) error {
	_, err := q.db.Exec(ctx, updatePetAvatar, arg.Petid, arg.DataImage, arg.OriginalImage)
	return err
}
