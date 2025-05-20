-- -- +migrate Up
-- CREATE TYPE movement_type_enum AS ENUM ('import', 'export');

-- CREATE TABLE product_stock_movements (
--     movement_id BIGSERIAL PRIMARY KEY,
--     product_id BIGINT NOT NULL,
--     movement_type movement_type_enum NOT NULL,
--     quantity INT NOT NULL CHECK (quantity > 0), -- Ensure quantity is positive
--     reason TEXT,
--     movement_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     price NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
--     FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
-- );

-- CREATE INDEX idx_product_stock_movements_product_id ON product_stock_movements(product_id);
-- CREATE INDEX idx_product_stock_movements_movement_date ON product_stock_movements(movement_date);

