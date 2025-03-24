-- name: GetAvailableRooms :many
SELECT id, name, type, status, current_appointment_id, available_at
FROM rooms
WHERE status = 'available' 
LIMIT $1 OFFSET $2;
  
-- name: AssignRoomToAppointment :exec
UPDATE rooms 
SET current_appointment_id = $2
WHERE id = $1;

-- name: ReleaseRoom :exec
UPDATE rooms
SET status = 'available',
    current_appointment_id = NULL,
    available_at = $1
WHERE id = $2;

-- name: GetRoomByID :one
SELECT * FROM rooms WHERE id = $1;

-- name: CreateRoom :one
INSERT INTO rooms (name, type, status, current_appointment_id, available_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateRoom :exec
UPDATE rooms
SET name = $2,
    type = $3,
    status = $4,
    current_appointment_id = $5,
    available_at = $6
WHERE id = $1;


-- name: DeleteRoom :exec
DELETE FROM rooms WHERE id = $1;