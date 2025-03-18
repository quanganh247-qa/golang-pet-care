-- name: CreateDoctor :one
INSERT INTO doctors (
    user_id,
    specialization,
    years_of_experience,
    education,
    certificate_number,
    bio) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetDoctor :one
SELECT 
    d.*,
    u.full_name AS name
FROM doctors d
JOIN users u ON d.user_id = u.id
WHERE d.id = $1;

-- name: GetDoctorByUserId :one
SELECT * FROM doctors
WHERE user_id = $1;

-- name: ListDoctors :many
SELECT 
    d.id AS doctor_id,
    u.username,
    u.full_name,
    u.email,
    u.role,
    d.specialization,
    d.years_of_experience,
    d.education,
    d.certificate_number,
    d.bio
FROM doctors d
JOIN users u ON d.user_id = u.id
ORDER BY u.full_name;

-- name: UpdateDoctor :one
UPDATE doctors
SET 
    specialization = COALESCE($2, specialization),
    years_of_experience = COALESCE($3, years_of_experience),
    education = COALESCE($4, education),
    certificate_number = COALESCE($5, certificate_number),
    bio = COALESCE($6, bio)
    WHERE id = $1
RETURNING *;

-- name: DeleteDoctor :exec
DELETE FROM doctors
WHERE id = $1;

-- name: GetAvailableDoctors :many
SELECT DISTINCT
    d.*,
    u.full_name
FROM doctors d
JOIN users u ON d.user_id = u.id
JOIN time_slots ts ON ts.doctor_id = d.id
WHERE ts.date = $1
AND ts.booked_patients < ts.max_patients;
