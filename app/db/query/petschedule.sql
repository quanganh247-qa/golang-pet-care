
-- name: GetAllSchedulesByPet :many
SELECT * FROM pet_schedule where pet_id = $3 ORDER BY pet_id LIMIT $1 OFFSET $2;

-- name: CreatePetSchedule :exec
INSERT INTO pet_schedule (pet_id,schedule_type, event_time, duration, activity_type, frequency, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7);