CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT, -- Nếu bạn muốn liên kết với người dùng đã tải lên
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);