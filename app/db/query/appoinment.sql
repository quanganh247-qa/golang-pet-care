
-- name: CreateAppointment :one
INSERT INTO appointments (petid, username, doctor_id, service_id, date, time_slot_id, state_id)
VALUES ($1, $2, $3, $4, $5, $6, (SELECT id FROM states WHERE state = 'Scheduled'))
RETURNING *;

-- name: UpdateTimeSlotBookedPatients :exec
UPDATE time_slots
SET booked_patients = booked_patients + 1
WHERE id = $1;

-- name: UpdateAppointmentByID :exec
UPDATE appointments SET 
    state_id = $2,
    notes = $3,
    reminder_send = $4,
    updated_at = now()
WHERE appointment_id = $1;

-- name: UpdateNotification :exec
UPDATE appointments
SET reminder_send = true
WHERE appointment_id = $1;

-- name: UpdateAppointmentStatus :exec
UPDATE appointments
SET state_id = $2
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


-- name: CountAppointmentsByDateAndTimeSlot :one
SELECT COUNT(*) 
FROM appointments 
WHERE date = $1 AND doctor_id = $2 AND status = 'completed';

-- name: GetAppointmentsByDoctor :many
SELECT 
    a.appointment_id,
    a.date,
    a.created_at,
    a.notes,
    a.reminder_send,
    d.id AS doctor_id,
    p.name AS pet_name,
    s.name AS service_name,
    ts.start_time,
    ts.end_time,
    ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id
FROM 
    appointments a
LEFT JOIN 
    doctors d ON a.doctor_id = d.id
LEFT JOIN 
    pets p ON a.petid = p.petid
LEFT JOIN 
    services s ON a.service_id = s.id
LEFT JOIN 
    time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN 
    states st ON a.state_id = st.id
WHERE 
    a.doctor_id = $1
ORDER BY a.created_at DESC;

-- name: ListAllAppointments :many
SELECT * FROM appointments;

-- name: GetAppointmentByStateId :many
SELECT * FROM appointments WHERE state_id = $1;

-- name: GetAllAppointments :many
SELECT 
    a.appointment_id, a.date, a.notes, a.reminder_send, a.created_at,
    p.name AS pet_name,
    d.id AS doctor_id,
    s.name AS service_name,
    ts.start_time, ts.end_time, ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id
FROM appointments a
LEFT JOIN pets p ON a.petid = p.petid
LEFT JOIN services s ON a.service_id = s.id
LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN doctors d ON a.doctor_id = d.id
LEFT JOIN states st ON a.state_id = st.id;

-- name: GetSOAPByAppointmentID :one
SELECT * FROM consultations WHERE appointment_id = $1;

-- name: GetAppointmentDetail :one
SELECT 
    s.name AS service_name,
    p.name AS pet_name,
    st.state AS state_name
FROM services s, pets p, states st
WHERE s.id = $1 AND p.petid = $2 AND st.id = $3;

-- name: GetAppointmentDetailByAppointmentID :one
SELECT 
    a.appointment_id, a.date, a.notes, a.reminder_send, a.created_at,
    d.id AS doctor_id,
    p.name AS pet_name,
    s.name AS service_name,
    ts.start_time, ts.end_time, ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id
FROM appointments a
LEFT JOIN pets p ON p.petid = a.petid
LEFT JOIN doctors d ON d.id = a.doctor_id
LEFT JOIN services s ON s.id = a.service_id
LEFT JOIN time_slots ts ON ts.id = a.time_slot_id
LEFT JOIN states st ON st.id = a.state_id
WHERE a.appointment_id = $1;

-- name: GetAppointmentsByUser :many
SELECT 
    a.appointment_id, a.date, a.created_at,
    p.name AS pet_name,
    d.id AS doctor_id,
    s.name AS service_name,
    ts.start_time, ts.end_time,
    st.state
FROM appointments a
LEFT JOIN pets p ON p.petid = a.petid
LEFT JOIN doctors d ON d.id = a.doctor_id
LEFT JOIN services s ON s.id = a.service_id
LEFT JOIN time_slots ts ON ts.id = a.time_slot_id
LEFT JOIN states st ON st.id = a.state_id
WHERE a.username = $1;