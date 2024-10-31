-- name: CreateFeedingSchedule :one
INSERT INTO FeedingSchedule (petid, meal_time, food_type, quantity, frequency, lastfed, notes, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING feeding_schedule_id, petid, meal_time, food_type, quantity, frequency, lastfed, notes, is_active;

-- name: GetFeedingScheduleByPetID :many
SELECT feeding_schedule_id, petid, meal_time, food_type, quantity, frequency, lastfed, notes, is_active
FROM FeedingSchedule
WHERE petid = $1
ORDER BY feeding_schedule_id;

-- name: ListActiveFeedingSchedules :many
SELECT feeding_schedule_id, petid, meal_time, food_type, quantity, frequency, lastfed, notes, is_active
FROM FeedingSchedule
WHERE is_active = true
ORDER BY petid;

-- name: UpdateFeedingSchedule :exec
UPDATE FeedingSchedule
SET meal_time = $2, food_type = $3, quantity = $4, frequency = $5, lastfed = $6, notes = $7, is_active = $8
WHERE feeding_schedule_id = $1;

-- name: DeleteFeedingSchedule :exec
DELETE FROM FeedingSchedule
WHERE feeding_schedule_id = $1;
