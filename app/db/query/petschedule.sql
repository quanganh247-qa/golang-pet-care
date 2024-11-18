-- name: GetAllSchedulesByPet :many
SELECT * FROM pet_schedule 
WHERE pet_id = $1
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
   notes
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ListPetSchedulesByUsername :many
SELECT pet_schedule.*, pet.name
FROM pet_schedule
LEFT JOIN pet ON pet_schedule.pet_id = pet.petid
LEFT JOIN users ON pet.username = users.username
WHERE users.username = $1
ORDER BY pet.petid, pet_schedule.reminder_datetime;