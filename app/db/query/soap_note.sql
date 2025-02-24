-- name: GetSoapNote :one
SELECT * FROM soap_notes WHERE appointment_id = $1;

-- name: CreateSoapNote :one
INSERT INTO soap_notes (appointment_id, note) VALUES ($1, $2) RETURNING *;

-- name: UpdateSoapNote :one
UPDATE soap_notes SET note = $2 WHERE id = $1 RETURNING *;