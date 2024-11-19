// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), $10)
RETURNING id
`

type CreateUserParams struct {
	Username        string      `json:"username"`
	HashedPassword  string      `json:"hashed_password"`
	FullName        string      `json:"full_name"`
	Email           string      `json:"email"`
	PhoneNumber     pgtype.Text `json:"phone_number"`
	Address         pgtype.Text `json:"address"`
	DataImage       []byte      `json:"data_image"`
	OriginalImage   pgtype.Text `json:"original_image"`
	Role            pgtype.Text `json:"role"`
	IsVerifiedEmail pgtype.Bool `json:"is_verified_email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int64, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
		arg.PhoneNumber,
		arg.Address,
		arg.DataImage,
		arg.OriginalImage,
		arg.Role,
		arg.IsVerifiedEmail,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getActiveDoctors = `-- name: GetActiveDoctors :many
SELECT 
  d.id,
  u.full_name AS name,
  d.specialization,
  d.years_of_experience,
  d.consultation_fee
FROM 
  Doctors d
JOIN 
  users u ON d.user_id = u.id
LEFT JOIN 
  DoctorSchedules ds ON d.id = ds.doctor_id
WHERE 
  d.is_active = true
  AND (ds.is_active = true OR ds.is_active IS NULL)
  AND ($1::VARCHAR IS NULL OR d.specialization = $1)
  AND ($2::INT IS NULL OR ds.day_of_week = $2)
ORDER BY 
  u.full_name
`

type GetActiveDoctorsParams struct {
	Column1 string `json:"column_1"`
	Column2 int32  `json:"column_2"`
}

type GetActiveDoctorsRow struct {
	ID                int64          `json:"id"`
	Name              string         `json:"name"`
	Specialization    pgtype.Text    `json:"specialization"`
	YearsOfExperience pgtype.Int4    `json:"years_of_experience"`
	ConsultationFee   pgtype.Numeric `json:"consultation_fee"`
}

func (q *Queries) GetActiveDoctors(ctx context.Context, arg GetActiveDoctorsParams) ([]GetActiveDoctorsRow, error) {
	rows, err := q.db.Query(ctx, getActiveDoctors, arg.Column1, arg.Column2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetActiveDoctorsRow{}
	for rows.Next() {
		var i GetActiveDoctorsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Specialization,
			&i.YearsOfExperience,
			&i.ConsultationFee,
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

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email, removed_at FROM users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.HashedPassword,
			&i.FullName,
			&i.Email,
			&i.PhoneNumber,
			&i.Address,
			&i.DataImage,
			&i.OriginalImage,
			&i.Role,
			&i.CreatedAt,
			&i.IsVerifiedEmail,
			&i.RemovedAt,
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

const getDoctor = `-- name: GetDoctor :one
SELECT 
  d.id,
  u.full_name AS name,
  d.specialization,
  d.years_of_experience,
  d.education,
  d.certificate_number,
  d.bio,
  d.consultation_fee
FROM
  Doctors d
JOIN
  users u ON d.user_id = u.id
WHERE
  d.id = $1
`

type GetDoctorRow struct {
	ID                int64          `json:"id"`
	Name              string         `json:"name"`
	Specialization    pgtype.Text    `json:"specialization"`
	YearsOfExperience pgtype.Int4    `json:"years_of_experience"`
	Education         pgtype.Text    `json:"education"`
	CertificateNumber pgtype.Text    `json:"certificate_number"`
	Bio               pgtype.Text    `json:"bio"`
	ConsultationFee   pgtype.Numeric `json:"consultation_fee"`
}

func (q *Queries) GetDoctor(ctx context.Context, id int64) (GetDoctorRow, error) {
	row := q.db.QueryRow(ctx, getDoctor, id)
	var i GetDoctorRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Specialization,
		&i.YearsOfExperience,
		&i.Education,
		&i.CertificateNumber,
		&i.Bio,
		&i.ConsultationFee,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email
FROM users
WHERE username = $1
`

type GetUserRow struct {
	ID              int64            `json:"id"`
	Username        string           `json:"username"`
	HashedPassword  string           `json:"hashed_password"`
	FullName        string           `json:"full_name"`
	Email           string           `json:"email"`
	PhoneNumber     pgtype.Text      `json:"phone_number"`
	Address         pgtype.Text      `json:"address"`
	DataImage       []byte           `json:"data_image"`
	OriginalImage   pgtype.Text      `json:"original_image"`
	Role            pgtype.Text      `json:"role"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	IsVerifiedEmail pgtype.Bool      `json:"is_verified_email"`
}

func (q *Queries) GetUser(ctx context.Context, username string) (GetUserRow, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i GetUserRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PhoneNumber,
		&i.Address,
		&i.DataImage,
		&i.OriginalImage,
		&i.Role,
		&i.CreatedAt,
		&i.IsVerifiedEmail,
	)
	return i, err
}

const insertDoctor = `-- name: InsertDoctor :one
INSERT INTO Doctors (
    user_id,
    specialization,
    years_of_experience,
    education,
    certificate_number,
    bio,
    consultation_fee
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, user_id, specialization, years_of_experience, education, certificate_number, bio, consultation_fee
`

type InsertDoctorParams struct {
	UserID            int64          `json:"user_id"`
	Specialization    pgtype.Text    `json:"specialization"`
	YearsOfExperience pgtype.Int4    `json:"years_of_experience"`
	Education         pgtype.Text    `json:"education"`
	CertificateNumber pgtype.Text    `json:"certificate_number"`
	Bio               pgtype.Text    `json:"bio"`
	ConsultationFee   pgtype.Numeric `json:"consultation_fee"`
}

func (q *Queries) InsertDoctor(ctx context.Context, arg InsertDoctorParams) (Doctor, error) {
	row := q.db.QueryRow(ctx, insertDoctor,
		arg.UserID,
		arg.Specialization,
		arg.YearsOfExperience,
		arg.Education,
		arg.CertificateNumber,
		arg.Bio,
		arg.ConsultationFee,
	)
	var i Doctor
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Specialization,
		&i.YearsOfExperience,
		&i.Education,
		&i.CertificateNumber,
		&i.Bio,
		&i.ConsultationFee,
	)
	return i, err
}

const insertDoctorSchedule = `-- name: InsertDoctorSchedule :one
INSERT INTO DoctorSchedules (
    doctor_id,
    day_of_week,
    start_time,
    end_time,
    is_active,
    max_appointments
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, doctor_id, day_of_week, start_time, end_time, is_active, max_appointments
`

type InsertDoctorScheduleParams struct {
	DoctorID        int64            `json:"doctor_id"`
	DayOfWeek       pgtype.Int4      `json:"day_of_week"`
	StartTime       pgtype.Timestamp `json:"start_time"`
	EndTime         pgtype.Timestamp `json:"end_time"`
	IsActive        pgtype.Bool      `json:"is_active"`
	MaxAppointments pgtype.Int4      `json:"max_appointments"`
}

func (q *Queries) InsertDoctorSchedule(ctx context.Context, arg InsertDoctorScheduleParams) (Doctorschedule, error) {
	row := q.db.QueryRow(ctx, insertDoctorSchedule,
		arg.DoctorID,
		arg.DayOfWeek,
		arg.StartTime,
		arg.EndTime,
		arg.IsActive,
		arg.MaxAppointments,
	)
	var i Doctorschedule
	err := row.Scan(
		&i.ID,
		&i.DoctorID,
		&i.DayOfWeek,
		&i.StartTime,
		&i.EndTime,
		&i.IsActive,
		&i.MaxAppointments,
	)
	return i, err
}

const verifiedUser = `-- name: VerifiedUser :one
UPDATE users
SET is_verified_email = $2
WHERE username = $1
RETURNING id, username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email, removed_at
`

type VerifiedUserParams struct {
	Username        string      `json:"username"`
	IsVerifiedEmail pgtype.Bool `json:"is_verified_email"`
}

func (q *Queries) VerifiedUser(ctx context.Context, arg VerifiedUserParams) (User, error) {
	row := q.db.QueryRow(ctx, verifiedUser, arg.Username, arg.IsVerifiedEmail)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PhoneNumber,
		&i.Address,
		&i.DataImage,
		&i.OriginalImage,
		&i.Role,
		&i.CreatedAt,
		&i.IsVerifiedEmail,
		&i.RemovedAt,
	)
	return i, err
}
