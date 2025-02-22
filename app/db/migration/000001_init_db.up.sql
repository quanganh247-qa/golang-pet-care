-- 1. Core User Management Tables
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

CREATE TABLE device_tokens (
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

CREATE TABLE notifications (
  notificationID BIGSERIAL PRIMARY KEY,
  username varchar NOT NULL,
  title VARCHAR(100) NOT NULL,
  description TEXT,
  datetime TIMESTAMP NOT NULL,
  is_read BOOLEAN DEFAULT false
);

-- 2. Medical Staff Tables
CREATE TABLE doctors (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  specialization VARCHAR(100),
  years_of_experience INT,
  education TEXT,
  certificate_number VARCHAR(50),
  bio TEXT,
  consultation_fee FLOAT8
);

CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE time_slots (
    id BIGSERIAL PRIMARY KEY,
    doctor_id INT NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    max_patients INT NULL,
    booked_patients INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES doctors(id),
    CONSTRAINT unique_slot UNIQUE (doctor_id, date, start_time)
);

-- 3. Pet Management Tables
CREATE TABLE pets (
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

CREATE TABLE pet_logs (
    log_id BIGSERIAL PRIMARY KEY,
	petid int8 NOT NULL,
	datetime timestamp NULL,
	title varchar NULL,
	notes text NULL,
	CONSTRAINT newtable_pet_fk FOREIGN KEY (petid) REFERENCES pets(petid)
);

CREATE TABLE vaccinations (
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
    pet_id BIGINT REFERENCES pets(petid),
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

-- 4. Medical Records & Treatment Tables
CREATE TABLE medical_records (
	id bigserial NOT NULL,
	pet_id bigint NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT medical_records_pk PRIMARY KEY (id)
);

CREATE TABLE medical_history (
	id bigserial NOT NULL,
	medical_record_id bigint NULL,
	"condition" varchar NULL,
	diagnosis_date timestamp NULL,
	notes text NULL,
  treatment int8 NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT medical_history_pk PRIMARY KEY (id)
);

CREATE TABLE allergies (
	id bigserial NOT NULL,
	medical_record_id bigint NULL,
	allergen jsonb NULL,
	severity varchar NULL,
	reaction jsonb NULL,
	notes text NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT allergies_pk PRIMARY KEY (id)
);

CREATE TABLE diseases (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    symptoms JSONB, -- Store symptoms as JSON array
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE medicines (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    usage TEXT,
    dosage TEXT,
    frequency TEXT,
    duration TEXT,
    side_effects TEXT,
    medical_record_id int8 NULL,
	prescribing_vet varchar NULL,
	start_date date NULL,
	end_date date NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE disease_medicines (
    disease_id BIGINT REFERENCES diseases(id),
    medicine_id BIGINT REFERENCES medicines(id),
    PRIMARY KEY (disease_id, medicine_id)
);

CREATE TABLE pet_treatments (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT REFERENCES pets(petid),
    disease_id BIGINT REFERENCES diseases(id),
    start_date DATE,
    end_date DATE,
    status VARCHAR(50),  -- CHECK (status IN ('ongoing', 'completed', 'paused', 'cancelled')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE treatment_phases (
    id BIGSERIAL PRIMARY KEY,
    treatment_id BIGINT REFERENCES pet_treatments(id),
    phase_name VARCHAR(255),
    description TEXT,
    status VARCHAR(50), -- CHECK (status IN ('pending', 'active', 'completed')),
    start_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE phase_medicines (
    phase_id BIGINT REFERENCES treatment_phases(id),
    medicine_id BIGINT REFERENCES medicines(id),
    dosage TEXT,
    frequency TEXT,
    duration TEXT,
    notes TEXT,
    created_at timestamptz NULL,
    PRIMARY KEY (phase_id, medicine_id)
);

-- 5. Services and Appointments
CREATE TABLE services (
	id bigserial PRIMARY KEY,
	"name" varchar(255) NULL,
	description text NULL,
	duration int2 NULL,
	"cost" float8 NULL,
	category varchar(255) NULL,
	notes text NULL,
	created_at timestamp DEFAULT now() NULL,
	updated_at timestamp NULL
);

CREATE TABLE appointments (
  appointment_id BIGSERIAL PRIMARY KEY,
  petid BIGINT,
  username VARCHAR,
  doctor_id BIGINT,
  service_id BIGINT,
  date timestamp DEFAULT (now()),
  notes TEXT,
  reminder_send BOOLEAN DEFAULT false,
  time_slot_id BIGINT,
  payment_status VARCHAR(20),
  created_at timestamp DEFAULT (now())
);

-- 6. Billing and Payments
CREATE TABLE checkouts (
  checkout_id BIGSERIAL PRIMARY KEY,
  petid BIGINT,
  doctor_id BIGINT,
  date timestamp DEFAULT (now()),
  total_tmount float8 NOT NULL,
  payment_status VARCHAR(20),
  payment_method VARCHAR(50),
  notes TEXT
);

CREATE TABLE checkout_services (
  checkoutService_ID BIGSERIAL PRIMARY KEY,
  checkoutID BIGINT,
  serviceID BIGINT,
  quantity INT DEFAULT 1,
  unitPrice float8,
  subtotal float8
);

-- 7. E-commerce Tables
CREATE TABLE products (
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

CREATE TABLE carts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cart_items (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT REFERENCES carts(id),
    product_id BIGINT,
    quantity INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
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

-- 8. Foreign Key Constraints
ALTER TABLE pets ADD CONSTRAINT pet_users_fk FOREIGN KEY (username) REFERENCES users (username);

ALTER TABLE vaccinations ADD CONSTRAINT vaccination_pet_fk FOREIGN KEY (petID) REFERENCES pets (petid);

ALTER TABLE appointments ADD CONSTRAINT appointment_pet_fk FOREIGN KEY (petid) REFERENCES pets (petid);

ALTER TABLE appointments ADD CONSTRAINT appointment_service_fk FOREIGN KEY (service_id) REFERENCES services (id);

ALTER TABLE checkout_services ADD CONSTRAINT cs_checkout_fk FOREIGN KEY (checkoutID) REFERENCES checkouts (checkout_id);

ALTER TABLE checkout_services ADD CONSTRAINT cs_service_fk FOREIGN KEY (serviceID) REFERENCES services (id);

ALTER TABLE doctors ADD CONSTRAINT fk_doctor_user FOREIGN KEY (user_id) REFERENCES users (id);
