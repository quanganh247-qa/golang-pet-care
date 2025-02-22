-- name: CreateAppointment :one
INSERT INTO appointments
( petid, doctor_id, username, service_id, "date", payment_status, notes, reminder_send, time_slot_id, created_at)
VALUES( 
    $1, $2, $3, $4, $5, $6, $7, $8, $9, now()
) RETURNING *;


-- name: UpdateTimeSlotBookedPatients :exec
UPDATE time_slots
SET booked_patients = booked_patients + 1
WHERE id = $1 AND doctor_id = $2;


-- name: UpdateNotification :exec
UPDATE appointments
SET reminder_send = true
WHERE appointment_id = $1;

-- name: UpdateAppointmentStatus :exec
UPDATE appointments
SET payment_status = $2
WHERE appointment_id = $1;

-- name: GetAppointmentsOfDoctorWithDetails :many
SELECT 
    a.appointment_id as appointment_id,
    p.name as pet_name,
    s.name as service_name,
    ts.start_time,
    ts.end_time
FROM appointments a
    LEFT JOIN doctors d ON a.doctor_id = d.id
    LEFT JOIN pets p ON a.petid = p.petid
    LEFT JOIN services s ON a.service_id = s.id
    LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
WHERE d.id = $1
AND LOWER(a.status) <> 'completed'
ORDER BY ts.start_time ASC;

-- name: GetAppointmentDetailById :one
SELECT * from appointments WHERE appointment_id = $1;

-- name: GetAppointmentsByUser :many
SELECT 
    p.*, s.*, a.*, ts.*
FROM 
    appointments a
JOIN 
    pets p ON a.petid = p.petid 
JOIN 
    services s ON a.service_id = s.id 
JOIN 
    time_slots ts ON a.time_slot_id = ts.id
WHERE 
    a.username = $1 and a.status <> 'completed';

-- name: CountAppointmentsByDateAndTimeSlot :one
SELECT COUNT(*) 
FROM appointments 
WHERE date = $1 AND doctor_id = $2 AND status = 'completed';

-- name: GetAppointmentsByDoctor :many
SELECT 
    a.*,
    d.id AS doctor_id,
    p.name AS pet_name,
    s.name AS service_name,
    ts.start_time,
    ts.end_time
FROM 
    appointments a
JOIN 
    doctors d ON a.doctor_id = d.id
JOIN 
    pets p ON a.petid = p.petid
JOIN 
    services as s ON a.service_id = s.id
JOIN 
    time_slots ts ON a.time_slot_id = ts.id
WHERE 
    a.doctor_id = $1;

-- name: ListAllAppointments :many
SELECT * FROM appointments;
