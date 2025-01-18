
CREATE TABLE Doctors (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  specialization VARCHAR(100),
  years_of_experience INT,
  education TEXT,
  certificate_number VARCHAR(50),
  bio TEXT,
  consultation_fee FLOAT8
);


-- 2. Bảng Departments (Khoa/Chuyên khoa)
CREATE TABLE Departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- 6. Bảng Time Slots (Các slot thời gian khám)
CREATE TABLE TimeSlots (
    id BIGSERIAL PRIMARY KEY,
    doctor_id INT NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    status VARCHAR(20) CHECK (status IN ('available', 'booked')) DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES Doctors(id),
    CONSTRAINT unique_slot UNIQUE (doctor_id, date, start_time)
);
