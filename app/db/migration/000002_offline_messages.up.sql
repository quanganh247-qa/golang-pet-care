-- Create offline_messages table
CREATE TABLE IF NOT EXISTS offline_messages (
    id BIGSERIAL PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    message_type VARCHAR(50) NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    delivered_at TIMESTAMP WITH TIME ZONE NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    retry_count INT NOT NULL DEFAULT 0,
    
    CONSTRAINT fk_offline_messages_username
        FOREIGN KEY (username)
        REFERENCES users(username)
        ON DELETE CASCADE
);

-- Create indexes for faster queries
CREATE INDEX idx_offline_messages_client_id ON offline_messages(client_id);
CREATE INDEX idx_offline_messages_username ON offline_messages(username);
CREATE INDEX idx_offline_messages_status ON offline_messages(status); 