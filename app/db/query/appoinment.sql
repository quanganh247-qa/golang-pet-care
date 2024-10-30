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