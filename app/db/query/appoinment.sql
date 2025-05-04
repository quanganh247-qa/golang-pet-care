-- name: CreateAppointment :one
INSERT INTO public.appointments (
    petid, 
    username, 
    doctor_id, 
    service_id, 
    "date", 
    reminder_send, 
    time_slot_id, 
    created_at, 
    state_id, 
    appointment_reason, 
    priority, 
    arrival_time, 
    room_id, 
    confirmation_sent
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), (SELECT id FROM public.states WHERE state = 'Scheduled' LIMIT 1), $8, $9, $10, $11, $12
) RETURNING *;


-- name: UpdateTimeSlotBookedPatients :exec
UPDATE time_slots
SET booked_patients = booked_patients + 1
WHERE id = $1;

-- name: CheckinAppointment :exec
UPDATE appointments
SET state_id = (SELECT id FROM states WHERE state = 'Checked In')
WHERE appointment_id = $1;

-- name: UpdateAppointmentByID :exec
UPDATE appointments SET 
    state_id = $2,
    reminder_send = $3,
    appointment_reason = $4,
    priority = $5,
    arrival_time = $6,
    room_id = $7,
    confirmation_sent = $8,
    updated_at = now(),
    notes = $9
WHERE appointment_id = $1;

-- name: UpdateNotification :exec
UPDATE appointments
SET reminder_send = true
WHERE appointment_id = $1;

-- name: UpdateAppointmentStatus :exec
UPDATE appointments
SET state_id = $2
WHERE appointment_id = $1;

-- name: GetAppointmentsOfDoctorWithDetails :many
SELECT 
    a.appointment_id as appointment_id,
    p.name as pet_name,
    s.name as service_name,
    ts.start_time,
    ts.end_time
FROM appointments a
    LEFT JOIN doctors d ON a.doctor_id = d.id
    LEFT JOIN pets p ON a.petid = p.petid
    LEFT JOIN services s ON a.service_id = s.id
    LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
WHERE d.id = $1
AND LOWER(a.status) <> 'completed'
ORDER BY ts.start_time ASC;


-- name: CountAppointmentsByDateAndTimeSlot :one
SELECT COUNT(*) 
FROM appointments 
WHERE date = $1 AND doctor_id = $2 AND status = 'completed';

-- name: GetAppointmentsByDoctor :many
SELECT 
    a.appointment_id,
    a.date,
    a.created_at,
    a.reminder_send,
    d.id AS doctor_id,
    p.name AS pet_name,
    s.name AS service_name,
    ts.start_time,
    ts.end_time,
    ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id
FROM 
    appointments a
LEFT JOIN 
    doctors d ON a.doctor_id = d.id
LEFT JOIN 
    pets p ON a.petid = p.petid
LEFT JOIN 
    services s ON a.service_id = s.id
LEFT JOIN 
    time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN 
    states st ON a.state_id = st.id
WHERE 
    a.doctor_id = $1
ORDER BY a.created_at DESC;

-- name: ListAllAppointments :many
SELECT * FROM appointments;

-- name: GetAppointmentByStateId :many
SELECT * FROM appointments WHERE state_id = $1;


-- name: GetAllAppointmentsWithDateOption :many
SELECT 
    a.appointment_id,
    a.date,
    a.reminder_send,
    a.created_at,
    a.appointment_reason,
    a.priority,
    a.arrival_time,
    a.notes,
    p.petid AS pet_id,
    p.name AS pet_name,
    p.breed AS pet_breed,
    d.id AS doctor_id,
    s.name AS service_name,
    s.duration AS service_duration,
    ts.start_time, ts.end_time, ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id,
    u.full_name AS owner_name,
    u.phone_number AS owner_phone,
    u.email AS owner_email,
    u.address AS owner_address,
    r.name AS room_name
FROM appointments a
LEFT JOIN pets p ON a.petid = p.petid
LEFT JOIN services s ON a.service_id = s.id
LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN doctors d ON a.doctor_id = d.id
LEFT JOIN users u ON a.username = u.username
LEFT JOIN states st ON a.state_id = st.id
LEFT JOIN rooms r ON a.room_id = r.id
WHERE DATE(a.date) = DATE($1)
AND ($4 = 'false' OR st.state IN ('Confirmed', 'Scheduled'))
LIMIT $2 OFFSET $3;


-- name: GetAllAppointments :many
SELECT 
    a.appointment_id,
    a.date,
    a.reminder_send,
    a.created_at,
    a.appointment_reason,
    a.priority,
    a.arrival_time,
    p.petid as pet_id,
    p.name AS pet_name,
    p.breed AS pet_breed,
    u.full_name AS doctor_name,
    d.id AS doctor_id,
    s.name AS service_name,
    s.duration AS service_duration,
    st.state AS state_name,
    st.id AS state_id,
    u.full_name AS owner_name,
    u.phone_number AS owner_phone,
    u.email AS owner_email,
    u.address AS owner_address,
    r.name AS room_name,
    ts.start_time, ts.end_time, ts.id AS time_slot_id
FROM appointments a
LEFT JOIN pets p ON a.petid = p.petid
LEFT JOIN services s ON a.service_id = s.id
LEFT JOIN rooms r ON a.room_id = r.id
LEFT JOIN doctors d ON a.doctor_id = d.id
LEFT JOIN users u ON a.username = u.username
LEFT JOIN states st ON a.state_id = st.id
LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
ORDER BY a.created_at DESC LIMIT $1 OFFSET $2;

-- name: CountAllAppointmentsByDate :one
SELECT COUNT(*)
FROM appointments
WHERE DATE(date) = DATE($1);

-- name: GetAllAppointmentsByDate :many
SELECT 
    a.appointment_id,
    a.date ,
    a.reminder_send,
    a.created_at,
    a.appointment_reason,
    a.priority,
    a.arrival_time,
    a.notes,
    p.petid as pet_id,
    p.name AS pet_name,
    p.breed AS pet_breed,
    d.id AS doctor_id,
    s.name AS service_name,
    s.duration AS service_duration,
    ts.start_time, ts.end_time, ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id,
    u.full_name AS owner_name,
    u.phone_number AS owner_phone,
    u.email AS owner_email,
    u.address AS owner_address,
    r.name AS room_name
FROM appointments a
LEFT JOIN pets p ON a.petid = p.petid
LEFT JOIN services s ON a.service_id = s.id
LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN doctors d ON a.doctor_id = d.id
LEFT JOIN users u ON a.username = u.username
LEFT JOIN states st ON a.state_id = st.id
LEFT JOIN rooms r ON a.room_id = r.id
WHERE DATE(a.date) = DATE($1) AND st.state IN ('Confirmed', 'Scheduled')
LIMIT $2 OFFSET $3 ;


-- name: GetAppointmentDetail :one
SELECT 
    s.name AS service_name,
    p.name AS pet_name,
    st.state AS state_name
FROM services s, pets p, states st
WHERE s.id = $1 AND p.petid = $2 AND st.id = $3;

-- name: GetAppointmentDetailByAppointmentID :one
SELECT 
    a.*,
    p.petid as pet_id,
    p.name AS pet_name,
    p.breed AS pet_breed,
    d.id AS doctor_id,
    s.name AS service_name,
    s.duration AS service_duration,
    s.cost AS service_amount,
    ts.start_time, ts.end_time, ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id,
    u.full_name AS owner_name,
    u.phone_number AS owner_phone,
    u.email AS owner_email,
    u.address AS owner_address
FROM appointments a
LEFT JOIN pets p ON a.petid = p.petid
LEFT JOIN services s ON a.service_id = s.id
LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN doctors d ON a.doctor_id = d.id
LEFT JOIN users u ON a.username = u.username
LEFT JOIN states st ON a.state_id = st.id
WHERE a.appointment_id = $1;

-- name: GetAppointmentsByUser :many
SELECT 
    a.*,
    p.petid as pet_id,
    p.name AS pet_name,
    p.breed AS pet_breed,
    d.id AS doctor_id,
    s.name AS service_name,
    s.duration AS service_duration,
    s.cost AS service_amount,
    ts.start_time, ts.end_time, ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id,
    u.full_name AS owner_name,
    u.phone_number AS owner_phone,
    u.email AS owner_email,
    u.address AS owner_address
FROM appointments a
LEFT JOIN pets p ON a.petid = p.petid
LEFT JOIN services s ON a.service_id = s.id
LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN doctors d ON a.doctor_id = d.id
LEFT JOIN users u ON a.username = u.username
LEFT JOIN states st ON a.state_id = st.id
WHERE a.username = $1;

-- name: GetAppointmentsQueue :many
SELECT * FROM public.appointments 
WHERE state_id = (SELECT id FROM public.states WHERE state = 'Checked In'  LIMIT 1) 
and doctor_id = $1 and Date(arrival_time) = Date($2)
ORDER BY arrival_time ASC;


-- name: GetAvailableRoomsForDuration :many
SELECT r.*
FROM rooms r
WHERE r.status = 'available'
AND NOT EXISTS (
    SELECT 1
    FROM appointments a
    LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
    LEFT JOIN services s ON a.service_id = s.id
    WHERE a.room_id = r.id
    AND ts.start_time < $1 + interval '1 minute' * $2
    AND ts.start_time + interval '1 minute' * s.duration > $1
);


-- name: GetHistoryAppointmentsByPetID :many
SELECT 
    appointments.*,
    s.name AS service_name,
    r.name AS room_name
FROM appointments
LEFT JOIN services s ON appointments.service_id = s.id
LEFT JOIN rooms r ON appointments.room_id = r.id
WHERE petid = $1 AND state_id = (SELECT id FROM states WHERE state = 'Completed' LIMIT 1)
ORDER BY date DESC;


-- name: GetAppointmentDistribution :many
SELECT 
    s.id as service_id,
    s.name as service_name,
    COUNT(a.appointment_id) as appointment_count,
    ROUND((COUNT(a.appointment_id) * 100.0 / NULLIF((SELECT COUNT(*) FROM appointments WHERE date BETWEEN $1::date AND $2::date), 0)), 2) as percentage
FROM 
    public.services s
LEFT JOIN 
    public.appointments a ON s.id = a.service_id
    AND a.date BETWEEN $1::date AND $2::date
GROUP BY 
    s.id, s.name
ORDER BY 
    COUNT(a.appointment_id) DESC;

-- name: CountPatientsInMonth :one
SELECT COUNT(DISTINCT pet_id) 
FROM appointments
WHERE date BETWEEN $1 AND $2;

-- -- name: UpdateRoomStatus :exec
-- UPDATE rooms
-- SET status = $2,
--     current_appointment_id = $3,
--     available_at = $4
-- WHERE id = $1;



-- name: GetAppointmentByState :many
SELECT 
    a.appointment_id,
    a.date,
    a.created_at,
    a.reminder_send,
    a.appointment_reason,
    d.id AS doctor_id,
    p.name AS pet_name,
    s.name AS service_name,
    ts.start_time,
    ts.end_time,
    ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id
FROM
    appointments a
LEFT JOIN doctors d ON a.doctor_id = d.id
LEFT JOIN pets p ON a.petid = p.petid
LEFT JOIN services s ON a.service_id = s.id
LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN states st ON a.state_id = st.id
WHERE st.state = $1
ORDER BY a.created_at DESC;
