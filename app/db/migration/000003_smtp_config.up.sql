-- Create smtp_configs table to store SMTP configurations
CREATE TABLE IF NOT EXISTS smtp_configs (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    smtp_host VARCHAR(255) NOT NULL DEFAULT 'smtp.gmail.com',
    smtp_port VARCHAR(10) NOT NULL DEFAULT '587',
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create index for faster queries
CREATE INDEX idx_smtp_configs_email ON smtp_configs(email);
CREATE INDEX idx_smtp_configs_is_default ON smtp_configs(is_default);
