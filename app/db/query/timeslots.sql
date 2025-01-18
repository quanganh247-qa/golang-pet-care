-- name: CreateTimeSlot :one
INSERT INTO TimeSlots
( doctor_id, "date", start_time, end_time, created_at, updated_at, status)
VALUES( 
    $1, $2, $3, $4, now(), now(), 'available'
) RETURNING *;

-- name: GetTimeSlotsByDoctorAndDate :many
SELECT * from TimeSlots WHERE doctor_id = $1 AND "date" = $2 ORDER BY start_time ASC;

-- name: GetTimeSlotById :one
SELECT * from TimeSlots WHERE id = $1;

-- name: UpdateTimeSlotStatus :exec
UPDATE TimeSlots
SET status = $2
WHERE id = $1;
