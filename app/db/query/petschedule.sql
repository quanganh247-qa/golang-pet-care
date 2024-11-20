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
) VALUES ($1, $2, $3, $4, $5, $6, $7, false);

-- name: ListPetSchedulesByUsername :many
SELECT pet_schedule.*, pet.name
FROM pet_schedule
LEFT JOIN pet ON pet_schedule.pet_id = pet.petid
LEFT JOIN users ON pet.username = users.username
WHERE users.username = $1 and pet_schedule.removedat is null
ORDER BY pet.petid, pet_schedule.reminder_datetime;

-- name: ActiveReminder :exec
UPDATE pet_schedule
SET is_active = $2
WHERE id = $1;

-- name: DeletePetSchedule :exec
Update pet_schedule
SET removedat = now()
WHERE id = $1;