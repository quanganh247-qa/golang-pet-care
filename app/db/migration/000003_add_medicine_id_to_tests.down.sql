-- Remove the foreign key constraint first
ALTER TABLE tests DROP CONSTRAINT IF EXISTS tests_medicine_id_fkey;

-- Remove the index
DROP INDEX IF EXISTS idx_tests_medicine_id;

-- Remove the medicine_id column
ALTER TABLE tests DROP COLUMN IF EXISTS medicine_id; 