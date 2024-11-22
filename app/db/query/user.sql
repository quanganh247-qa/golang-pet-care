-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email)
<<<<<<< HEAD
<<<<<<< HEAD
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), false)
RETURNING *;
=======
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), $10)
=======
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), false)
>>>>>>> eefcc96 (date time in log)
RETURNING id;
>>>>>>> 0fb3f30 (user images)

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
SET full_name = $2, email = $3, phone_number = $4, address = $5, data_image = $6, original_image = $7
WHERE username = $1
RETURNING *;

>>>>>>> 473cd1d (uplaod image method)
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


-- name: InsertDoctor :one
INSERT INTO Doctors (
    user_id,
    specialization,
    years_of_experience,
    education,
    certificate_number,
    bio
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: GetAllRole :many
SELECT distinct (role) FROM users;