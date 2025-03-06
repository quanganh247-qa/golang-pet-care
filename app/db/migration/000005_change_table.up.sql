CREATE TABLE consultations (
    id SERIAL PRIMARY KEY,
    appointment_id BIGINT REFERENCES appointments(appointment_id) ON DELETE CASCADE,
    subjective TEXT, -- Thông tin từ chủ thú cưng (S)
    objective TEXT, -- Dữ liệu khám lâm sàng (O)
    assessment TEXT, -- Chẩn đoán sơ bộ (A)
    plan TEXT, -- Phác đồ điều trị (P)
    created_at TIMESTAMP DEFAULT NOW()
);
