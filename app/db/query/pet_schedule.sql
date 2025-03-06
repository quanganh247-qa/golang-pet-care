-- name: GetAllSchedulesByPet :many
SELECT * FROM pet_schedule 
WHERE pet_id = $1 and removedat is null
ORDER BY reminder_datetime 
LIMIT $2 OFFSET $3;

-- name: CreatePetSchedule :exec
INSERT INTO pet_schedule (
   pet_id,
   title,
   reminder_datetime,
   event_repeat,
   end_type,
   end_date,
   notes,
   is_active
) VALUES ($1, $2, $3, $4, $5, $6, $7, true);

-- name: ListPetSchedulesByUsername :many
SELECT pet_schedule.*, pets.name
FROM pet_schedule
LEFT JOIN pets ON pet_schedule.pet_id = pets.petid
LEFT JOIN users ON pets.username = users.username
WHERE users.username = $1 and pet_schedule.removedat is null
ORDER BY pets.petid, pet_schedule.reminder_datetime;

-- name: ActiveReminder :exec
UPDATE pet_schedule
SET is_active = $2
WHERE id = $1;

-- name: DeletePetSchedule :exec
Update pet_schedule
SET removedat = now()
WHERE id = $1;

-- name: UpdatePetSchedule :exec
UPDATE pet_schedule
SET title = $2,
    reminder_datetime = $3,
    event_repeat = $4,
    end_type = $5,
    end_date = $6,
    notes = $7,
    is_active = $8
WHERE id = $1;

-- name: GetPetScheduleById :one
SELECT * FROM pet_schedule
WHERE id = $1;