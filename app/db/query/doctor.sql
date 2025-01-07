-- name: CreateDoctorSchedule :one
INSERT INTO doctorschedules (
    doctor_id,
    day_of_week,
    shift,
    start_time,
    end_time,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;



-- name: GetDoctorSchedules :many
SELECT * FROM doctorschedules
WHERE doctor_id = $1 AND is_active = true;



-- name: UpdateDoctorSchedule :one
UPDATE doctorschedules
SET 
    start_time = COALESCE($1, start_time),
    end_time = COALESCE($2, end_time),
    is_active = COALESCE($3, is_active),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $4
RETURNING *;


