-- -- name: CreateFeedingSchedule :one
-- INSERT INTO FeedingSchedule (petID, mealTime, foodType, quantity, frequency, lastFed, notes, isActive)
-- VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
-- RETURNING feedingScheduleID, petID, mealTime, foodType, quantity, frequency, lastFed, notes, isActive;

-- -- name: GetFeedingScheduleByPetID :many
-- SELECT feedingScheduleID, petID, mealTime, foodType, quantity, frequency, lastFed, notes, isActive
-- FROM FeedingSchedule
-- WHERE petID = $1
-- ORDER BY feedingScheduleID;

-- -- name: ListActiveFeedingSchedules :many
-- SELECT feedingScheduleID, petID, mealTime, foodType, quantity, frequency, lastFed, notes, isActive
-- FROM FeedingSchedule
-- WHERE isActive = true
-- ORDER BY petID;

-- -- name: UpdateFeedingSchedule :exec
-- UPDATE FeedingSchedule
-- SET mealTime = $2, foodType = $3, quantity = $4, frequency = $5, lastFed = $6, notes = $7, isActive = $8
-- WHERE feedingScheduleID = $1;

-- -- name: DeleteFeedingSchedule :exec
-- DELETE FROM FeedingSchedule
-- WHERE feedingScheduleID = $1;
