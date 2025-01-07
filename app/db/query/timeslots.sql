-- name: InsertTimeslot :one
INSERT INTO Timeslots (
    doctor_id,
    start_time,
    end_time,
    day
) VALUES (
    $1, $2, $3, now()
) RETURNING *;

