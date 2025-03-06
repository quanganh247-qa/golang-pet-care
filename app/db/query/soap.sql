-- name: CreateSOAP :one
INSERT INTO consultations (appointment_id, subjective, objective, assessment, plan)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

<<<<<<< HEAD
-- name: GetSOAPByAppointmentID :one
SELECT * FROM consultations WHERE appointment_id = $1;


-- name: UpdateSOAP :one
UPDATE consultations SET subjective = $2, objective = $3, assessment = $4
WHERE appointment_id = $1 RETURNING *;

=======
-- name: GetSOAP :one
SELECT * FROM consultations WHERE appointment_id = $1;

-- name: UpdateSOAP :one
UPDATE consultations SET subjective = $2, objective = $3, assessment = $4, plan = $5 
WHERE appointment_id = $1 RETURNING *;
>>>>>>> e859654 (Elastic search)
