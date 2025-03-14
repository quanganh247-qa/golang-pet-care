-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email)
<<<<<<< HEAD
<<<<<<< HEAD
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), false)
RETURNING *;
<<<<<<< HEAD
=======
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), $10)
=======
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), false)
>>>>>>> eefcc96 (date time in log)
RETURNING id;
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> ada3717 (Docker file)

-- name: GetUser :one
SELECT id, username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email
FROM users
WHERE username = $1;
<<<<<<< HEAD


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

=======
>>>>>>> 0fb3f30 (user images)


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;


-- name: GetAllUsers :many
SELECT * FROM users ;

<<<<<<< HEAD
<<<<<<< HEAD
-- name: UpdateUser :one
UPDATE users
SET full_name = $2, email = $3, phone_number = $4, address = $5
=======
-- name: VerifiedUser :one
UPDATE users
<<<<<<< HEAD
SET is_verified_email = $2
>>>>>>> 6610455 (feat: redis queue)
=======
SET is_verified_email = true
>>>>>>> edfe5ad (OTP verifycation)
WHERE username = $1
RETURNING *;

<<<<<<< HEAD
-- name: UpdateAvatarUser :one
UPDATE users
SET data_image = $2, original_image = $3
WHERE username = $1
RETURNING *;


=======
-- name: UpdateUser :one
UPDATE users
SET full_name = $2, email = $3, phone_number = $4, address = $5
WHERE username = $1
RETURNING *;

<<<<<<< HEAD
>>>>>>> 473cd1d (uplaod image method)
=======
-- name: UpdateAvatarUser :one
UPDATE users
SET data_image = $2, original_image = $3
WHERE username = $1
RETURNING *;


>>>>>>> a415f25 (new data)
-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE username = $1 RETURNING *;

-- name: VerifiedUser :one
UPDATE users
SET is_verified_email = true
WHERE username = $1
RETURNING *;

=======
>>>>>>> 1f24c18 (feat: OTP with redis)
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;


<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 3003e08 (update sqlc)
-- name: InsertDoctor :one
INSERT INTO Doctors (
    user_id,
    specialization,
    years_of_experience,
    education,
    certificate_number,
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
    bio
) VALUES (
    $1, $2, $3, $4, $5, $6
=======
=======
>>>>>>> 3003e08 (update sqlc)
    bio,
    consultation_fee
=======
    bio
>>>>>>> ada3717 (Docker file)
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;
<<<<<<< HEAD

<<<<<<< HEAD
-- name: InsertDoctorSchedule :one
INSERT INTO DoctorSchedules (
    doctor_id,
    day_of_week,
    start_time,
    end_time,
    is_active
  ) VALUES (
    $1, $2, $3, $4, $5
>>>>>>> e9037c6 (update sqlc)
) RETURNING *;

<<<<<<< HEAD
=======
-- name: GetDoctor :one
SELECT 
  d.id,
  u.full_name AS name,
  d.specialization,
  d.years_of_experience,
  d.education,
  d.certificate_number,
  d.bio,
  d.consultation_fee
FROM
  Doctors d
JOIN
  users u ON d.user_id = u.id
WHERE
  d.id = $1;
=======
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
=======
>>>>>>> 3003e08 (update sqlc)

<<<<<<< HEAD

<<<<<<< HEAD
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
>>>>>>> 6f3ea8a (update sqlc)

-- -- name: GetDoctorById :one
-- select * from Doctors where id = $1;

<<<<<<< HEAD
-- name: GetActiveDoctors :many
SELECT 
  d.id,
  u.full_name AS name,
  d.specialization,
  d.years_of_experience,
  d.consultation_fee
FROM 
  Doctors d
JOIN 
  users u ON d.user_id = u.id
LEFT JOIN 
  DoctorSchedules ds ON d.id = ds.doctor_id
WHERE 
  d.is_active = true
  AND (ds.is_active = true OR ds.is_active IS NULL)
  AND ($1::VARCHAR IS NULL OR d.specialization = $1)
  AND ($2::INT IS NULL OR ds.day_of_week = $2)
ORDER BY 
  u.full_name;
>>>>>>> cfbe865 (updated service response)
=======
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
>>>>>>> 6f3ea8a (update sqlc)
=======
-- name: GetDoctor :one
SELECT 
  d.id,
  u.full_name AS name,
  d.specialization,
  d.years_of_experience,
  d.education,
  d.certificate_number,
  d.bio,
  d.consultation_fee
FROM
  Doctors d
JOIN
  users u ON d.user_id = u.id
WHERE
  d.id = $1;

-- name: GetDoctorById :one
select * from Doctors where id = $1;
>>>>>>> 685da65 (latest update)

<<<<<<< HEAD
-- name: GetAllRole :many
SELECT distinct (role) FROM users;
=======

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 685da65 (latest update)
-- name: GetDoctors :many
SELECT 
    d.id AS doctor_id,
    u.username,
    u.full_name,
    u.role,
    d.specialization,
    d.years_of_experience,
    d.education,
    d.certificate_number,
    d.bio,
    d.consultation_fee
FROM 
    Doctors d
JOIN 
    users u ON d.user_id = u.id
ORDER BY 
<<<<<<< HEAD
<<<<<<< HEAD
    u.full_name;
>>>>>>> e30b070 (Get list appoinment by user)
=======
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
>>>>>>> 6f3ea8a (update sqlc)
=======
    u.full_name;
>>>>>>> 685da65 (latest update)
=======
    u.full_name;


-- name: GetAllDoctors :many
SELECT * FROM Doctors WHERE is_active is true;
>>>>>>> 33fcf96 (Big update)
=======
>>>>>>> ffc9071 (AI suggestion)
=======

>>>>>>> ada3717 (Docker file)
