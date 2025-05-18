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
WHERE doctor_id = $1;

-- name: DeleteShiftsByDate :exec
DELETE FROM shifts
WHERE doctor_id = $1;

-- name: CreateShift :one
INSERT INTO shifts (doctor_id,date)
VALUES ($1, $2)
RETURNING id, doctor_id, date, created_at;

-- name: GetShifts :many
SELECT id, doctor_id, date, created_at
FROM shifts
ORDER BY created_at;

-- name: GetShiftByDoctorId :many
SELECT s.id, s.doctor_id, s.date, s.created_at, t.start_time, t.end_time
FROM shifts s
LEFT JOIN time_slots t ON s.id = t.shift_id
WHERE s.doctor_id = $1;


-- name: DeleteShift :exec
DELETE FROM shifts
WHERE id = $1;