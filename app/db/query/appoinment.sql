-- name: CreateAppointment :one
INSERT INTO Appointment
( petid, doctor_id, service_id, "date", status, notes, reminder_send, time_slot_id, created_at)
VALUES( 
    $1, $2, $3, $4, $5, $6, $7, $8, now()
) RETURNING *;

-- name: UpdateNotification :exec
UPDATE Appointment
SET reminder_send = true
WHERE appointment_id = $1;

-- name: UpdateAppointmentStatus :exec
UPDATE Appointment
SET status = $2
WHERE appointment_id = $1;

-- name: GetAppointmentsOfDoctorWithDetails :many
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
ORDER BY ts.start_time ASC;

-- name: GetAppointmentDetailById :one
SELECT * from Appointment WHERE appointment_id = $1;

-- name: GetAppointmentsByUser :many
SELECT 
    p.*, s.*, a.*, ts.*
FROM 
    appointment a
JOIN 
    pet p ON a.petid = p.petid 
JOIN 
    service s ON a.service_id = s.serviceid 
JOIN 
    timeslots ts ON a.time_slot_id = ts.id
WHERE 
    a.username = $1 and a.status <> 'completed';

-- name: GetAppointmentsByDoctor :one
SELECT COUNT(*) 
FROM appointment 
WHERE date = $1 AND doctor_id = $2 AND status = 'completed';