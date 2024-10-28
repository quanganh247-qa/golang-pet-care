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
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING id, username, hashed_password, full_name, email, phone_number, address, avatar, role, created_at, is_verified_email, removed_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PhoneNumber,
		&i.Address,
		&i.Avatar,
		&i.Role,
		&i.CreatedAt,
		&i.IsVerifiedEmail,
		&i.RemovedAt,
	)
	return i, err
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
SELECT id, username, hashed_password, full_name, email, phone_number, address, avatar, role, created_at, is_verified_email, removed_at FROM users
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
			&i.Avatar,
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
SELECT id, username, hashed_password, full_name, email, phone_number, address, avatar, role, created_at, is_verified_email, removed_at FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PhoneNumber,
		&i.Address,
		&i.Avatar,
		&i.Role,
		&i.CreatedAt,
		&i.IsVerifiedEmail,
		&i.RemovedAt,
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
	DoctorID        int64       `json:"doctor_id"`
	DayOfWeek       pgtype.Int4 `json:"day_of_week"`
	StartTime       pgtype.Time `json:"start_time"`
	EndTime         pgtype.Time `json:"end_time"`
	IsActive        pgtype.Bool `json:"is_active"`
	MaxAppointments pgtype.Int4 `json:"max_appointments"`
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

const updateUser = `-- name: UpdateUser :one
UPDATE users 
SET 
  hashed_password = COALESCE($1, hashed_password),
  full_name =  COALESCE($2,full_name),
  email = COALESCE($3,email),
  is_verified_email = COALESCE($4,is_verified_email)
WHERE
  username = $5
RETURNING id, username, hashed_password, full_name, email, phone_number, address, avatar, role, created_at, is_verified_email, removed_at
`

type UpdateUserParams struct {
	HashedPassword  pgtype.Text `json:"hashed_password"`
	FullName        pgtype.Text `json:"full_name"`
	Email           pgtype.Text `json:"email"`
	IsVerifiedEmail pgtype.Bool `json:"is_verified_email"`
	Username        string      `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
		arg.IsVerifiedEmail,
		arg.Username,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PhoneNumber,
		&i.Address,
		&i.Avatar,
		&i.Role,
		&i.CreatedAt,
		&i.IsVerifiedEmail,
		&i.RemovedAt,
	)
	return i, err
}
