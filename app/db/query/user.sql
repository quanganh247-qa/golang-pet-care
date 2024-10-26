-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: UpdateUser :one
UPDATE users 
SET 
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  full_name =  COALESCE(sqlc.narg(full_name),full_name),
  email = COALESCE(sqlc.narg(email),email),
  is_verified_email = COALESCE(sqlc.narg(is_verified_email),is_verified_email)
WHERE
  username = sqlc.arg(username)
RETURNING *;


-- name: InsertDoctor :one
INSERT INTO Doctors (
    user_id,
    specialization,
    years_of_experience,
    education,
    certificate_number,
    bio,
    consultation_fee
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: InsertDoctorSchedule :one
INSERT INTO DoctorSchedules (
    doctor_id,
    day_of_week,
    start_time,
    end_time,
    is_active,
    max_appointments
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

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

-- name: GetAvailableTimeSlots :many
SELECT 
  ts.id,
  ts.start_time,
  ts.end_time,
  ts.duration_minutes
FROM 
  TimeSlots ts
LEFT JOIN 
  DoctorTimeOff dto ON ts.id = dto.time_slot_id AND dto.doctor_id = $1 AND dto.start_datetime::date = $2::date
LEFT JOIN 
  Appointment a ON ts.id = a.time_slot_id AND a.doctor_id = $1 AND a.date::date = $2::date
WHERE 
  dto.id IS NULL
  AND (SELECT COUNT(*) FROM Appointment a2 WHERE a2.time_slot_id = ts.id AND a2.doctor_id = $1 AND a2.date::date = $2::date) < (SELECT max_appointments FROM DoctorSchedules ds WHERE ds.doctor_id = $1 AND ds.day_of_week = EXTRACT(DOW FROM $2::date))
ORDER BY 
  ts.start_time;