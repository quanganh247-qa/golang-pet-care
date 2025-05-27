-- Add medicine_id column to tests table for vaccine inventory management
ALTER TABLE tests ADD COLUMN medicine_id int8 NULL;

-- Add foreign key constraint to medicines table
ALTER TABLE tests ADD CONSTRAINT tests_medicine_id_fkey 
    FOREIGN KEY (medicine_id) REFERENCES medicines(id) ON DELETE SET NULL;

-- Add index for medicine_id for better query performance
CREATE INDEX idx_tests_medicine_id ON tests USING btree (medicine_id);

-- Add comment to clarify the purpose
COMMENT ON COLUMN tests.medicine_id IS 'References medicine inventory for vaccines (null for regular tests)'; 