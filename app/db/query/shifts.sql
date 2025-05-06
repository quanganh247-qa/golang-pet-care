-- name: GetDoctors :many
SELECT id, specialization
FROM doctors;

-- name: GetTimeSlotsByDoctorAndDate :many
SELECT id, doctor_id, date, start_time, end_time, max_patients, booked_patients
FROM time_slots
WHERE doctor_id = $1 AND date = $2;

-- name: CountAppointmentsByDoctorAndDate :one
SELECT COUNT(*) as count
FROM appointments
WHERE doctor_id = $1 AND date >= $2 AND date < $3;

-- name: GetAppointmentsByTimeSlot :many
SELECT a.appointment_id, a.doctor_id, a.time_slot_id, a.service_id, s.category
FROM appointments a
JOIN services s ON a.service_id = s.id
WHERE a.time_slot_id = $1;

-- name: CountShiftsByDoctorAndDate :one
SELECT COUNT(*) as count
FROM shifts
WHERE doctor_id = $1 AND start_time >= $2 AND end_time <= $3;

-- name: DeleteShiftsByDate :exec
DELETE FROM shifts
WHERE start_time >= $1 AND end_time <= $2;

-- name: CreateShift :one
INSERT INTO shifts (doctor_id, start_time, end_time, assigned_patients)
VALUES ($1, $2, $3, $4)
RETURNING id, doctor_id, start_time, end_time, assigned_patients, created_at;

-- name: GetShifts :many
SELECT id, doctor_id, start_time, end_time, assigned_patients, created_at
FROM shifts
ORDER BY start_time;

-- name: GetShiftByDoctorId :many
SELECT id, doctor_id, start_time, end_time, assigned_patients, created_at
FROM shifts
WHERE doctor_id = $1;


-- name: DeleteShift :exec
DELETE FROM shifts
WHERE id = $1;