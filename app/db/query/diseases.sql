
-- 1. Query cơ bản để lấy thông tin bệnh và thuốc điều trị
-- name: GetDiceaseAndMedicinesInfo :many
SELECT 
    d.id AS disease_id,
    d.name AS disease_name,
    d.description AS disease_description,
    d.symptoms,
    m.id AS medicine_id,
    m.name AS medicine_name,
    m.usage AS medicine_usage,
    m.dosage,
    m.frequency,
    m.duration,
    m.side_effects
FROM diseases d
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
LEFT JOIN medicines m ON d.id = m.disease_id
=======
LEFT JOIN disease_medicines dm ON d.id = dm.disease_id
=======
>>>>>>> e859654 (Elastic search)
LEFT JOIN medicines m ON dm.medicine_id = m.id
>>>>>>> 6c35562 (dicease and treatment plan)
=======
LEFT JOIN medicines m ON d.id = m.disease_id
>>>>>>> ada3717 (Docker file)
WHERE LOWER(d.name) LIKE LOWER($1);


-- name: GetDiseaseTreatmentPlanWithPhases :many
SELECT 
<<<<<<< HEAD
<<<<<<< HEAD
    d.id AS disease_id,
    d.name AS disease_name,
    d.description AS disease_description,
    d.symptoms,
<<<<<<< HEAD
    tp.phase_name,
    tp.description AS phase_description,
    m.id AS medicine_id,
    m.usage AS medicine_usage,
    m.name AS medicine_name,
    m.description AS medicine_description,
    pm.phase_id AS phase_id,
=======
=======
    d.id AS disease_id,
>>>>>>> 9ee4f0a (fix bug ratelimit)
    d.name AS disease_name,
    d.description AS disease_description,
    d.symptoms,
    tp.phase_number,
=======
>>>>>>> 3bf345d (happy new year)
    tp.phase_name,
    tp.description AS phase_description,
    m.id AS medicine_id,
    m.usage AS medicine_usage,
    m.name AS medicine_name,
    m.description AS medicine_description,
<<<<<<< HEAD
>>>>>>> 6c35562 (dicease and treatment plan)
=======
    pm.phase_id AS phase_id,
>>>>>>> 9ee4f0a (fix bug ratelimit)
    COALESCE(pm.dosage, m.dosage) AS dosage,
    COALESCE(pm.frequency, m.frequency) AS frequency,
    COALESCE(pm.duration, m.duration) AS duration,
    m.side_effects,
    pm.notes AS medicine_notes
FROM diseases d
JOIN treatment_phases tp ON d.id = tp.disease_id
JOIN phase_medicines pm ON tp.id = pm.phase_id
JOIN medicines m ON pm.medicine_id = m.id
WHERE LOWER(d.name) LIKE LOWER($1)
<<<<<<< HEAD
<<<<<<< HEAD
ORDER BY tp.start_date, m.name;



-- name: GetTreatmentByDiseaseId :many
SELECT 
    d.id AS disease_id,
    d.name AS disease_name,
    d.description AS disease_description,
    d.symptoms,
    tp.id AS phase_id,
    tp.phase_name AS phase_name,
    tp.description AS phase_description,
    m.id AS medicine_id,
    m.usage AS medicine_usage,
    m.name AS medicine_name,
    m.description AS medicine_description,
    COALESCE(pm.dosage, m.dosage) AS dosage,
    COALESCE(pm.frequency, m.frequency) AS frequency,
    COALESCE(pm.duration, m.duration) AS duration,
    m.side_effects
FROM diseases d
JOIN treatment_phases tp ON d.id = tp.disease_id
JOIN phase_medicines pm ON tp.id = pm.phase_id
JOIN medicines m ON pm.medicine_id = m.id
WHERE d.id = $1  LIMIT $2 OFFSET $3;

-- name: GetDiseaseByID :one
<<<<<<< HEAD
<<<<<<< HEAD
SELECT * FROM diseases WHERE id = $1 LIMIT 1;

-- name: CreateDisease :one
INSERT INTO diseases (name, description, symptoms, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *;
=======
ORDER BY tp.phase_number, m.name;
<<<<<<< HEAD
>>>>>>> 6c35562 (dicease and treatment plan)
=======
=======
ORDER BY tp.start_date, m.name;
>>>>>>> 3bf345d (happy new year)



-- name: GetTreatmentByDiseaseId :many
SELECT 
    d.id AS disease_id,
    d.name AS disease_name,
    d.description AS disease_description,
    d.symptoms,
    tp.id AS phase_id,
    tp.phase_name AS phase_name,
    tp.description AS phase_description,
    m.id AS medicine_id,
    m.usage AS medicine_usage,
    m.name AS medicine_name,
    m.description AS medicine_description,
    COALESCE(pm.dosage, m.dosage) AS dosage,
    COALESCE(pm.frequency, m.frequency) AS frequency,
    COALESCE(pm.duration, m.duration) AS duration,
    m.side_effects
FROM diseases d
JOIN treatment_phases tp ON d.id = tp.disease_id
<<<<<<< HEAD
WHERE d.id = $1  LIMIT $2 OFFSET $3;
>>>>>>> 6a85052 (get treatment by disease)
=======
JOIN phase_medicines pm ON tp.id = pm.phase_id
JOIN medicines m ON pm.medicine_id = m.id
WHERE d.id = $1  LIMIT $2 OFFSET $3;
>>>>>>> 169d7f8 (get treatment by disease)
=======
SELECT * FROM diseases WHERE id = $1 LIMIT 1;
>>>>>>> 2fe5baf (treatment phase)
=======
SELECT * FROM diseases WHERE id = $1 LIMIT 1;

-- name: CreateDisease :one
INSERT INTO diseases (name, description, symptoms, created_at, updated_at) VALUES ($1, $2, $3, now(), now()) RETURNING *;
>>>>>>> e859654 (Elastic search)
