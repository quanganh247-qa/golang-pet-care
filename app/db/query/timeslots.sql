
-- name: DeleteDoctorSchedule :exec
DELETE FROM doctorschedules WHERE id = $1;

-- name: DeleteTimeSlot :exec
DELETE FROM timeslots WHERE id = $1;

-- name: UpdateTimeSlot :one
UPDATE timeslots
SET 
    max_patients = COALESCE($1, max_patients),
    slot_status = COALESCE($2, slot_status),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $3
RETURNING *;

-- name: CreateTimeSlot :one
INSERT INTO timeslots (
    doctor_id,
    schedule_id,
    date,
    start_time,
    end_time,
    max_patients,
    slot_status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetDoctorTimeSlots :many
SELECT * FROM timeslots
WHERE doctor_id = $1 AND date = $2;