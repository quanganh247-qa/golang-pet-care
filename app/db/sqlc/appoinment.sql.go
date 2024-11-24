// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: appoinment.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAppointment = `-- name: CreateAppointment :one
INSERT INTO Appointment (
    doctor_id,
    petid,
    service_id,
    time_slot_id,
    date,
    status
) VALUES (
    $1, $2, $3, $4, $5,'pending'
) RETURNING appointment_id, petid, doctor_id, service_id, date, status, notes, reminder_send, time_slot_id, created_at
`

type CreateAppointmentParams struct {
	DoctorID   pgtype.Int8      `json:"doctor_id"`
	Petid      pgtype.Int8      `json:"petid"`
	ServiceID  pgtype.Int8      `json:"service_id"`
	TimeSlotID pgtype.Int8      `json:"time_slot_id"`
	Date       pgtype.Timestamp `json:"date"`
}

func (q *Queries) CreateAppointment(ctx context.Context, arg CreateAppointmentParams) (Appointment, error) {
	row := q.db.QueryRow(ctx, createAppointment,
		arg.DoctorID,
		arg.Petid,
		arg.ServiceID,
		arg.TimeSlotID,
		arg.Date,
	)
	var i Appointment
	err := row.Scan(
		&i.AppointmentID,
		&i.Petid,
		&i.DoctorID,
		&i.ServiceID,
		&i.Date,
		&i.Status,
		&i.Notes,
		&i.ReminderSend,
		&i.TimeSlotID,
		&i.CreatedAt,
	)
	return i, err
}

const getAppointmentDetailById = `-- name: GetAppointmentDetailById :one
SELECT appointment_id, petid, doctor_id, service_id, date, status, notes, reminder_send, time_slot_id, created_at from Appointment WHERE appointment_id = $1
`

func (q *Queries) GetAppointmentDetailById(ctx context.Context, appointmentID int64) (Appointment, error) {
	row := q.db.QueryRow(ctx, getAppointmentDetailById, appointmentID)
	var i Appointment
	err := row.Scan(
		&i.AppointmentID,
		&i.Petid,
		&i.DoctorID,
		&i.ServiceID,
		&i.Date,
		&i.Status,
		&i.Notes,
		&i.ReminderSend,
		&i.TimeSlotID,
		&i.CreatedAt,
	)
	return i, err
}

const getAppointmentsOfDoctorWithDetails = `-- name: GetAppointmentsOfDoctorWithDetails :many
SELECT 
    a.appointment_id as appointment_id,
    p.name as pet_name,
    s.name as service_name,
    ts.start_time,
    ts.end_time
FROM Appointment a
    LEFT JOIN Doctors d ON a.doctor_id = d.id
    LEFT JOIN Pet p ON a.petid = p.petid
    LEFT JOIN Service s ON a.service_id = s.serviceid
    LEFT JOIN TimeSlots ts ON a.time_slot_id = ts.id
WHERE d.id = $1
AND LOWER(a.status) <> 'completed'
ORDER BY ts.start_time ASC
`

type GetAppointmentsOfDoctorWithDetailsRow struct {
	AppointmentID int64            `json:"appointment_id"`
	PetName       pgtype.Text      `json:"pet_name"`
	ServiceName   pgtype.Text      `json:"service_name"`
	StartTime     pgtype.Timestamp `json:"start_time"`
	EndTime       pgtype.Timestamp `json:"end_time"`
}

func (q *Queries) GetAppointmentsOfDoctorWithDetails(ctx context.Context, id int64) ([]GetAppointmentsOfDoctorWithDetailsRow, error) {
	rows, err := q.db.Query(ctx, getAppointmentsOfDoctorWithDetails, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAppointmentsOfDoctorWithDetailsRow{}
	for rows.Next() {
		var i GetAppointmentsOfDoctorWithDetailsRow
		if err := rows.Scan(
			&i.AppointmentID,
			&i.PetName,
			&i.ServiceName,
			&i.StartTime,
			&i.EndTime,
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

const updateAppointmentStatus = `-- name: UpdateAppointmentStatus :exec
UPDATE Appointment
SET status = $2
WHERE appointment_id = $1
`

type UpdateAppointmentStatusParams struct {
	AppointmentID int64       `json:"appointment_id"`
	Status        pgtype.Text `json:"status"`
}

func (q *Queries) UpdateAppointmentStatus(ctx context.Context, arg UpdateAppointmentStatusParams) error {
	_, err := q.db.Exec(ctx, updateAppointmentStatus, arg.AppointmentID, arg.Status)
	return err
}

const updateNotification = `-- name: UpdateNotification :exec
UPDATE Appointment
SET reminder_send = true
WHERE appointment_id = $1
`

func (q *Queries) UpdateNotification(ctx context.Context, appointmentID int64) error {
	_, err := q.db.Exec(ctx, updateNotification, appointmentID)
	return err
}
