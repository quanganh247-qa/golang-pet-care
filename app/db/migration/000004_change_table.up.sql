CREATE TABLE IF NOT EXISTS states (
    id BIGSERIAL PRIMARY KEY,
    state VARCHAR(255) NOT NULL CHECK (state IN ('Scheduled', 'Confirmed', 'In Progress', 'Completed', 'Closed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);

ALTER TABLE appointments ADD COLUMN state_id INT;
ALTER TABLE appointments ADD FOREIGN KEY (state_id) REFERENCES states(id);
