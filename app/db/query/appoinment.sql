-- name: CreateAppointment :one
INSERT INTO Appointment (
    doctor_id,
    petid,
    service_id,
    time_slot_id,
    status
) VALUES (
    $1, $2, $3, $4, 'pending'
) RETURNING *;

-- name: UpdateNotification :exec
UPDATE Appointment
SET reminder_send = true
WHERE appointment_id = $1;

-- name: UpdateAppointmentStatus :exec
UPDATE Appointment
SET status = $2
WHERE appointment_id = $1;

-- name: GetAppointmentsOfDoctor :many
SELECT * FROM Appointment as a
left join Doctors as d on a.doctor_id = d.id
WHERE d.id = $1 and a.status <> 'completed';