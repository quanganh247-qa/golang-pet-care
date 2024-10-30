-- name: InsertTimeslot :one
INSERT INTO Timeslots (
    doctor_id,
    start_time,
    end_time,
    day
) VALUES (
    $1, $2, $3, now()
) RETURNING *;

-- name: GetTimeslotsAvailable :many
SELECT 
    ts.doctor_id ,
    ts.start_time,
    ts.end_time
FROM 
    TimeSlots ts
JOIN 
    Doctors d ON ts.doctor_id = d.id
WHERE 
    d.id = $1  -- Replace $1 with the doctor's ID you are querying for
    AND ts.day::date = $2  -- Replace $2 with the specific date (YYYY-MM-DD)
    AND ts.is_active = true;

-- name: GetAllTimeSlots :many
SELECT 
    ts.doctor_id ,
    ts.start_time,
    ts.end_time
FROM 
    TimeSlots ts
JOIN 
    Doctors d ON ts.doctor_id = d.id
WHERE 
    d.id = $1  -- Replace $1 with the doctor's ID you are querying for
    AND ts.day::date = $2;  -- Replace $2 with the specific date (YYYY-MM-DD)

-- name: UpdateDoctorAvailable :exec
UPDATE TimeSlots
SET is_active = $1
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 59d4ef2 (modify type of filed in dtb)
WHERE id = $2;

-- name: GetTimeSlotByID :one
SELECT 
    ts.doctor_id ,
    ts.start_time,
    ts.end_time
FROM
    TimeSlots ts
WHERE
<<<<<<< HEAD
    ts.id = $1;  -- Replace $1 with the specific time slot ID you are querying for
=======
WHERE id = $2;
>>>>>>> 24ea3ee (time slot of doctor api)
=======
    ts.id = $1;  -- Replace $1 with the specific time slot ID you are querying for
>>>>>>> 59d4ef2 (modify type of filed in dtb)
