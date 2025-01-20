<<<<<<< HEAD
<<<<<<< HEAD
-- name: CreateTimeSlot :one
INSERT INTO time_slots 
(doctor_id, "date", start_time, end_time, max_patients, booked_patients, created_at, updated_at)
VALUES (
   $1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
) RETURNING *;

<<<<<<< HEAD
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
SELECT booked_patients, max_patients, start_time, end_time
FROM time_slots
WHERE id = $1
FOR UPDATE;
=======
>>>>>>> ae87825 (updated)
=======

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

=======
>>>>>>> 685da65 (latest update)
-- name: CreateTimeSlot :one
INSERT INTO TimeSlots
( doctor_id, "date", start_time, end_time, created_at, updated_at, status)
VALUES( 
    $1, $2, $3, $4, now(), now(), 'available'
) RETURNING *;

<<<<<<< HEAD
-- name: GetDoctorTimeSlots :many
SELECT * FROM timeslots
WHERE doctor_id = $1 AND date = $2;
>>>>>>> e9037c6 (update sqlc)
=======
-- name: GetTimeSlotsByDoctorAndDate :many
SELECT * from TimeSlots WHERE doctor_id = $1 AND "date" = $2 ORDER BY start_time ASC;

-- name: GetTimeSlot :one
SELECT * FROM timeslots
WHERE id = $1 AND date = $2 AND doctor_id = $3
FOR UPDATE; -- Khóa bản ghi để tránh race condition

-- name: GetTimeSlotById :one
SELECT * from TimeSlots WHERE id = $1;

<<<<<<< HEAD
-- name: UpdateTimeSlotStatus :exec
UPDATE TimeSlots
SET status = $2
WHERE id = $1;
>>>>>>> 685da65 (latest update)
=======
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
>>>>>>> b393bb9 (add service and add permission)
