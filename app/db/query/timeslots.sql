-- name: CreateTimeSlot :one
INSERT INTO time_slots 
(doctor_id, shift_id, "date", start_time, end_time, max_patients, booked_patients, created_at, updated_at)
VALUES (
   $1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
) RETURNING *;

-- -- name: GetTimeSlotsByDoctorAndDate :many
-- SELECT * from time_slots WHERE doctor_id = $1 AND "date" = $2 ORDER BY start_time ASC;

-- name: GetTimeSlot :one
SELECT * FROM time_slots
WHERE id = $1 AND date = $2 AND doctor_id = $3
FOR UPDATE; -- Lock record to avoid race condition

-- name: GetTimeSlotById :one
SELECT * from time_slots WHERE id = $1;

-- name: GetAvailableTimeSlots :many
SELECT id, start_time, end_time, max_patients, booked_patients
FROM time_slots
WHERE doctor_id = $1
AND date = $2
AND booked_patients < max_patients
ORDER BY start_time;

-- name: GetTimeSlotForUpdate :one
SELECT booked_patients, max_patients, start_time, end_time, date
FROM time_slots
WHERE id = $1
FOR UPDATE;


-- name: GetTimeSlotsByShiftID :many
SELECT * FROM time_slots
WHERE shift_id = $1
ORDER BY start_time;