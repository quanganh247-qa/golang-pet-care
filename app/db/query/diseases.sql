
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
LEFT JOIN disease_medicines dm ON d.id = dm.disease_id
LEFT JOIN medicines m ON dm.medicine_id = m.id
WHERE LOWER(d.name) LIKE LOWER($1);


-- name: GetDiseaseTreatmentPlanWithPhases :many
SELECT 
    d.name AS disease_name,
    d.description AS disease_description,
    d.symptoms,
    tp.phase_number,
    tp.phase_name,
    tp.description AS phase_description,
    tp.duration AS phase_duration,
    tp.notes AS phase_notes,
    m.id AS medicine_id,
    m.usage AS medicine_usage,
    m.name AS medicine_name,
    m.description AS medicine_description,
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
ORDER BY tp.phase_number, m.name;
