-- -- name: SaveState :one
-- INSERT INTO oauth_states (state, username, created_at) VALUES ($1, $2, $3) RETURNING *;

-- -- name: GetState :one
-- SELECT state, username, created_at
--     FROM oauth_states
--     WHERE state = $1 LIMIT 1;

-- -- name: DeleteState :exec
-- DELETE FROM oauth_states WHERE state = $1 ;

-- -- name: CleanExpiredStates :exec
-- DELETE FROM oauth_states WHERE created_at < $1;