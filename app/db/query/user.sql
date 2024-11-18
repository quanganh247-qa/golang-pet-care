-- name: CreateUser :one
INSERT INTO users (username, hashed_password, full_name, email, phone_number, address, data_image, original_image, role, created_at, is_verified_email)
<<<<<<< HEAD
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), false)
RETURNING *;
=======
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), $10)
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
-- name: UpdateUser :one
UPDATE users
SET full_name = $2, email = $3, phone_number = $4, address = $5
=======
-- name: VerifiedUser :one
UPDATE users
SET is_verified_email = $2
>>>>>>> 6610455 (feat: redis queue)
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