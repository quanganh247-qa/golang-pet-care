-- -- name: CreateActivityLog :one
-- INSERT INTO ActivityLog (petID, activityType, startTime, duration, notes)
-- VALUES ($1, $2, $3, $4, $5)
-- RETURNING *;

-- -- name: GetActivityLogByID :one
-- SELECT * FROM ActivityLog WHERE logID = $1;

-- -- name: UpdateActivityLog :exec
-- UPDATE ActivityLog
-- SET activityType = $2, startTime = $3, duration = $4, notes = $5
-- WHERE logID = $1;

-- -- name: DeleteActivityLog :exec
-- DELETE FROM ActivityLog WHERE logID = $1;

-- -- name: ListActivityLogs :many
-- SELECT * FROM ActivityLog WHERE petID = $1 ORDER BY startTime DESC LIMIT $2 OFFSET $3;
