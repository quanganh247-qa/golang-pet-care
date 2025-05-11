-- Migration to add pet weight history tracking
CREATE TABLE public.pet_weight_history (
    id bigserial NOT NULL,
    pet_id bigint NOT NULL,
    weight_kg float8 NOT NULL,
    recorded_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes text,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT pet_weight_history_pkey PRIMARY KEY (id),
    CONSTRAINT pet_weight_history_pet_id_fkey FOREIGN KEY (pet_id) REFERENCES public.pets(petid) ON DELETE CASCADE
);

-- Create indexes for improved query performance
CREATE INDEX idx_pet_weight_history_pet_id ON public.pet_weight_history (pet_id);
CREATE INDEX idx_pet_weight_history_recorded_at ON public.pet_weight_history (recorded_at);

-- Function to add new weight record and update the pet's current weight
CREATE OR REPLACE FUNCTION add_pet_weight_record(
    p_pet_id bigint,
    p_weight_kg float8,
    p_notes text DEFAULT NULL
) RETURNS json
LANGUAGE plpgsql
AS $$
DECLARE
    result json;
BEGIN
    -- Insert the new weight record
    INSERT INTO pet_weight_history (pet_id, weight_kg, notes)
    VALUES (p_pet_id, p_weight_kg, p_notes);
    
    -- Update the pet's current weight
    UPDATE pets 
    SET weight = p_weight_kg::float8, updated_at = NOW()
    WHERE petid = p_pet_id;
    
    -- Return success
    RETURN json_build_object(
        'success', true,
        'message', 'Weight record added successfully'
    );
EXCEPTION WHEN OTHERS THEN
    RETURN json_build_object(
        'success', false,
        'message', 'Failed to add weight record: ' || SQLERRM
    );
END;
$$; 