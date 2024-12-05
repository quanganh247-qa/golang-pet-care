CREATE TABLE Products (
    product_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price float8 NOT NULL,
    stock_quantity INT DEFAULT 0,
    category VARCHAR(100),
    data_image BYTEA,
    original_image VARCHAR(255),
    created_at TIMESTAMP DEFAULT now(),
    is_available BOOLEAN DEFAULT true,
    removed_at TIMESTAMP DEFAULT NULL
);

-- Index for faster lookup
CREATE INDEX idx_products_category ON Products (category);
CREATE INDEX idx_products_is_available ON Products (is_available);


CREATE TABLE Cart (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL, -- Liên kết tới bảng Users
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE CartItem (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL, -- Liên kết tới bảng Cart
    product_id BIGINT NOT NULL, -- Liên kết tới sản phẩm (hoặc thú cưng)
    quantity INT DEFAULT 1, -- Số lượng
    unit_price FLOAT8 NOT NULL, -- Giá sản phẩm tại thời điểm thêm vào giỏ
    total_price FLOAT8 GENERATED ALWAYS AS (quantity * unit_price) STORED, -- Tổng giá
    FOREIGN KEY (cart_id) REFERENCES Cart (id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES Products (product_id) ON DELETE CASCADE,
    CONSTRAINT cartitem_cart_id_product_id_unique UNIQUE (cart_id, product_id)
);

CREATE TABLE Orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL, -- Liên kết tới bảng Users
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Ngày đặt hàng
    total_amount FLOAT8 NOT NULL, -- Tổng tiền của đơn hàng
    payment_status VARCHAR(20) DEFAULT 'pending', -- Trạng thái thanh toán (pending, paid, canceled)
    cart_items JSONB,
    shipping_address VARCHAR(255), -- Địa chỉ giao hàng
    notes TEXT, -- Ghi chú khách hàng
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- CREATE TABLE OrderItem (
--     id BIGSERIAL PRIMARY KEY,
--     order_id BIGINT NOT NULL, -- Liên kết tới bảng Orders
--     product_id BIGINT NOT NULL, -- Liên kết tới sản phẩm
--     quantity INT NOT NULL, -- Số lượng
--     unit_price FLOAT8 NOT NULL, -- Giá tại thời điểm đặt hàng
--     total_price FLOAT8 GENERATED ALWAYS AS (quantity * unit_price) STORED, -- Tổng giá
--     FOREIGN KEY (order_id) REFERENCES Orders (id) ON DELETE CASCADE
-- );

CREATE INDEX idx_cart_user_id ON Cart (user_id);
CREATE INDEX idx_cart_item_cart_id ON CartItem (cart_id);
CREATE INDEX idx_orders_user_id ON Orders (user_id);
-- CREATE INDEX idx_order_item_order_id ON OrderItem (order_id);
