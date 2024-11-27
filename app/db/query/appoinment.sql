-- name: CreateAppointment :one
INSERT INTO Appointment (
    doctor_id,
    petid,
    service_id,
    time_slot_id,
    date,
    status
) VALUES (
    $1, $2, $3, $4, $5,'pending'
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

-- name: GetAppointmentsByPetOfUser :many
SELECT 
    a.appointment_id,
    a.petid,
    a.doctor_id,
    a.service_id,
    a.date,
    a.status,
    a.notes,
    a.reminder_send,
    a.time_slot_id,
    a.created_at
FROM 
    Appointment a
JOIN 
    Pet p ON a.petid = p.petid
JOIN 
    users u ON p.username = u.username
WHERE 
    u.username = $1 and p.is_active is true
ORDER BY 
    a.date DESC;