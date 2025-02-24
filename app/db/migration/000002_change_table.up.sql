CREATE TABLE IF NOT EXISTS task_states (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    task_id INT NOT NULL,
    state VARCHAR(255) NOT NULL CHECK (state IN ('Scheduled', 'Checked-in', 'In Progress', 'Pending Checkout', 'Closed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
    
-- Add state_id into appointments table
ALTER TABLE appointments ADD COLUMN state_id uuid;


-- Add soap note to appointments table
CREATE TABLE IF NOT EXISTS soap_notes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    appointment_id INT NOT NULL,
    note JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);