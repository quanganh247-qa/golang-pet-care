-- name: GetAvailableRooms :many
SELECT id, name, type, status, current_appointment_id, available_at
FROM rooms
WHERE status = 'available' 
  AND (available_at IS NULL OR available_at <= $1)
ORDER BY id;

-- name: AssignRoomToAppointment :exec
UPDATE rooms 
SET status = 'occupied',
    current_appointment_id = $1,
    available_at = $2
WHERE id = $3;

-- name: ReleaseRoom :exec
UPDATE rooms
SET status = 'available',
    current_appointment_id = NULL,
    available_at = $1
WHERE id = $2;

-- name: GetRoomByID :one
SELECT * FROM rooms WHERE id = $1;