-- Drop function
DROP FUNCTION IF EXISTS add_pet_weight_record;

-- Drop indexes
DROP INDEX IF EXISTS idx_pet_weight_history_pet_id;
DROP INDEX IF EXISTS idx_pet_weight_history_recorded_at;

-- Drop table
DROP TABLE IF EXISTS pet_weight_history; 