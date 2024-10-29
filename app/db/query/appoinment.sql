-- name: CreateAppointment :one
INSERT INTO Appointments (
    doctor_id,
    patient_id,
    ServiceID,
    time_slot_id
    status
) VALUES (
    $1, $2, $3, $4, 'pending'
) RETURNING *;