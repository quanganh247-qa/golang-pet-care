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
  secret_code varchar NOT NULL,
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
    end_type VARCHAR(50),
    end_date DATE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
CREATE TABLE Notification (
  notificationID BIGSERIAL PRIMARY KEY,
  petID BIGINT,
  title VARCHAR(100) NOT NULL,
  body TEXT,
  dueDate TIMESTAMP NOT NULL,
  repeatInterval VARCHAR(50),
  isCompleted BOOLEAN DEFAULT false,
  notificationSent BOOLEAN DEFAULT false
);

-- Modify notification_history table to include username
CREATE TABLE NotificationHistory (
    id BIGSERIAL PRIMARY KEY,
    notification_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    device_token_id BIGINT,
    title VARCHAR(255) NOT NULL,
    body TEXT,
    data JSONB,
    sent_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    delivered_at TIMESTAMPTZ,
    opened_at TIMESTAMPTZ,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (notification_id) REFERENCES Notification(notificationID) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (device_token_id) REFERENCES DeviceTokens(id) ON DELETE SET NULL
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

-- Create indexes
CREATE INDEX idx_diseases_name ON diseases(name);
CREATE INDEX idx_medicines_name ON medicines(name);

ALTER TABLE Pet ADD CONSTRAINT pet_users_fk FOREIGN KEY (username) REFERENCES users (username);

ALTER TABLE Vaccination ADD CONSTRAINT vaccination_pet_fk FOREIGN KEY (petID) REFERENCES Pet (petid);

ALTER TABLE Notification ADD CONSTRAINT not_pet_fk FOREIGN KEY (petID) REFERENCES Pet (petid);

ALTER TABLE Service ADD CONSTRAINT service_type_fk FOREIGN KEY (typeID) REFERENCES ServiceType (typeID);

ALTER TABLE Appointment ADD CONSTRAINT appointment_pet_fk FOREIGN KEY (petid) REFERENCES Pet (petid);

ALTER TABLE Appointment ADD CONSTRAINT appointment_service_fk FOREIGN KEY (service_id) REFERENCES Service (serviceID);

ALTER TABLE CheckoutService ADD CONSTRAINT cs_checkout_fk FOREIGN KEY (checkoutID) REFERENCES Checkout (checkout_id);

ALTER TABLE CheckoutService ADD CONSTRAINT cs_service_fk FOREIGN KEY (serviceID) REFERENCES Service (serviceID);

ALTER TABLE Doctors ADD CONSTRAINT fk_doctor_user FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE DoctorSchedules ADD CONSTRAINT fk_schedule_doctor FOREIGN KEY (doctor_id) REFERENCES Doctors (id);







-- -- 1. Query cơ bản để lấy thông tin bệnh và thuốc điều trị
-- SELECT 
--     d.id AS disease_id,
--     d.name AS disease_name,
--     d.description AS disease_description,
--     d.symptoms,
--     m.id AS medicine_id,
--     m.name AS medicine_name,
--     m.usage AS medicine_usage,
--     m.dosage,
--     m.frequency,
--     m.duration,
--     m.side_effects
-- FROM diseases d
-- LEFT JOIN disease_medicines dm ON d.id = dm.disease_id
-- LEFT JOIN medicines m ON dm.medicine_id = m.id
-- WHERE LOWER(d.name) LIKE LOWER('%nấm da%');

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



-- -- Query lấy phác đồ điều trị đầy đủ
-- SELECT 
--     d.name AS disease_name,
--     d.description AS disease_description,
--     d.symptoms,
--     tp.phase_number,
--     tp.phase_name,
--     tp.description AS phase_description,
--     tp.duration AS phase_duration,
--     tp.notes AS phase_notes,
--     m.name AS medicine_name,
--     m.description AS medicine_description,
--     COALESCE(pm.dosage, m.dosage) AS dosage,
--     COALESCE(pm.frequency, m.frequency) AS frequency,
--     COALESCE(pm.duration, m.duration) AS duration,
--     m.side_effects,
--     pm.notes AS medicine_notes
-- FROM diseases d
-- JOIN treatment_phases tp ON d.id = tp.disease_id
-- JOIN phase_medicines pm ON tp.id = pm.phase_id
-- JOIN medicines m ON pm.medicine_id = m.id
-- WHERE LOWER(d.name) LIKE LOWER($1)
-- ORDER BY tp.phase_number, m.name;

-- -- 3. Query để lấy tổng quan điều trị
-- SELECT 
--     d.name AS disease_name,
--     d.description,
--     d.symptoms,
--     json_agg(
--         json_build_object(
--             'phase_number', tp.phase_number,
--             'phase_name', tp.phase_name,
--             'duration', tp.duration,
--             'medicines', (
--                 SELECT json_agg(
--                     json_build_object(
--                         'name', m.name,
--                         'dosage', COALESCE(pm.dosage, m.dosage),
--                         'frequency', COALESCE(pm.frequency, m.frequency),
--                         'duration', COALESCE(pm.duration, m.duration),
--                         'notes', pm.notes
--                     )
--                 )
--                 FROM phase_medicines pm
--                 JOIN medicines m ON pm.medicine_id = m.id
--                 WHERE pm.phase_id = tp.id
--             )
--         )
--     ) AS treatment_phases
-- FROM diseases d
-- JOIN treatment_phases tp ON d.id = tp.disease_id
-- WHERE LOWER(d.name) LIKE LOWER($1)
-- GROUP BY d.id, d.name, d.description, d.symptoms;

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

-- -- Query lấy lịch sử điều trị
-- SELECT 
--     p.name AS pet_name,
--     d.name AS disease_name,
--     pt.start_date,
--     pt.end_date,
--     pt.status,
--     json_agg(
--         json_build_object(
--             'phase_name', tp.phase_name,
--             'start_date', tpr.start_date,
--             'end_date', tpr.end_date,
--             'status', tpr.status,
--             'notes', tpr.notes
--         )
--     ) AS progress
-- FROM pet_treatments pt
-- JOIN Pets p ON pt.pet_id = p.id
-- JOIN diseases d ON pt.disease_id = d.id
-- JOIN treatment_progress tpr ON pt.id = tpr.treatment_id
-- JOIN treatment_phases tp ON tpr.phase_id = tp.id
-- WHERE p.id = $1
-- GROUP BY p.id, p.name, d.name, pt.start_date, pt.end_date, pt.status
-- ORDER BY pt.start_date DESC;


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

-- INSERT INTO pet_schedule (pet_id, schedule_type, event_time, food_type, quantity, frequency, notes, is_active)
-- VALUES
-- -- Feeding schedule for pet with ID 1
-- (1, 'Feeding', '2023-11-15 08:00:00', 'Hạt khô cao cấp', 100, '2 lần/ngày', 'Cho ăn vào buổi sáng, cung cấp thêm nước', true),
-- (1, 'Feeding', '2023-11-15 18:00:00', 'Thịt gà nấu chín', 150, 'Hàng ngày', 'Cho ăn vào buổi tối, tránh thức ăn quá lạnh', true),

-- -- Feeding schedule for pet with ID 2
-- (2, 'Feeding', '2023-11-15 07:00:00', 'Thức ăn hạt dành cho mèo', 80, 'Hàng ngày', 'Phù hợp cho mèo con', true),
-- (2, 'Feeding', '2023-11-15 19:00:00', 'Pate thịt', 120, 'Hàng ngày', 'Phục vụ buổi tối', true);

-- INSERT INTO pet_schedule (pet_id, schedule_type, event_time, duration, activity_type, frequency, notes, is_active)
-- VALUES
-- -- Activity schedule for pet with ID 1
-- (1, 'Activity', '2023-11-15 09:00:00', '00:30:00', 'Đi dạo', 'Hàng ngày', 'Dạo quanh công viên', true),
-- (1, 'Activity', '2023-11-15 17:00:00', '01:00:00', 'Chơi đùa', 'Hàng ngày', 'Chơi với bóng và đồ chơi', true),

-- -- Activity schedule for pet with ID 2
-- (2, 'Activity', '2023-11-15 10:00:00', '00:45:00', 'Tập thể dục', 'Hàng tuần', 'Bài tập nhẹ nhàng để duy trì sức khỏe', true),
-- (2, 'Activity', '2023-11-15 16:00:00', '00:30:00', 'Chơi với đồ chơi', 'Hàng ngày', 'Chơi với đồ chơi mèo', true);


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

INSERT INTO Notification (petID, title, body, dueDate, repeatInterval, isCompleted, notificationSent)
VALUES 
(1, 'Nhắc lịch uống thuốc', 'Milo cần uống thuốc đúng giờ để điều trị bệnh dị ứng', '2023-10-12 08:00:00', 'Hàng ngày', false, false),
(2, 'Lịch tiêm phòng định kỳ', 'Luna cần tiêm phòng bạch cầu vào tháng tới', '2023-11-15 09:00:00', 'Hàng năm', false, false);


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