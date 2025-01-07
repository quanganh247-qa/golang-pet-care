
CREATE TABLE Doctors (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  specialization VARCHAR(100),
  years_of_experience INT,
  education TEXT,
  certificate_number VARCHAR(50),
  bio TEXT,
  consultation_fee FLOAT8,
);


-- 2. Bảng Departments (Khoa/Chuyên khoa)
CREATE TABLE Departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 5. Bảng Doctor Schedules (Lịch làm việc của bác sĩ)
CREATE TABLE DoctorSchedules (
    id BIGSERIAL PRIMARY KEY,
    doctor_id INT NOT NULL,
    day_of_week TEXT CHECK (day_of_week IN ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday')),
    shift TEXT CHECK (shift IN ('morning', 'afternoon')) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES Doctors(id),
    CONSTRAINT unique_doctor_schedule UNIQUE (doctor_id, day_of_week, shift)
);

-- 6. Bảng Time Slots (Các slot thời gian khám)
CREATE TABLE TimeSlots (
    id BIGSERIAL PRIMARY KEY,
    doctor_id INT NOT NULL,
    schedule_id INT NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    max_patients INT,
    slot_status BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES Doctors(id),
    FOREIGN KEY (schedule_id) REFERENCES DoctorSchedules(id),
    CONSTRAINT unique_slot UNIQUE (doctor_id, date, start_time)
);
