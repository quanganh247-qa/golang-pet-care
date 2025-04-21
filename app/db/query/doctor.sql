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
    u.full_name AS name,
    u.role,
    u.email
FROM doctors d
JOIN users u ON d.user_id = u.id
WHERE d.id = $1;

-- name: GetDoctorByUsername :one
SELECT 
    d.*,
    u.full_name AS name,
    u.role,
    u.email
FROM doctors d
JOIN users u ON d.user_id = u.id
WHERE u.username = $1;


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
    u.data_image,
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


-- name: CreateLeaveRequest :one
INSERT INTO leave_requests (
    doctor_id,
    start_date,
    end_date,
    leave_type,
    reason,
    status
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetLeaveRequests :many
SELECT * FROM leave_requests
WHERE doctor_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateLeaveRequest :one
UPDATE leave_requests
SET status = $2,
    updated_at = CURRENT_TIMESTAMP,
    reviewed_by = $3,
    review_notes = $4
WHERE id = $1
RETURNING *;

-- name: GetDoctorAttendance :many
SELECT s.start_time::date as work_date,
       COUNT(a.appointment_id) as total_appointments,
       COUNT(CASE WHEN a.state_id = (SELECT id FROM states WHERE state = 'Completed') THEN 1 END) as completed_appointments,
       EXTRACT(epoch FROM (s.end_time - s.start_time))/3600 as work_hours
FROM shifts s
LEFT JOIN appointments a ON a.doctor_id = s.doctor_id 
    AND a.date BETWEEN s.start_time AND s.end_time
WHERE s.doctor_id = $1
    AND s.start_time BETWEEN $2 AND $3
GROUP BY s.start_time::date
ORDER BY s.start_time::date;

-- name: GetDoctorWorkload :one
SELECT 
    COUNT(a.appointment_id) as total_appointments,
    COUNT(CASE WHEN a.state_id = (SELECT id FROM states WHERE state = 'Completed') THEN 1 END) as completed_appointments,
    ROUND(AVG(EXTRACT(epoch FROM (s.end_time - s.start_time))/3600), 2) as avg_work_hours_per_day,
    COUNT(DISTINCT s.start_time::date) as total_work_days
FROM shifts s
LEFT JOIN appointments a ON a.doctor_id = s.doctor_id 
    AND a.date BETWEEN s.start_time AND s.end_time
WHERE s.doctor_id = $1
    AND s.start_time BETWEEN $2 AND $3
GROUP BY s.doctor_id;
