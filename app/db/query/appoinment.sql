<<<<<<< HEAD
<<<<<<< HEAD
-- name: CreateAppointment :one
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
=======
INSERT INTO Appointment (
    doctor_id,
    petid,
    service_id,
    time_slot_id,
    date,
    status
) VALUES (
    $1, $2, $3, $4, $5,'pending'
>>>>>>> cfbe865 (updated service response)
=======
INSERT INTO Appointment
( petid, doctor_id,username, service_id, "date", payment_status, notes, reminder_send, time_slot_id, created_at)
VALUES( 
<<<<<<< HEAD
    $1, $2, $3, $4, $5, $6, $7, $8, now()
>>>>>>> 685da65 (latest update)
=======
    $1, $2, $3, $4, $5, $6, $7, $8, $9,now()
>>>>>>> b393bb9 (add service and add permission)
=======
INSERT INTO appointments
( petid, doctor_id, username, service_id, "date", notes, reminder_send, time_slot_id, created_at, state_id)
VALUES( 
<<<<<<< HEAD
<<<<<<< HEAD
    $1, $2, $3, $4, $5, $6, $7, $8, $9, now()
>>>>>>> 33fcf96 (Big update)
=======
    $1, $2, $3, $4, $5, $6, $7, $8, $9, now(), $10
>>>>>>> ffc9071 (AI suggestion)
=======
    $1, $2, $3, $4, $5, $6, $7, $8, now(), $9
>>>>>>> e859654 (Elastic search)
) RETURNING *;
=======
>>>>>>> ada3717 (Docker file)

=======
>>>>>>> 4ccd381 (Update appointment flow)
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
<<<<<<< HEAD
<<<<<<< HEAD
UPDATE time_slots
SET booked_patients = booked_patients + 1
WHERE id = $1;
<<<<<<< HEAD

-- name: CheckinAppointment :exec
UPDATE appointments
SET state_id = (SELECT id FROM states WHERE state = 'Checked In')
WHERE appointment_id = $1;

<<<<<<< HEAD
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
=======
UPDATE timeslots
=======
UPDATE time_slots
>>>>>>> 33fcf96 (Big update)
SET booked_patients = booked_patients + 1
WHERE id = $1 AND doctor_id = $2;
=======
>>>>>>> ada3717 (Docker file)

<<<<<<< HEAD
>>>>>>> b393bb9 (add service and add permission)
=======
=======
>>>>>>> 4ccd381 (Update appointment flow)
-- name: UpdateAppointmentByID :exec
UPDATE appointments SET 
    state_id = $2,
    reminder_send = $3,
    appointment_reason = $4,
    priority = $5,
    arrival_time = $6,
    room_id = $7,
    confirmation_sent = $8,
    updated_at = now()
WHERE appointment_id = $1;
>>>>>>> e859654 (Elastic search)

-- name: UpdateNotification :exec
UPDATE appointments
SET reminder_send = true
WHERE appointment_id = $1;

-- name: UpdateAppointmentStatus :exec
<<<<<<< HEAD
<<<<<<< HEAD
UPDATE appointments
SET state_id = $2
<<<<<<< HEAD
=======
UPDATE Appointment
=======
UPDATE appointments
>>>>>>> 33fcf96 (Big update)
SET payment_status = $2
>>>>>>> b393bb9 (add service and add permission)
=======
>>>>>>> e859654 (Elastic search)
WHERE appointment_id = $1;

-- name: GetAppointmentsOfDoctorWithDetails :many
SELECT 
    a.appointment_id as appointment_id,
    p.name as pet_name,
    s.name as service_name,
    ts.start_time,
    ts.end_time
<<<<<<< HEAD
<<<<<<< HEAD
FROM appointments a
    LEFT JOIN doctors d ON a.doctor_id = d.id
    LEFT JOIN pets p ON a.petid = p.petid
    LEFT JOIN services s ON a.service_id = s.id
    LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
=======
FROM Appointment a
    LEFT JOIN Doctors d ON a.doctor_id = d.id
    LEFT JOIN Pet p ON a.petid = p.petid
    LEFT JOIN services s ON a.service_id = s.id
    LEFT JOIN TimeSlots ts ON a.time_slot_id = ts.id
>>>>>>> b393bb9 (add service and add permission)
=======
FROM appointments a
    LEFT JOIN doctors d ON a.doctor_id = d.id
    LEFT JOIN pets p ON a.petid = p.petid
    LEFT JOIN services s ON a.service_id = s.id
    LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
>>>>>>> 33fcf96 (Big update)
WHERE d.id = $1
AND LOWER(a.status) <> 'completed'
ORDER BY ts.start_time ASC;

<<<<<<< HEAD
<<<<<<< HEAD

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


-- name: GetAllAppointments :many
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
    a.appointment_id, a.date, a.created_at,
    p.name AS pet_name,
    d.id AS doctor_id,
    s.name AS service_name,
    ts.start_time, ts.end_time,
    st.state
FROM appointments a
LEFT JOIN pets p ON p.petid = a.petid
LEFT JOIN doctors d ON d.id = a.doctor_id
LEFT JOIN services s ON s.id = a.service_id
LEFT JOIN time_slots ts ON ts.id = a.time_slot_id
LEFT JOIN states st ON st.id = a.state_id
WHERE a.username = $1;

-- name: GetAppointmentsQueue :many
SELECT * FROM public.appointments 
WHERE state_id <> (SELECT id FROM public.states WHERE state = 'Scheduled' LIMIT 1) and doctor_id = $1
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
=======
-- name: GetAppointmentDetailById :one
<<<<<<< HEAD
SELECT * from Appointment WHERE appointment_id = $1;
<<<<<<< HEAD
>>>>>>> 7e35c2e (get appointment detail)
=======
=======
SELECT * from appointments WHERE appointment_id = $1;
>>>>>>> 33fcf96 (Big update)

-- name: GetAppointmentsByUser :many
SELECT 
    p.*, s.*, a.*, ts.*
FROM 
    appointments a
JOIN 
    pets p ON a.petid = p.petid 
JOIN 
    services s ON a.service_id = s.id 
JOIN 
    time_slots ts ON a.time_slot_id = ts.id
WHERE 
<<<<<<< HEAD
    u.username = $1 and p.is_active is true
ORDER BY 
    a.date DESC;
>>>>>>> e30b070 (Get list appoinment by user)
=======
    a.username = $1 and a.status <> 'completed';
=======
>>>>>>> dc47646 (Optimize SQL query)

-- name: CountAppointmentsByDateAndTimeSlot :one
SELECT COUNT(*) 
FROM appointments 
WHERE date = $1 AND doctor_id = $2 AND status = 'completed';
<<<<<<< HEAD
>>>>>>> 685da65 (latest update)
=======

-- name: GetAppointmentsByDoctor :many
SELECT 
    a.appointment_id,
    a.date,
    a.created_at,
    a.notes,
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
<<<<<<< HEAD
    a.doctor_id = $1;
<<<<<<< HEAD
>>>>>>> b393bb9 (add service and add permission)
=======
=======
    a.doctor_id = $1
ORDER BY a.created_at DESC;
>>>>>>> dc47646 (Optimize SQL query)

-- name: ListAllAppointments :many
SELECT * FROM appointments;
<<<<<<< HEAD
>>>>>>> 33fcf96 (Big update)
=======

-- name: GetAppointmentByStateId :many
SELECT * FROM appointments WHERE state_id = $1;

-- name: GetAllAppointments :many
SELECT 
    a.appointment_id, a.date, a.reminder_send, a.created_at,
    p.name AS pet_name,
    d.id AS doctor_id,
    s.name AS service_name,
    ts.start_time, ts.end_time, ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id
FROM appointments a
LEFT JOIN pets p ON a.petid = p.petid
LEFT JOIN services s ON a.service_id = s.id
LEFT JOIN time_slots ts ON a.time_slot_id = ts.id
LEFT JOIN doctors d ON a.doctor_id = d.id
LEFT JOIN states st ON a.state_id = st.id;

<<<<<<< HEAD
>>>>>>> ffc9071 (AI suggestion)
=======
-- name: GetSOAPByAppointmentID :one
SELECT * FROM consultations WHERE appointment_id = $1;
<<<<<<< HEAD
>>>>>>> e859654 (Elastic search)
=======

<<<<<<< HEAD
>>>>>>> ada3717 (Docker file)
=======
-- name: GetAppointmentDetail :one
SELECT 
    s.name AS service_name,
    p.name AS pet_name,
    st.state AS state_name
FROM services s, pets p, states st
WHERE s.id = $1 AND p.petid = $2 AND st.id = $3;

-- name: GetAppointmentDetailByAppointmentID :one
SELECT 
    a.appointment_id, a.date,  a.reminder_send, a.created_at, a.appointment_reason, a.priority, a.arrival_time, a.room_id, a.confirmation_sent,
    d.id AS doctor_id,
    p.name AS pet_name,
    s.name AS service_name,
    ts.start_time, ts.end_time, ts.id AS time_slot_id,
    st.state AS state_name,
    st.id AS state_id
FROM appointments a
LEFT JOIN pets p ON p.petid = a.petid
LEFT JOIN doctors d ON d.id = a.doctor_id
LEFT JOIN services s ON s.id = a.service_id
LEFT JOIN time_slots ts ON ts.id = a.time_slot_id
LEFT JOIN states st ON st.id = a.state_id
WHERE a.appointment_id = $1;

-- name: GetAppointmentsByUser :many
SELECT 
    a.appointment_id, a.date, a.created_at,
    p.name AS pet_name,
    d.id AS doctor_id,
    s.name AS service_name,
    ts.start_time, ts.end_time,
    st.state
FROM appointments a
LEFT JOIN pets p ON p.petid = a.petid
LEFT JOIN doctors d ON d.id = a.doctor_id
LEFT JOIN services s ON s.id = a.service_id
LEFT JOIN time_slots ts ON ts.id = a.time_slot_id
LEFT JOIN states st ON st.id = a.state_id
WHERE a.username = $1;
<<<<<<< HEAD
>>>>>>> dc47646 (Optimize SQL query)
=======

-- name: GetAppointmentsQueue :many
SELECT * FROM public.appointments 
WHERE state_id <> (SELECT id FROM public.states WHERE state = 'Scheduled' LIMIT 1)
ORDER BY arrival_time ASC;
>>>>>>> 4ccd381 (Update appointment flow)
