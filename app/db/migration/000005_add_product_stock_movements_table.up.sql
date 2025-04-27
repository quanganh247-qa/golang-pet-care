-- +migrate Up
CREATE TYPE movement_type_enum AS ENUM ('import', 'export');

CREATE TABLE product_stock_movements (
    movement_id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    movement_type movement_type_enum NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0), -- Ensure quantity is positive
    reason TEXT,
    movement_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    price NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);

CREATE INDEX idx_product_stock_movements_product_id ON product_stock_movements(product_id);
CREATE INDEX idx_product_stock_movements_movement_date ON product_stock_movements(movement_date);

-- Create conversations table
CREATE TABLE IF NOT EXISTS conversations (
  id BIGSERIAL PRIMARY KEY,
  type VARCHAR(20) NOT NULL CHECK (type IN ('private', 'group')),
  name VARCHAR(255), -- only for group conversations
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create conversation_participants table
CREATE TABLE IF NOT EXISTS conversation_participants (
  conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
  left_at TIMESTAMP,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (conversation_id, user_id)
);

-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
  id BIGSERIAL PRIMARY KEY,
  conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id BIGINT NOT NULL REFERENCES users(id) ON DELETE SET NULL,
  content TEXT NOT NULL,
  message_type VARCHAR(20) NOT NULL DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'file')),
  metadata JSONB,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create message_read_status table to track read status for each message
CREATE TABLE IF NOT EXISTS message_read_status (
  message_id BIGINT NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  read_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (message_id, user_id)
);

-- Add indexes for better performance
CREATE INDEX idx_conversation_participants_user_id ON conversation_participants(user_id);
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_message_read_status_user_id ON message_read_status(user_id);