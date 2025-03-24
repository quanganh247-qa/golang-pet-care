-- name: GetAvailableRooms :many
SELECT id, name, type, status, current_appointment_id, available_at
FROM rooms
WHERE status = 'available' 
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
LIMIT $1 OFFSET $2;
  
-- name: AssignRoomToAppointment :exec
UPDATE rooms 
SET current_appointment_id = $2
WHERE id = $1;
=======
=======
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
  AND (available_at IS NULL OR available_at <= $1)
ORDER BY id;

-- name: AssignRoomToAppointment :exec
UPDATE rooms 
SET status = 'occupied',
    current_appointment_id = $1,
    available_at = $2
WHERE id = $3;
<<<<<<< HEAD
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
LIMIT $1 OFFSET $2;
  
-- name: AssignRoomToAppointment :exec
UPDATE rooms 
SET current_appointment_id = $2
WHERE id = $1;
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
LIMIT $1 OFFSET $2;
  
-- name: AssignRoomToAppointment :exec
UPDATE rooms 
SET current_appointment_id = $2
WHERE id = $1;
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)

-- name: ReleaseRoom :exec
UPDATE rooms
SET status = 'available',
    current_appointment_id = NULL,
    available_at = $1
WHERE id = $2;

-- name: GetRoomByID :one
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
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
<<<<<<< HEAD
<<<<<<< HEAD
DELETE FROM rooms WHERE id = $1;
=======
SELECT * FROM rooms WHERE id = $1;
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
DELETE FROM rooms WHERE id = $1;
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
SELECT * FROM rooms WHERE id = $1;
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
DELETE FROM rooms WHERE id = $1;
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
