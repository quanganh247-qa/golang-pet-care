-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), false)
RETURNING id;

-- name: GetUser :one
SELECT id, username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email
FROM users
WHERE username = $1;


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;


-- name: GetAllUsers :many
SELECT * FROM users ;

-- name: UpdateUser :one
UPDATE users
SET full_name = $2, email = $3, phone_number = $4, address = $5
WHERE username = $1
RETURNING *;

-- name: UpdateAvatarUser :one
UPDATE users
SET data_image = $2, original_image = $3
WHERE username = $1
RETURNING *;


-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE username = $1 RETURNING *;

-- name: VerifiedUser :one
UPDATE users
SET is_verified_email = true
WHERE username = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;


-- -- name: InsertDoctor :one
-- INSERT INTO Doctors (
--     user_id,
--     specialization,
--     years_of_experience,
--     education,
--     certificate_number,
--     bio,
--     consultation_fee
-- ) VALUES (
--     $1, $2, $3, $4, $5, $6, $7
-- ) RETURNING *;

-- -- name: InsertDoctorSchedule :one
-- INSERT INTO DoctorSchedules (
--     doctor_id,
--     day_of_week,
--     start_time,
--     end_time,
--     is_active
--   ) VALUES (
--     $1, $2, $3, $4, $5
-- ) RETURNING *;

-- -- name: GetDoctor :one
-- SELECT 
--   d.id,
--   u.full_name AS name,
--   d.specialization,
--   d.years_of_experience,
--   d.education,
--   d.certificate_number,
--   d.bio,
--   d.consultation_fee
-- FROM
--   Doctors d
-- JOIN
--   users u ON d.user_id = u.id
-- WHERE
--   d.id = $1;

-- -- name: GetDoctorById :one
-- select * from Doctors where id = $1;

-- -- name: GetActiveDoctors :many
-- SELECT 
--   d.id,
--   u.full_name AS name,
--   d.specialization,
--   d.years_of_experience,
--   d.consultation_fee
-- FROM 
--   Doctors d
-- JOIN 
--   users u ON d.user_id = u.id
-- LEFT JOIN 
--   DoctorSchedules ds ON d.id = ds.doctor_id
-- WHERE 
--   d.is_active = true
--   AND (ds.is_active = true OR ds.is_active IS NULL)
--   AND ($1::VARCHAR IS NULL OR d.specialization = $1)
--   AND ($2::INT IS NULL OR ds.day_of_week = $2)
-- ORDER BY 
--   u.full_name;


-- -- name: GetDoctors :many
-- SELECT 
--     d.id AS doctor_id,
--     u.username,
--     u.full_name,
--     u.role,
--     d.specialization,
--     d.years_of_experience,
--     d.education,
--     d.certificate_number,
--     d.bio,
--     d.consultation_fee
-- FROM 
--     Doctors d
-- JOIN 
--     users u ON d.user_id = u.id
-- ORDER BY 
--     u.full_name;