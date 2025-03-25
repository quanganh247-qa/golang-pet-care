-- name: CreateTimeSlot :one
INSERT INTO TimeSlots
( doctor_id, "date", start_time, end_time, created_at, updated_at, status)
VALUES( 
    $1, $2, $3, $4, now(), now(), 'available'
) RETURNING *;

-- name: GetTimeSlotsByDoctorAndDate :many
SELECT * from TimeSlots WHERE doctor_id = $1 AND "date" = $2 ORDER BY start_time ASC;

-- name: GetTimeSlot :one
SELECT * FROM timeslots
WHERE id = $1 AND date = $2 AND doctor_id = $3
FOR UPDATE; -- Khóa bản ghi để tránh race condition

-- name: GetTimeSlotById :one
SELECT * from TimeSlots WHERE id = $1;

-- name: GetAvailableTimeSlots :many
SELECT 
    id,
    start_time,
    end_time,
    booked_patients,
    max_patients
FROM 
    timeslots
WHERE 
    doctor_id = $1 
    AND date = $2 
    AND booked_patients < max_patients;