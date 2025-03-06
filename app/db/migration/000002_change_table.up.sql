
-- Add soap note to appointments table
CREATE TABLE IF NOT EXISTS soap_notes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    appointment_id INT NOT NULL,
    note JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);


CREATE TABLE notification_preferences (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    topic VARCHAR(100) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(username, topic),
    FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
);

-- Các topic mặc định
INSERT INTO notification_preferences (username, topic, enabled)
SELECT username, 'appointment', true FROM users;

INSERT INTO notification_preferences (username, topic, enabled)
SELECT username, 'medical_record', true FROM users;

INSERT INTO notification_preferences (username, topic, enabled)
SELECT username, 'vaccination', true FROM users;

INSERT INTO notification_preferences (username, topic, enabled)
SELECT username, 'treatment', true FROM users;



CREATE TABLE notifications (
  id BIGSERIAL PRIMARY KEY,
  username varchar NOT NULL,
  title VARCHAR(100) NOT NULL,
  content TEXT,
  is_read BOOLEAN DEFAULT false,
  related_id INT,
  related_type VARCHAR(255),
  datetime TIMESTAMP NOT NULL,
  notify_type VARCHAR(255),
  FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
);
