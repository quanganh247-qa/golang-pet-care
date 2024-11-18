
-- name: GetAllSchedulesByPet :many
SELECT * FROM pet_schedule where pet_id = $3 ORDER BY event_time LIMIT $1 OFFSET $2;

-- name: CreatePetSchedule :exec
INSERT INTO pet_schedule (pet_id,schedule_type, event_time, duration, activity_type, frequency, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: ListPetSchedulesByUsername :many
SELECT pet_schedule.*, pet.name
FROM pet_schedule
LEFT JOIN pet ON pet_schedule.pet_id = pet.petid
LEFT JOIN users ON pet.username = users.username
WHERE users.username = $1
ORDER BY pet.petid, pet_schedule.event_time;