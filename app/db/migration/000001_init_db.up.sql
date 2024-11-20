-- Create users table with proper constraints
CREATE TABLE users (
    id BIGSERIAL NOT NULL,
    username VARCHAR NOT NULL,
    hashed_password VARCHAR NOT NULL,
    full_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE, -- Added unique constraint on email
    phone_number VARCHAR,
    address VARCHAR,
    data_image BYTEA ,
    original_image VARCHAR(255),
    role VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    is_verified_email BOOL DEFAULT false,
    removed_at TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (id),  
    CONSTRAINT users_username_key UNIQUE (username)
);


CREATE TABLE verify_emails (
  id BIGSERIAL NOT NULL,
  username varchar NOT NULL,
  email varchar NOT NULL,
  secret_code int8 NOT NULL,
  is_used bool NOT NULL DEFAULT false,
  created_at timestamp DEFAULT (now()),
  expired_at timestamp DEFAULT (now()+'00:15:00'::interval),
  PRIMARY KEY (id)
);

CREATE TABLE Pet (
  petid BIGSERIAL NOT NULL,
  name varchar(100) NOT NULL,
  type varchar(50) NOT NULL,
  breed varchar(100),
  age int4,
  gender varchar(10),
  healthnotes text,
  weight float8,
  birth_date date,
  username varchar NOT NULL,
  microchip_number varchar(50),
  last_checkup_date date,
  is_active BOOLEAN DEFAULT true,
  data_image BYTEA ,
  original_image VARCHAR(255),
  PRIMARY KEY (petid)
);


CREATE TABLE Vaccination (
  vaccinationID BIGSERIAL PRIMARY KEY,
  petID BIGINT,
  vaccineName VARCHAR(100) NOT NULL,
  dateAdministered timestamp NOT NULL,
  nextDueDate timestamp,
  vaccineProvider VARCHAR(100),
  batchNumber VARCHAR(50),
  notes TEXT
);

CREATE TABLE pet_schedule (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT REFERENCES Pet(petid),
    title VARCHAR(255),
    reminder_datetime timestamp,
    event_repeat VARCHAR(50),
    end_type bool DEFAULT false,
    end_date DATE,
    notes TEXT,
    is_active BOOLEAN DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    removedat TIMESTAMP DEFAULT NULL
);


CREATE TABLE ServiceType (
  typeID BIGSERIAL PRIMARY KEY,
  serviceTypeName varchar NOT NULL,
  description TEXT,
  iconURL VARCHAR(255)
);

CREATE TABLE Service (
  serviceID BIGSERIAL PRIMARY KEY,
  typeID BIGINT,
  name varchar NOT NULL,
  price float8,
  duration INTERVAL,
  description TEXT,
  isAvailable BOOLEAN DEFAULT true,
  removed_at timestamp
);

CREATE TABLE Appointment (
  appointment_id BIGSERIAL PRIMARY KEY,
  petid BIGINT,
  doctor_id BIGINT,
  service_id BIGINT,
  date timestamp DEFAULT (now()),
  status VARCHAR(20),
  notes TEXT,
  reminder_send BOOLEAN DEFAULT false,
  time_slot_id BIGINT,
  created_at timestamp DEFAULT (now())
);

CREATE TABLE Checkout (
  checkout_id BIGSERIAL PRIMARY KEY,
  petid BIGINT,
  doctor_id BIGINT,
  date timestamp DEFAULT (now()),
  total_tmount float8 NOT NULL,
  payment_status VARCHAR(20),
  payment_method VARCHAR(50),
  notes TEXT
);

CREATE TABLE CheckoutService (
  checkoutService_ID BIGSERIAL PRIMARY KEY,
  checkoutID BIGINT,
  serviceID BIGINT,
  quantity INT DEFAULT 1,
  unitPrice float8,
  subtotal float8
);


CREATE TABLE Doctors (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  specialization VARCHAR(100),
  years_of_experience INT,
  education TEXT,
  certificate_number VARCHAR(50),
  bio TEXT,
  consultation_fee DECIMAL(10,2)
);

CREATE TABLE TimeSlots (
  id BIGSERIAL PRIMARY KEY,
  doctor_id BIGINT NOT NULL,
  start_time timestamp NOT NULL,
  end_time timestamp NOT NULL,
  is_active BOOLEAN DEFAULT true,
  day date NOT NULL,
  FOREIGN KEY (doctor_id) REFERENCES Doctors (id)
);

CREATE TABLE DoctorSchedules (
  id BIGSERIAL PRIMARY KEY,
  doctor_id BIGINT NOT NULL,
  day_of_week INT,
  start_time timestamp NOT NULL,
  end_time timestamp NOT NULL,
  is_active BOOLEAN DEFAULT true,
  max_appointments INT DEFAULT 1
);


-- Create device tokens table with proper foreign key reference
CREATE TABLE DeviceTokens (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR NOT NULL,
  token VARCHAR NOT NULL UNIQUE,
  device_type VARCHAR(50),
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  last_used_at TIMESTAMP,
  expired_at TIMESTAMP,
  CONSTRAINT fk_device_tokens_username 
      FOREIGN KEY (username) 
      REFERENCES users(username) 
      ON DELETE CASCADE
);

-- Create notifications table with necessary fields
CREATE TABLE notifications (
  notificationID BIGSERIAL PRIMARY KEY,
  username varchar NOT NULL,
  title VARCHAR(100) NOT NULL,
  description TEXT,
  datetime TIMESTAMP NOT NULL,
  is_read BOOLEAN DEFAULT false
);


-- Create diseases table
CREATE TABLE diseases (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    symptoms JSONB, -- Store symptoms as JSON array
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create medicines table
CREATE TABLE medicines (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    usage TEXT,
    dosage TEXT,
    frequency TEXT,
    duration TEXT,
    side_effects TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create junction table for many-to-many relationship
CREATE TABLE disease_medicines (
    disease_id BIGINT REFERENCES diseases(id),
    medicine_id BIGINT REFERENCES medicines(id),
    PRIMARY KEY (disease_id, medicine_id)
);


-- 2. Query chi tiết hơn với thông tin phác đồ điều trị theo từng giai đoạn
CREATE TABLE treatment_phases (
    id BIGSERIAL PRIMARY KEY,
    disease_id BIGINT REFERENCES diseases(id),
    phase_number INT,
    phase_name VARCHAR(255),
    description TEXT,
    duration VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE phase_medicines (
    phase_id BIGINT REFERENCES treatment_phases(id),
    medicine_id BIGINT REFERENCES medicines(id),
    dosage TEXT,
    frequency TEXT,
    duration TEXT,
    notes TEXT,
    PRIMARY KEY (phase_id, medicine_id)
);
-- 4. Query để lấy lịch sử điều trị của một thú cưng
CREATE TABLE pet_treatments (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT REFERENCES Pet(petid),
    disease_id BIGINT REFERENCES diseases(id),
    start_date DATE,
    end_date DATE,
    status VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE treatment_progress (
    id BIGSERIAL PRIMARY KEY,
    treatment_id BIGINT REFERENCES pet_treatments(id),
    phase_id BIGINT REFERENCES treatment_phases(id),
    start_date DATE,
    end_date DATE,
    status VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pet_logs (
    log_id BIGSERIAL PRIMARY KEY,
	petid int8 NOT NULL,
	datetime timestamp NULL,
	title varchar NULL,
	notes text NULL,
	CONSTRAINT newtable_pet_fk FOREIGN KEY (petid) REFERENCES pet(petid)
);



ALTER TABLE Pet ADD CONSTRAINT pet_users_fk FOREIGN KEY (username) REFERENCES users (username);

ALTER TABLE Vaccination ADD CONSTRAINT vaccination_pet_fk FOREIGN KEY (petID) REFERENCES Pet (petid);

ALTER TABLE Service ADD CONSTRAINT service_type_fk FOREIGN KEY (typeID) REFERENCES ServiceType (typeID);

ALTER TABLE Appointment ADD CONSTRAINT appointment_pet_fk FOREIGN KEY (petid) REFERENCES Pet (petid);

ALTER TABLE Appointment ADD CONSTRAINT appointment_service_fk FOREIGN KEY (service_id) REFERENCES Service (serviceID);

ALTER TABLE CheckoutService ADD CONSTRAINT cs_checkout_fk FOREIGN KEY (checkoutID) REFERENCES Checkout (checkout_id);

ALTER TABLE CheckoutService ADD CONSTRAINT cs_service_fk FOREIGN KEY (serviceID) REFERENCES Service (serviceID);

ALTER TABLE Doctors ADD CONSTRAINT fk_doctor_user FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE DoctorSchedules ADD CONSTRAINT fk_schedule_doctor FOREIGN KEY (doctor_id) REFERENCES Doctors (id);


-- Index for users table
CREATE INDEX idx_users_created_at ON users (created_at);

-- Index for Pet table
CREATE INDEX idx_pet_username ON Pet (username);
CREATE INDEX idx_pet_is_active ON Pet (is_active);
CREATE INDEX idx_pet_birth_date ON Pet (birth_date);

-- Index for Vaccination table
CREATE INDEX idx_vaccination_pet_id ON Vaccination (petID);
CREATE INDEX idx_vaccination_date_administered ON Vaccination (dateAdministered);

-- Index for pet_schedule table
CREATE INDEX idx_pet_schedule_pet_id ON pet_schedule (pet_id);
CREATE INDEX idx_pet_schedule_reminder_datetime ON pet_schedule (reminder_datetime);

-- Index for Service table
CREATE INDEX idx_service_type_id ON Service (typeID);
CREATE INDEX idx_service_is_available ON Service (isAvailable);

-- Index for Appointment table
CREATE INDEX idx_appointment_pet_id ON Appointment (petid);
CREATE INDEX idx_appointment_service_id ON Appointment (service_id);
CREATE INDEX idx_appointment_doctor_id ON Appointment (doctor_id);
CREATE INDEX idx_appointment_date ON Appointment (date);

-- Index for Checkout table
CREATE INDEX idx_checkout_pet_id ON Checkout (petid);
CREATE INDEX idx_checkout_doctor_id ON Checkout (doctor_id);
CREATE INDEX idx_checkout_date ON Checkout (date);

-- Index for Doctors table
CREATE INDEX idx_doctors_user_id ON Doctors (user_id);

-- Index for TimeSlots table
CREATE INDEX idx_timeslots_doctor_id ON TimeSlots (doctor_id);
CREATE INDEX idx_timeslots_day ON TimeSlots (day);



-- Index for diseases table
CREATE INDEX idx_diseases_name ON diseases (name);

-- Index for medicines table
CREATE INDEX idx_medicines_name ON medicines (name);

-- Index for pet_treatments table
CREATE INDEX idx_pet_treatments_pet_id ON pet_treatments (pet_id);
CREATE INDEX idx_pet_treatments_disease_id ON pet_treatments (disease_id);
CREATE INDEX idx_pet_treatments_start_date ON pet_treatments (start_date);

-- Index for treatment_phases table
CREATE INDEX idx_treatment_phases_disease_id ON treatment_phases (disease_id);

-- Index for phase_medicines table
CREATE INDEX idx_phase_medicines_phase_id ON phase_medicines (phase_id);
CREATE INDEX idx_phase_medicines_medicine_id ON phase_medicines (medicine_id);

-- Index for treatment_progress table
CREATE INDEX idx_treatment_progress_treatment_id ON treatment_progress (treatment_id);

-- Index for pet_logs table
CREATE INDEX idx_pet_logs_pet_id ON pet_logs (petid);
CREATE INDEX idx_pet_logs_datetime ON pet_logs (datetime);


INSERT INTO users (username, hashed_password, full_name, email, phone_number, address, original_image, role, is_verified_email)
VALUES 
('nguyenhoa', '$2a$10$tLnDg/6/QNu/nD3bIcoR2OtUqNUci4jkzlswN6cHRxhJ4QuEOvXHW', 'Nguyễn Văn Hòa', 'hoa.nguyen@example.com', '0123456789', '123 Đường ABC, Hà Nội', 'hoa_avatar.png', 'user', true),
('trangnguyen', 'hashed_password_example', 'Trần Thị Trang', 'trang.nguyen@example.com', '0987654321', '456 Đường XYZ, TP HCM',  'trang_avatar.png', 'admin', false);

INSERT INTO Pet (name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date,original_image)
VALUES 
('Milo', 'Chó', 'Poodle', 3, 'Đực', 'Dị ứng nhẹ với phấn hoa', 5.2, '2020-04-15', 'nguyenhoa', 'MICRO123456', '2023-10-01',  'milo.png'),
('Luna', 'Mèo', 'Anh lông ngắn', 2, 'Cái', 'Không có vấn đề sức khỏe nghiêm trọng', 3.8, '2021-06-10', 'trangnguyen', 'MICRO654321', '2023-09-20','luna.png');

INSERT INTO Vaccination (petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes)
VALUES 
(1, 'Rabies Vaccine', '2023-01-15', '2024-01-15', 'Trung tâm thú y Hà Nội', 'BATCH001', 'Tiêm phòng bệnh dại định kỳ cho chó'),
(2, 'Feline Leukemia Vaccine', '2023-05-10', '2024-05-10', 'Trung tâm thú y TP HCM', 'BATCH002', 'Tiêm phòng bệnh bạch cầu cho mèo');



INSERT INTO ServiceType (serviceTypeName, description, iconURL)
VALUES 
('Khám sức khỏe', 'Dịch vụ khám sức khỏe định kỳ cho thú cưng', 'icon_health_check.png'),
('Tiêm phòng', 'Dịch vụ tiêm phòng cho thú cưng', 'icon_vaccination.png');

INSERT INTO Service (typeID, name, price, duration, description)
VALUES 
(1, 'Khám sức khỏe tổng quát', 300000, '00:30:00', 'Khám sức khỏe và kiểm tra tổng quát tình trạng của thú cưng'),
(2, 'Tiêm phòng bệnh dại', 150000, '00:15:00', 'Tiêm phòng bệnh dại cho thú cưng');

INSERT INTO Appointment (petid, doctor_id, service_id, date, status, notes, reminder_send, time_slot_id)
VALUES 
(1, 1, 1, '2023-10-12 09:00:00', 'Scheduled', 'Lần khám sức khỏe định kỳ', false, 1),
(2, 2, 2, '2023-10-15 14:00:00', 'Completed', 'Tiêm phòng bệnh bạch cầu cho mèo Luna', true, 2);

INSERT INTO Doctors (user_id, specialization, years_of_experience, education, certificate_number, bio, consultation_fee)
VALUES 
(1, 'Bác sĩ thú y chuyên khoa da liễu', 10, 'Đại học Nông Lâm Hà Nội', 'CERT12345', 'Chuyên gia chăm sóc da và điều trị nấm da cho thú cưng', 500000),
(2, 'Bác sĩ thú y chuyên khoa nội tiết', 8, 'Đại học Y Dược TP HCM', 'CERT67890', 'Chuyên gia điều trị các bệnh về nội tiết cho thú cưng', 450000);

INSERT INTO TimeSlots (doctor_id, start_time, end_time, is_active, day)
VALUES 
(1, '2023-10-12 09:00:00', '2023-10-12 09:30:00', true, '2023-10-12'),
(2, '2023-10-15 14:00:00', '2023-10-15 14:15:00', true, '2023-10-15');

INSERT INTO notifications (
  username, title, description,datetime)
VALUES 
('nguyenhoa', 'Nhắc lịch uống thuốc', 'Milo cần uống thuốc đúng giờ để điều trị bệnh dị ứng', '2023-10-12 08:00:00'),
( 'nguyenhoa','Lịch tiêm phòng định kỳ', 'Luna cần tiêm phòng bạch cầu vào tháng tới', '2023-11-15 09:00:00');


-- Sample data
INSERT INTO diseases (name, description, symptoms)
VALUES (
    'Nấm da',
    'Bệnh nấm da là một bệnh phổ biến ở thú cưng, đặc biệt là chó và mèo',
    '["Ngứa nhiều", "Da đỏ", "Rụng lông từng mảng", "Vảy da"]'
);

INSERT INTO medicines (name, description, usage, dosage, frequency, duration, side_effects)
VALUES
(
    'Ketoconazole',
    'Thuốc kháng nấm dạng uống',
    'Uống sau khi ăn',
    '5-10mg/kg thể trọng',
    '1 lần/ngày',
    '2-4 tuần',
    'Có thể gây buồn nôn, chán ăn'
),
(
    'Miconazole',
    'Thuốc kháng nấm dạng bôi',
    'Bôi trực tiếp lên vùng da bị nấm',
    'Bôi một lớp mỏng',
    '2 lần/ngày',
    '2-4 tuần',
    'Có thể gây kích ứng da nhẹ'
);

-- Link diseases with medicines
INSERT INTO disease_medicines (disease_id, medicine_id)
VALUES
(1, 1),
(1, 2);

-- Insert sample data
INSERT INTO treatment_phases (disease_id, phase_number, phase_name, description, duration, notes)
VALUES 
(1, 1, 'Giai đoạn cấp tính', 'Điều trị ban đầu để kiểm soát các triệu chứng', '1-2 tuần', 'Cần theo dõi sát trong giai đoạn này'),
(1, 2, 'Giai đoạn duy trì', 'Tiếp tục điều trị để ngăn ngừa tái phát', '2-4 tuần', 'Có thể điều chỉnh liều dựa trên đáp ứng');

INSERT INTO phase_medicines (phase_id, medicine_id, dosage, frequency, duration, notes)
VALUES 
(1, 1, '10mg/kg', '2 lần/ngày', '1 tuần', 'Uống sau bữa ăn'),
(1, 2, 'Bôi lớp mỏng', '3 lần/ngày', '1 tuần', 'Tránh để thú cưng liếm thuốc'),
(2, 1, '5mg/kg', '1 lần/ngày', '3 tuần', 'Uống sau bữa ăn'),
(2, 2, 'Bôi lớp mỏng', '2 lần/ngày', '3 tuần', 'Tiếp tục theo dõi phản ứng của da');