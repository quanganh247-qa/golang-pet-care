<<<<<<< HEAD



-- public.checkouts definition

-- Drop table

-- DROP TABLE public.checkouts;

CREATE TABLE public.checkouts (
	checkout_id bigserial NOT NULL,
	petid int8 NULL,
	doctor_id int8 NULL,
	"date" timestamp DEFAULT now() NULL,
	total_tmount float8 NOT NULL,
	payment_status varchar(20) NULL,
	payment_method varchar(50) NULL,
	notes text NULL,
	CONSTRAINT checkouts_pkey PRIMARY KEY (checkout_id)
);


-- public.diseases definition

-- Drop table

-- DROP TABLE public.diseases;

CREATE TABLE public.diseases (
	id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	description text NULL,
	symptoms jsonb NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT diseases_pkey PRIMARY KEY (id)
);


-- public.medical_history definition

-- Drop table

-- DROP TABLE public.medical_history;

CREATE TABLE public.medical_history (
	id bigserial NOT NULL,
	medical_record_id int8 NULL,
	"condition" varchar NULL,
	diagnosis_date timestamp NULL,
	notes text NULL,
	treatment int8 NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT medical_history_pk PRIMARY KEY (id)
);


-- public.medical_records definition

-- Drop table

-- DROP TABLE public.medical_records;

CREATE TABLE public.medical_records (
	id bigserial NOT NULL,
	pet_id int8 NULL,
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT medical_records_pk PRIMARY KEY (id)
);


-- public.medicines definition

-- Drop table

-- DROP TABLE public.medicines;

CREATE TABLE public.medicines (
	id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	description text NULL,
	"usage" text NULL,
	dosage text NULL,
	frequency text NULL,
	duration text NULL,
	side_effects text NULL,
	start_date date NULL,
	end_date date NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	expiration_date date NULL,
	quantity int8 NULL,
	CONSTRAINT medicines_pkey PRIMARY KEY (id)
);


-- public.products definition

-- Drop table

-- DROP TABLE public.products;

CREATE TABLE public.products (
	product_id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	description text NULL,
	price float8 NOT NULL,
	stock_quantity int4 DEFAULT 0 NULL,
	category varchar(100) NULL,
	data_image bytea NULL,
	original_image varchar(255) NULL,
	created_at timestamp DEFAULT now() NULL,
	is_available bool DEFAULT true NULL,
	removed_at timestamp NULL,
	CONSTRAINT products_pkey PRIMARY KEY (product_id)
);


-- public.services definition

-- Drop table

-- DROP TABLE public.services;

CREATE TABLE public.services (
	id bigserial NOT NULL,
	"name" varchar(255) NULL,
	description text NULL,
	duration int2 NULL,
	"cost" float8 NULL,
	category varchar(255) NULL,
	priority int2 DEFAULT 1 NULL,
	created_at timestamp DEFAULT now() NULL,
	CONSTRAINT services_pkey PRIMARY KEY (id)
);


-- public.states definition

-- Drop table

-- DROP TABLE public.states;

CREATE TABLE public.states (
	id bigserial NOT NULL,
	state varchar(255) NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT states_pkey PRIMARY KEY (id),
	CONSTRAINT states_state_check CHECK (((state)::text = ANY ((ARRAY['Scheduled'::character varying, 'Confirmed'::character varying, 'Checked In'::character varying, 'Waiting'::character varying, 'In Progress'::character varying, 'Completed'::character varying, 'Closed'::character varying])::text[])))
);

-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	hashed_password varchar NOT NULL,
	full_name varchar NOT NULL,
	email varchar NOT NULL,
	phone_number varchar NULL,
	address varchar NULL,
	data_image bytea NULL,
	original_image varchar(255) NULL,
	"role" varchar(20) NULL,
	status varchar(20) NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	is_verified_email bool DEFAULT false NULL,
	removed_at timestamp NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
);


-- public.verify_emails definition

-- Drop table

-- DROP TABLE public.verify_emails;

CREATE TABLE public.verify_emails (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	email varchar NOT NULL,
	secret_code int8 NOT NULL,
	is_used bool DEFAULT false NOT NULL,
	created_at timestamp DEFAULT now() NULL,
	expired_at timestamp DEFAULT (now() + '00:15:00'::interval) NULL,
	CONSTRAINT verify_emails_pkey PRIMARY KEY (id)
);


-- public.carts definition

-- Drop table

-- DROP TABLE public.carts;

CREATE TABLE public.carts (
	id bigserial NOT NULL,
	user_id int8 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT carts_pkey PRIMARY KEY (id),
	CONSTRAINT carts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);


-- public.checkout_services definition

-- Drop table

-- DROP TABLE public.checkout_services;

CREATE TABLE public.checkout_services (
	checkoutservice_id bigserial NOT NULL,
	checkoutid int8 NULL,
	serviceid int8 NULL,
	quantity int4 DEFAULT 1 NULL,
	unitprice float8 NULL,
	subtotal float8 NULL,
	CONSTRAINT checkout_services_pkey PRIMARY KEY (checkoutservice_id),
	CONSTRAINT cs_checkout_fk FOREIGN KEY (checkoutid) REFERENCES public.checkouts(checkout_id),
	CONSTRAINT cs_service_fk FOREIGN KEY (serviceid) REFERENCES public.services(id)
);


-- public.device_tokens definition

-- Drop table

-- DROP TABLE public.device_tokens;

CREATE TABLE public.device_tokens (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	"token" varchar NOT NULL,
	device_type varchar(50) NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	last_used_at timestamp NULL,
	expired_at timestamp NULL,
	CONSTRAINT device_tokens_pkey PRIMARY KEY (id),
	CONSTRAINT device_tokens_token_key UNIQUE (token),
	CONSTRAINT fk_device_tokens_username FOREIGN KEY (username) REFERENCES public.users(username) ON DELETE CASCADE
);


-- public.doctors definition

-- Drop table

-- DROP TABLE public.doctors;

CREATE TABLE public.doctors (
	id bigserial NOT NULL,
	user_id int8 NOT NULL,
	specialization varchar(100) NULL,
	years_of_experience int4 NULL,
	education text NULL,
	certificate_number varchar(50) NULL,
	bio text NULL,
	CONSTRAINT doctors_pkey PRIMARY KEY (id),
	CONSTRAINT fk_doctor_user FOREIGN KEY (user_id) REFERENCES public.users(id)
);


-- public.files definition

-- Drop table

-- DROP TABLE public.files;

CREATE TABLE public.files (
	id bigserial NOT NULL,
	file_name varchar(255) NOT NULL,
	file_path varchar(255) NOT NULL,
	file_size int8 NOT NULL,
	file_type varchar(50) NOT NULL,
	uploaded_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	user_id int8 NULL,
	CONSTRAINT files_pkey PRIMARY KEY (id),
	CONSTRAINT files_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL
);


-- public.notifications definition

-- Drop table

-- DROP TABLE public.notifications;

CREATE TABLE public.notifications (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	title varchar(100) NOT NULL,
	"content" text NULL,
	is_read bool DEFAULT false NULL,
	related_id int4 NULL,
	related_type varchar(255) NULL,
	datetime timestamp NOT NULL,
	notify_type varchar(255) NULL,
	CONSTRAINT notifications_pkey PRIMARY KEY (id),
	CONSTRAINT notifications_username_fkey FOREIGN KEY (username) REFERENCES public.users(username) ON DELETE CASCADE
);


-- public.orders definition

-- Drop table

-- DROP TABLE public.orders;

CREATE TABLE public.orders (
	id bigserial NOT NULL,
	user_id int8 NOT NULL,
	order_date timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	total_amount float8 NOT NULL,
	payment_status varchar(20) DEFAULT 'pending'::character varying NULL,
	cart_items jsonb NULL,
	shipping_address varchar(255) NULL,
	notes text NULL,
	CONSTRAINT orders_pkey PRIMARY KEY (id),
	CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE
);


-- public.pets definition

-- Drop table

-- DROP TABLE public.pets;

CREATE TABLE public.pets (
	petid bigserial NOT NULL,
	"name" varchar(100) NOT NULL,
	"type" varchar(50) NOT NULL,
	breed varchar(100) NULL,
	age int4 NULL,
	gender varchar(10) NULL,
	healthnotes text NULL,
	weight float8 NULL,
	birth_date date NULL,
	username varchar NOT NULL,
	microchip_number varchar(50) NULL,
	last_checkup_date date NULL,
	is_active bool DEFAULT true NULL,
	data_image bytea NULL,
	original_image varchar(255) NULL,
	CONSTRAINT pets_pkey PRIMARY KEY (petid),
	CONSTRAINT pet_users_fk FOREIGN KEY (username) REFERENCES public.users(username)
);


-- public.vaccinations definition

-- Drop table

-- DROP TABLE public.vaccinations;

CREATE TABLE public.vaccinations (
	vaccinationid bigserial NOT NULL,
	petid int8 NULL,
	vaccinename varchar(100) NOT NULL,
	dateadministered timestamp NOT NULL,
	nextduedate timestamp NULL,
	vaccineprovider varchar(100) NULL,
	batchnumber varchar(50) NULL,
	notes text NULL,
	CONSTRAINT vaccinations_pkey PRIMARY KEY (vaccinationid),
	CONSTRAINT vaccination_pet_fk FOREIGN KEY (petid) REFERENCES public.pets(petid)
);



-- public.cart_items definition

-- Drop table

-- DROP TABLE public.cart_items;

CREATE TABLE public.cart_items (
	id bigserial NOT NULL,
	cart_id int8 NOT NULL,
	product_id int8 NOT NULL,
	quantity int4 DEFAULT 1 NULL,
	unit_price float8 NOT NULL,
	total_price float8 GENERATED ALWAYS AS (quantity::double precision * unit_price) STORED NULL,
	CONSTRAINT cart_items_pkey PRIMARY KEY (id),
	CONSTRAINT cartitem_cart_id_product_id_unique UNIQUE (cart_id, product_id),
	CONSTRAINT cart_items_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES public.carts(id) ON DELETE CASCADE,
	CONSTRAINT cart_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(product_id) ON DELETE CASCADE
);


-- public.consultations definition

-- Drop table

-- DROP TABLE public.consultations;

CREATE TABLE public.consultations (
	id serial4 NOT NULL,
	appointment_id int8 NULL,
	subjective text NULL,
	objective text NULL,
	assessment text NULL,
	plan int8 NULL,
	created_at timestamp DEFAULT now() NULL,
	CONSTRAINT consultations_pkey PRIMARY KEY (id),
	CONSTRAINT consultations_appointment_id_fkey FOREIGN KEY (appointment_id) REFERENCES public.appointments(appointment_id) ON DELETE CASCADE
);


-- public.pet_logs definition

-- Drop table

-- DROP TABLE public.pet_logs;

CREATE TABLE public.pet_logs (
	log_id bigserial NOT NULL,
	petid int8 NOT NULL,
	datetime timestamp NULL,
	title varchar NULL,
	notes text NULL,
	CONSTRAINT pet_logs_pkey PRIMARY KEY (log_id),
	CONSTRAINT newtable_pet_fk FOREIGN KEY (petid) REFERENCES public.pets(petid)
);


-- public.pet_schedule definition

-- Drop table

-- DROP TABLE public.pet_schedule;

CREATE TABLE public.pet_schedule (
	id bigserial NOT NULL,
	pet_id int8 NULL,
	title varchar(255) NULL,
	reminder_datetime timestamp NULL,
	event_repeat varchar(50) NULL,
	end_type bool DEFAULT false NULL,
	end_date date NULL,
	notes text NULL,
	is_active bool NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	removedat timestamp NULL,
	CONSTRAINT pet_schedule_pkey PRIMARY KEY (id),
	CONSTRAINT pet_schedule_pet_id_fkey FOREIGN KEY (pet_id) REFERENCES public.pets(petid)
);


-- public.pet_treatments definition

-- Drop table

-- DROP TABLE public.pet_treatments;

CREATE TABLE public.pet_treatments (
	id bigserial NOT NULL,
	pet_id int8 NULL,
	disease_id int8 NULL,
	start_date date NULL,
	end_date date NULL,
	status varchar(50) NULL,
	"name" varchar NULL,
	"type" varchar NULL,
	description text NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	doctor_id int4 NULL,
	CONSTRAINT pet_treatments_pkey PRIMARY KEY (id),
	CONSTRAINT pet_treatments_disease_id_fkey FOREIGN KEY (disease_id) REFERENCES public.diseases(id),
	CONSTRAINT pet_treatments_pet_id_fkey FOREIGN KEY (pet_id) REFERENCES public.pets(petid)
);


-- public.treatment_phases definition

-- Drop table

-- DROP TABLE public.treatment_phases;

CREATE TABLE public.treatment_phases (
	id bigserial NOT NULL,
	treatment_id int8 NULL,
	phase_name varchar(255) NULL,
	description text NULL,
	status varchar(50) NULL,
	start_date date NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT treatment_phases_pkey PRIMARY KEY (id),
	CONSTRAINT treatment_phases_treatment_id_fkey FOREIGN KEY (treatment_id) REFERENCES public.pet_treatments(id)
);


-- public.phase_medicines definition

-- Drop table

-- DROP TABLE public.phase_medicines;

CREATE TABLE public.phase_medicines (
	phase_id int8 NOT NULL,
	medicine_id int8 NOT NULL,
	dosage text NULL,
	frequency text NULL,
	duration text NULL,
	notes text NULL,
	created_at timestamptz NULL,
	quantity int4 NULL,
	is_received bool DEFAULT true NULL,
	CONSTRAINT phase_medicines_pkey PRIMARY KEY (phase_id, medicine_id),
	CONSTRAINT phase_medicines_medicine_id_fkey FOREIGN KEY (medicine_id) REFERENCES public.medicines(id),
	CONSTRAINT phase_medicines_phase_id_fkey FOREIGN KEY (phase_id) REFERENCES public.treatment_phases(id)
);

-- public.clinics definition

-- Drop table

-- DROP TABLE public.clinics;

CREATE TABLE public.clinics (
	id bigserial NOT NULL,
	"name" varchar NULL,
	address text NULL,
	phone varchar NULL,
	CONSTRAINT clinics_pk PRIMARY KEY (id)
);


-- public.shifts definition

-- Drop table

-- DROP TABLE public.shifts;

CREATE TABLE public.shifts (
    id bigserial NOT NULL,
    doctor_id int8 NOT NULL,
    start_time timestamp NOT NULL,
    end_time timestamp NOT NULL,
	max_patients int4 DEFAULT 10 NULL,
    assigned_patients int4 DEFAULT 0 NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT shifts_pkey PRIMARY KEY (id),
    CONSTRAINT shifts_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id)
);
-- public.time_slots definition

-- Drop table

-- DROP TABLE public.time_slots;

CREATE TABLE public.time_slots (
	id bigserial NOT NULL,
	doctor_id int4 NOT NULL,
	"date" date NOT NULL,
	start_time time NOT NULL,
	end_time time NOT NULL,
	max_patients int4 NULL,
	booked_patients int4 DEFAULT 0 NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	shift_id bigserial NOT NULL,
	CONSTRAINT time_slots_pkey PRIMARY KEY (id),
	CONSTRAINT unique_slot UNIQUE (doctor_id, date, start_time)
);


-- public.time_slots foreign keys

ALTER TABLE public.time_slots ADD CONSTRAINT time_slots_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id);
ALTER TABLE public.time_slots ADD CONSTRAINT time_slots_shift_id_fkey FOREIGN KEY (shift_id) REFERENCES public.shifts(id);


-- -- public.appointments definition

-- -- Drop table

-- -- DROP TABLE public.appointments;

CREATE TABLE public.appointments (
    appointment_id bigserial NOT NULL,
    petid int8 NULL,
    username varchar NULL,
    doctor_id int8 NULL,
    service_id int8 NULL,
    "date" timestamp DEFAULT now() NULL,
    reminder_send bool DEFAULT false NULL,
    time_slot_id int8 NULL,
    created_at timestamp DEFAULT now() NULL,
    state_id int4 NULL,
    appointment_reason text NULL,
    priority varchar(20) NULL,
    arrival_time timestamp NULL,
    room_id int8 NULL,
    confirmation_sent bool DEFAULT false NULL,
	notes text NULL,
	updated_at timestamp null,
    CONSTRAINT appointments_pkey PRIMARY KEY (appointment_id),
    CONSTRAINT appointment_pet_fk FOREIGN KEY (petid) REFERENCES public.pets(petid),
    CONSTRAINT appointment_service_fk FOREIGN KEY (service_id) REFERENCES public.services(id),
    CONSTRAINT appointments_state_id_fkey FOREIGN KEY (state_id) REFERENCES public.states(id),
    CONSTRAINT appointments_room_id_fkey FOREIGN KEY (room_id) REFERENCES public.rooms(id)
);


CREATE TABLE public.rooms (
    id bigserial NOT NULL,
    name varchar(100) NOT NULL,
    type varchar(50) NOT NULL,
    status varchar(20) DEFAULT 'available' NULL,
    current_appointment_id int8 NULL,
    available_at timestamp NULL,
    CONSTRAINT rooms_pkey PRIMARY KEY (id),
    CONSTRAINT rooms_current_appointment_id_fkey 
        FOREIGN KEY (current_appointment_id) REFERENCES public.appointments(appointment_id)
);

CREATE TABLE public.pet_allergies (
	id bigserial NOT NULL,
	pet_id int8 NULL,
	"type" varchar NULL,
	detail text NULL,
	CONSTRAINT pet_alert_pk PRIMARY KEY (id)
);


-- Indexing for public.checkouts
CREATE INDEX idx_checkouts_petid ON public.checkouts (petid);
CREATE INDEX idx_checkouts_doctor_id ON public.checkouts (doctor_id);
CREATE INDEX idx_checkouts_date ON public.checkouts (date);

-- Indexing for public.diseases
CREATE INDEX idx_diseases_name ON public.diseases (name);

-- Indexing for public.medical_history
CREATE INDEX idx_medical_history_medical_record_id ON public.medical_history (medical_record_id);
CREATE INDEX idx_medical_history_diagnosis_date ON public.medical_history (diagnosis_date);

-- Indexing for public.medical_records
CREATE INDEX idx_medical_records_pet_id ON public.medical_records (pet_id);

-- Indexing for public.medicines
CREATE INDEX idx_medicines_name ON public.medicines (name);
CREATE INDEX idx_medicines_expiration_date ON public.medicines (expiration_date);

-- Indexing for public.products
CREATE INDEX idx_products_name ON public.products (name);
CREATE INDEX idx_products_category ON public.products (category);
CREATE INDEX idx_products_is_available ON public.products (is_available);

-- Indexing for public.services
CREATE INDEX idx_services_name ON public.services (name);
CREATE INDEX idx_services_category ON public.services (category);

-- Indexing for public.states
CREATE INDEX idx_states_state ON public.states (state);

-- Indexing for public.users
CREATE INDEX idx_users_full_name ON public.users (full_name);
CREATE INDEX idx_users_role ON public.users (role);

-- Indexing for public.verify_emails
CREATE INDEX idx_verify_emails_username ON public.verify_emails (username);
CREATE INDEX idx_verify_emails_expired_at ON public.verify_emails (expired_at);

-- Indexing for public.carts
CREATE INDEX idx_carts_user_id ON public.carts (user_id);

-- Indexing for public.checkout_services
CREATE INDEX idx_checkout_services_checkoutid ON public.checkout_services (checkoutid);
CREATE INDEX idx_checkout_services_serviceid ON public.checkout_services (serviceid);

-- Indexing for public.device_tokens
CREATE INDEX idx_device_tokens_username ON public.device_tokens (username);

-- Indexing for public.doctors
CREATE INDEX idx_doctors_user_id ON public.doctors (user_id);
CREATE INDEX idx_doctors_specialization ON public.doctors (specialization);

-- Indexing for public.files
CREATE INDEX idx_files_user_id ON public.files (user_id);
CREATE INDEX idx_files_uploaded_at ON public.files (uploaded_at);

-- Indexing for public.notifications
CREATE INDEX idx_notifications_username ON public.notifications (username);
CREATE INDEX idx_notifications_datetime ON public.notifications (datetime);
CREATE INDEX idx_notifications_is_read ON public.notifications (is_read);

-- Indexing for public.orders
CREATE INDEX idx_orders_user_id ON public.orders (user_id);
CREATE INDEX idx_orders_order_date ON public.orders (order_date);
CREATE INDEX idx_orders_payment_status ON public.orders (payment_status);

-- Indexing for public.pets
CREATE INDEX idx_pets_username ON public.pets (username);
CREATE INDEX idx_pets_name ON public.pets (name);
CREATE INDEX idx_pets_is_active ON public.pets (is_active);

-- Indexing for public.time_slots
CREATE INDEX idx_time_slots_doctor_id ON public.time_slots (doctor_id);
CREATE INDEX idx_time_slots_date ON public.time_slots (date);
CREATE INDEX idx_time_slots_doctor_id_date ON public.time_slots (doctor_id, date);

-- Indexing for public.vaccinations
CREATE INDEX idx_vaccinations_petid ON public.vaccinations (petid);
CREATE INDEX idx_vaccinations_dateadministered ON public.vaccinations (dateadministered);
CREATE INDEX idx_vaccinations_nextduedate ON public.vaccinations (nextduedate);

-- Indexing for public.appointments
CREATE INDEX idx_appointments_petid ON public.appointments (petid);
CREATE INDEX idx_appointments_username ON public.appointments (username);
CREATE INDEX idx_appointments_doctor_id ON public.appointments (doctor_id);
CREATE INDEX idx_appointments_service_id ON public.appointments (service_id);
CREATE INDEX idx_appointments_date ON public.appointments (date);
CREATE INDEX idx_appointments_time_slot_id ON public.appointments (time_slot_id);
CREATE INDEX idx_appointments_state_id ON public.appointments (state_id);

-- Indexing for public.cart_items
CREATE INDEX idx_cart_items_cart_id ON public.cart_items (cart_id);
CREATE INDEX idx_cart_items_product_id ON public.cart_items (product_id);

-- Indexing for public.consultations
CREATE INDEX idx_consultations_appointment_id ON public.consultations (appointment_id);

-- Indexing for public.pet_logs
CREATE INDEX idx_pet_logs_petid ON public.pet_logs (petid);
CREATE INDEX idx_pet_logs_datetime ON public.pet_logs (datetime);

-- Indexing for public.pet_schedule
CREATE INDEX idx_pet_schedule_pet_id ON public.pet_schedule (pet_id);
CREATE INDEX idx_pet_schedule_reminder_datetime ON public.pet_schedule (reminder_datetime);

-- Indexing for public.pet_treatments
CREATE INDEX idx_pet_treatments_pet_id ON public.pet_treatments (pet_id);
CREATE INDEX idx_pet_treatments_disease_id ON public.pet_treatments (disease_id);
CREATE INDEX idx_pet_treatments_doctor_id ON public.pet_treatments (doctor_id);

-- Indexing for public.treatment_phases
CREATE INDEX idx_treatment_phases_treatment_id ON public.treatment_phases (treatment_id);

-- Indexing for public.phase_medicines
CREATE INDEX idx_phase_medicines_phase_id ON public.phase_medicines (phase_id);
CREATE INDEX idx_phase_medicines_medicine_id ON public.phase_medicines (medicine_id);

-- Indexing for public.clinics
CREATE INDEX idx_clinics_name ON public.clinics (name);

-- Indexing for public.shifts
CREATE INDEX idx_shifts_doctor_id ON public.shifts (doctor_id);
CREATE INDEX idx_shifts_start_time ON public.shifts (start_time);
CREATE INDEX idx_shifts_end_time ON public.shifts (end_time);
=======
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


CREATE TABLE Appointment (
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

<<<<<<< HEAD
-- -- Create indexes for better query performance
-- CREATE INDEX idx_device_tokens_user ON DeviceTokens(user_id, username);
-- CREATE INDEX idx_users_username ON users(username);
-- CREATE INDEX idx_users_email ON users(email);
>>>>>>> 0fb3f30 (user images)
=======
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
    medical_record_id int8 NULL,
	prescribing_vet varchar NULL,
	start_date date NULL,
	end_date date NULL,
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
    PRIMARY KEY (phase_id, medicine_id)
);
-- 4. Query để lấy lịch sử điều trị của một thú cưng
CREATE TABLE pet_treatments (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT REFERENCES Pet(petid),
    disease_id BIGINT REFERENCES diseases(id),
    start_date DATE,
    end_date DATE,
    status VARCHAR(50),  -- CHECK (status IN ('ongoing', 'completed', 'paused', 'cancelled')),
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

-- Index for pet_logs table
CREATE INDEX idx_pet_logs_pet_id ON pet_logs (petid);
CREATE INDEX idx_pet_logs_datetime ON pet_logs (datetime);


-- --------
INSERT INTO users (username, hashed_password, full_name, email, phone_number, address, original_image, role, is_verified_email)
VALUES
    ('hoangduong', '$2a$10$somehashedpassword1234567890', 'Hoàng Dương', 'hoang.duong@example.com', '0912345678', '789 Đường DEF, Hà Nội', 'hoangduong_avatar.png', 'user', true),
    ('linhtran', '$2a$10$anotherhashedpassword0987654321', 'Linh Trần', 'linh.tran@example.com', '0981234567', '321 Đường GHI, TP HCM', 'linhtran_avatar.png', 'user', true),
    ('nguyenhoa', '$2a$10$tLnDg/6/QNu/nD3bIcoR2OtUqNUci4jkzlswN6cHRxhJ4QuEOvXHW', 'Nguyễn Văn Hòa', 'hoa.nguyen@example.com', '0123456789', '123 Đường ABC, Hà Nội', 'hoa_avatar.png', 'admin', true);

-- Insert sample pets
INSERT INTO Pet (name, type, breed, age, gender, healthnotes, weight, birth_date, username, microchip_number, last_checkup_date, original_image)
VALUES
    ('Milo', 'Chó', 'Poodle', 3, 'Đực', 'Dị ứng nhẹ với phấn hoa', 5.2, '2020-04-15', 'nguyenhoa', 'MICRO123456', '2023-10-01', 'milo.png'),
    ('Luna', 'Mèo', 'Anh lông ngắn', 2, 'Cái', 'Không có vấn đề sức khỏe nghiêm trọng', 3.8, '2021-06-10', 'hoangduong', 'MICRO654321', '2023-09-20', 'luna.png'),
    ('Rocky', 'Chó', 'Golden Retriever', 5, 'Đực', 'Khỏe mạnh, không có vấn đề sức khỏe', 28.0, '2018-08-05', 'hoangduong', 'MICRO112233', '2023-11-15', 'rocky.png'),
    ('Bella', 'Mèo', 'Maine Coon', 4, 'Cái', 'Viêm da dị ứng nhẹ', 6.5, '2019-12-25', 'linhtran', 'MICRO445566', '2023-11-10', 'bella.png');



-- Insert sample appointments
INSERT INTO Appointment (petid, doctor_id, service_id, status, notes)
VALUES
    (1, 1, 1, 'Scheduled', 'Khám tổng quát cho Milo'),
    (2, 2, 2, 'Completed', 'Tiêm phòng dại cho Luna');

-- Insert sample vaccinations
INSERT INTO Vaccination (petID, vaccineName, dateAdministered, nextDueDate, vaccineProvider, batchNumber, notes)
VALUES
    (1, 'Tiêm phòng dại', '2023-08-10', '2024-08-10', 'VaccineCo', 'BATCH123', 'Không có phản ứng phụ'),
    (2, 'Tiêm phòng Parvo', '2023-05-20', '2024-05-20', 'VaccineCo', 'BATCH456', 'Thú cưng khỏe mạnh');

   -- Insert sample diseases
INSERT INTO diseases (name, description, symptoms)
VALUES 
    ('Viêm da dị ứng', 'Bệnh viêm da dị ứng thường xảy ra khi thú cưng tiếp xúc với các chất gây dị ứng.', '["Ngứa dữ dội", "Mẩn đỏ", "Rụng lông", "Da khô"]'),
    ('Nấm da', 'Bệnh nấm da là một bệnh phổ biến ở thú cưng, đặc biệt là chó và mèo', '["Ngứa nhiều", "Da đỏ", "Rụng lông từng mảng", "Vảy da"]'),
    ('Bệnh dại', 'Bệnh dại là một căn bệnh do virus gây ra, ảnh hưởng đến hệ thần kinh của thú cưng, có thể gây tử vong.', '["Thay đổi hành vi", "Sốt cao", "Bại liệt", "Khó thở", "Chảy dãi"]'),
    ('Bệnh Parvo (Chó con)', 'Bệnh Parvo là một bệnh viêm ruột cấp tính do virus gây ra, thường gặp ở chó con.', '["Nôn mửa", "Tiêu chảy có máu", "Mệt mỏi", "Không muốn ăn uống", "Sốt"]'),
    ('Tiêu chảy', 'Tiêu chảy thường do nhiễm khuẩn hoặc thức ăn không phù hợp.', '["Phân lỏng", "Mất nước", "Nôn mửa", "Chán ăn"]'),
    ('Nhiễm khuẩn đường hô hấp', 'Nhiễm khuẩn đường hô hấp là bệnh thường gặp ở thú cưng khi hệ miễn dịch yếu.', '["Ho", "Hắt hơi", "Sốt", "Khó thở"]'),
    ('Viêm khớp', 'Bệnh viêm khớp thường gặp ở thú cưng lớn tuổi, gây đau đớn và khó di chuyển.', '["Khó di chuyển", "Lè lưỡi", "Lưng cong", "Cứng khớp", "Kêu rên khi di chuyển"]'),
    ('Sán lá gan', 'Sán lá gan là một loại ký sinh trùng có thể ảnh hưởng đến gan của thú cưng.', '["Nôn mửa", "Tiêu chảy", "Suy giảm cân", "Vàng da"]'),
    ('Sốt rét', 'Sốt rét ở thú cưng thường do ký sinh trùng Plasmodium hoặc các loại kí sinh trùng khác gây ra.', '["Sốt cao", "Mệt mỏi", "Chán ăn", "Vàng mắt"]');

-- Insert sample medicines
INSERT INTO medicines (name, description, usage, dosage, frequency, duration, side_effects)
VALUES
    ('Ketoconazole', 'Thuốc kháng nấm dạng uống', 'Uống sau khi ăn', '5-10mg/kg thể trọng', '1 lần/ngày', '2-4 tuần', 'Có thể gây buồn nôn, chán ăn'),
    ('Miconazole', 'Thuốc kháng nấm dạng bôi', 'Bôi trực tiếp lên vùng da bị nấm', 'Bôi một lớp mỏng', '2 lần/ngày', '2-4 tuần', 'Có thể gây kích ứng da nhẹ');

-- Insert sample treatment phases
INSERT INTO treatment_phases (disease_id, phase_number, phase_name, description, duration, notes)
VALUES
    (1, 1, 'Giai đoạn cấp tính', 'Điều trị ban đầu để kiểm soát các triệu chứng', '1-2 tuần', 'Cần theo dõi sát trong giai đoạn này'),
    (2, 2, 'Giai đoạn duy trì', 'Tiếp tục điều trị để ngăn ngừa tái phát', '2-4 tuần', 'Có thể điều chỉnh liều dựa trên đáp ứng');


<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
INSERT INTO phase_medicines (phase_id, medicine_id, dosage, frequency, duration, notes)
VALUES 
(1, 1, '10mg/kg', '2 lần/ngày', '1 tuần', 'Uống sau bữa ăn'),
(1, 2, 'Bôi lớp mỏng', '3 lần/ngày', '1 tuần', 'Tránh để thú cưng liếm thuốc'),
(2, 1, '5mg/kg', '1 lần/ngày', '3 tuần', 'Uống sau bữa ăn'),
<<<<<<< HEAD
(2, 2, 'Bôi lớp mỏng', '2 lần/ngày', '3 tuần', 'Tiếp tục theo dõi phản ứng của da');

-- Query lấy phác đồ điều trị đầy đủ
SELECT 
    d.name AS disease_name,
    d.description AS disease_description,
    d.symptoms,
    tp.phase_number,
    tp.phase_name,
    tp.description AS phase_description,
    tp.duration AS phase_duration,
    tp.notes AS phase_notes,
    m.name AS medicine_name,
    m.description AS medicine_description,
    COALESCE(pm.dosage, m.dosage) AS dosage,
    COALESCE(pm.frequency, m.frequency) AS frequency,
    COALESCE(pm.duration, m.duration) AS duration,
    m.side_effects,
    pm.notes AS medicine_notes
FROM diseases d
JOIN treatment_phases tp ON d.id = tp.disease_id
JOIN phase_medicines pm ON tp.id = pm.phase_id
JOIN medicines m ON pm.medicine_id = m.id
WHERE LOWER(d.name) LIKE LOWER($1)
ORDER BY tp.phase_number, m.name;

-- 3. Query để lấy tổng quan điều trị
SELECT 
    d.name AS disease_name,
    d.description,
    d.symptoms,
    json_agg(
        json_build_object(
            'phase_number', tp.phase_number,
            'phase_name', tp.phase_name,
            'duration', tp.duration,
            'medicines', (
                SELECT json_agg(
                    json_build_object(
                        'name', m.name,
                        'dosage', COALESCE(pm.dosage, m.dosage),
                        'frequency', COALESCE(pm.frequency, m.frequency),
                        'duration', COALESCE(pm.duration, m.duration),
                        'notes', pm.notes
                    )
                )
                FROM phase_medicines pm
                JOIN medicines m ON pm.medicine_id = m.id
                WHERE pm.phase_id = tp.id
            )
        )
    ) AS treatment_phases
FROM diseases d
JOIN treatment_phases tp ON d.id = tp.disease_id
WHERE LOWER(d.name) LIKE LOWER($1)
GROUP BY d.id, d.name, d.description, d.symptoms;

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

-- Query lấy lịch sử điều trị
SELECT 
    p.name AS pet_name,
    d.name AS disease_name,
    pt.start_date,
    pt.end_date,
    pt.status,
    json_agg(
        json_build_object(
            'phase_name', tp.phase_name,
            'start_date', tpr.start_date,
            'end_date', tpr.end_date,
            'status', tpr.status,
            'notes', tpr.notes
        )
    ) AS progress
FROM pet_treatments pt
JOIN Pets p ON pt.pet_id = p.id
JOIN diseases d ON pt.disease_id = d.id
JOIN treatment_progress tpr ON pt.id = tpr.treatment_id
JOIN treatment_phases tp ON tpr.phase_id = tp.id
WHERE p.id = $1
GROUP BY p.id, p.name, d.name, pt.start_date, pt.end_date, pt.status
ORDER BY pt.start_date DESC;
>>>>>>> 6c35562 (dicease and treatment plan)
=======
(2, 2, 'Bôi lớp mỏng', '2 lần/ngày', '3 tuần', 'Tiếp tục theo dõi phản ứng của da');
>>>>>>> 50f041a (update database design)
=======
>>>>>>> a415f25 (new data)
=======
=======
>>>>>>> bd5945b (get list products)
INSERT INTO medicines (name, description, usage, dosage, frequency, duration, side_effects)
VALUES
    ('Rabies Vaccine', 'Vaccine để phòng bệnh dại', 'Tiêm bắp', '1 liều tiêu chuẩn', 'Một lần', '1 năm', 'Phản ứng dị ứng nhẹ'),
    ('Parvovirus Vaccine', 'Vaccine để phòng bệnh Parvo', 'Tiêm bắp', '1 liều tiêu chuẩn', 'Một lần', '3 năm', 'Không phổ biến'),
    ('Amoxicillin', 'Kháng sinh phổ rộng để điều trị nhiễm khuẩn đường hô hấp.', 'Uống sau ăn', '10mg/kg', '2 lần/ngày', '1 tuần', 'Có thể gây tiêu chảy hoặc chán ăn'),
    ('Meloxicam', 'Thuốc giảm đau và kháng viêm, thường dùng cho bệnh viêm khớp.', 'Uống sau ăn', '0.1mg/kg', '1 lần/ngày', '7-14 ngày', 'Có thể gây buồn nôn, loét dạ dày');
<<<<<<< HEAD
=======
-- INSERT INTO medicines (name, description, usage, dosage, frequency, duration, side_effects)
-- VALUES
--     ('Rabies Vaccine', 'Vaccine để phòng bệnh dại', 'Tiêm bắp', '1 liều tiêu chuẩn', 'Một lần', '1 năm', 'Phản ứng dị ứng nhẹ'),
--     ('Parvovirus Vaccine', 'Vaccine để phòng bệnh Parvo', 'Tiêm bắp', '1 liều tiêu chuẩn', 'Một lần', '3 năm', 'Không phổ biến'),
--     ('Amoxicillin', 'Kháng sinh phổ rộng để điều trị nhiễm khuẩn đường hô hấp.', 'Uống sau ăn', '10mg/kg', '2 lần/ngày', '1 tuần', 'Có thể gây tiêu chảy hoặc chán ăn'),
--     ('Meloxicam', 'Thuốc giảm đau và kháng viêm, thường dùng cho bệnh viêm khớp.', 'Uống sau ăn', '0.1mg/kg', '1 lần/ngày', '7-14 ngày', 'Có thể gây buồn nôn, loét dạ dày');
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> bd5945b (get list products)

  
   INSERT INTO treatment_phases (disease_id, phase_number, phase_name, description, duration, notes)
VALUES
    -- Giai đoạn cho Bệnh dại
    (5, 1, 'Giai đoạn phòng bệnh', 'Tiêm vaccine phòng bệnh dại định kỳ', 'Hằng năm', 'Theo dõi sau tiêm để phát hiện phản ứng dị ứng'),

    -- Giai đoạn cho Bệnh Parvo
    (6, 1, 'Điều trị cấp tính', 'Điều trị triệu chứng và bổ sung nước điện giải', '1 tuần', 'Cách ly để tránh lây nhiễm cho thú cưng khác'),

    -- Giai đoạn cho Nhiễm khuẩn đường hô hấp
    (7, 1, 'Giai đoạn cấp tính', 'Điều trị kháng sinh và hỗ trợ miễn dịch', '7 ngày', 'Sử dụng kháng sinh phổ rộng để giảm viêm nhiễm'),
    (7, 2, 'Giai đoạn phục hồi', 'Tiếp tục điều trị để tăng cường miễn dịch và phục hồi chức năng hô hấp', '1-2 tuần', 'Sử dụng bổ sung vitamin nếu cần'),

    -- Giai đoạn cho Viêm khớp
    (8, 1, 'Giai đoạn giảm đau', 'Sử dụng thuốc giảm đau và kháng viêm để cải thiện vận động', '1-2 tuần', 'Theo dõi tình trạng đau khi di chuyển'),
    (8, 2, 'Giai đoạn duy trì', 'Điều trị lâu dài để cải thiện chất lượng cuộc sống', 'Hằng năm', 'Duy trì tập luyện nhẹ nhàng và chế độ ăn phù hợp');

  INSERT INTO pet_treatments (pet_id, disease_id, start_date, end_date, status, notes)
VALUES
    (1, 1, '2024-10-01', '2024-10-15', 'Completed', 'Điều trị viêm da dị ứng thành công'),
    (2, 2, '2024-11-01', '2024-11-20', 'In Progress', 'Điều trị nấm da giai đoạn đầu'),
    (3, 3, '2024-09-01', NULL, 'Ongoing', 'Phòng ngừa bệnh dại với vaccine định kỳ');

   
INSERT INTO phase_medicines (phase_id, medicine_id, dosage, frequency, duration, notes)
VALUES
    -- Giai đoạn điều trị cho Viêm da dị ứng
    (1, 1, '5-10mg/kg thể trọng', '1 lần/ngày', '2-4 tuần', 'Dùng Ketoconazole để giảm triệu chứng nấm'),
    (1, 2, 'Bôi một lớp mỏng', '2 lần/ngày', '2-4 tuần', 'Dùng Miconazole tại chỗ để kiểm soát triệu chứng'),

    -- Giai đoạn điều trị cho Nấm da
    (2, 1, '5-10mg/kg thể trọng', '1 lần/ngày', '2-4 tuần', 'Dùng Ketoconazole để tiêu diệt nấm'),
    (2, 2, 'Bôi một lớp mỏng', '2 lần/ngày', '2-4 tuần', 'Dùng Miconazole để ngăn tái phát'),

    -- Giai đoạn điều trị cho Bệnh dại
    (3, 3, '1 liều tiêu chuẩn', 'Một lần', 'Hằng năm', 'Tiêm vaccine Rabies để phòng bệnh dại'),

    -- Giai đoạn điều trị cho Bệnh Parvo
    (4, 4, '10mg/kg', '2 lần/ngày', '1 tuần', 'Dùng Amoxicillin để kiểm soát nhiễm khuẩn đường ruột'),

<<<<<<< HEAD
<<<<<<< HEAD
    -- Giai đoạn điều trị cho Viêm khớp
    (5, 5, '0.1mg/kg', '1 lần/ngày', '7-14 ngày', 'Dùng Meloxicam để giảm đau và viêm khớp');
<<<<<<< HEAD
>>>>>>> 7af4c7a (new data)
=======
--     -- Giai đoạn điều trị cho Viêm khớp
--     (5, 5, '0.1mg/kg', '1 lần/ngày', '7-14 ngày', 'Dùng Meloxicam để giảm đau và viêm khớp');
>>>>>>> c449ffc (feat: cart api)
=======
    -- Giai đoạn điều trị cho Viêm khớp
    (5, 5, '0.1mg/kg', '1 lần/ngày', '7-14 ngày', 'Dùng Meloxicam để giảm đau và viêm khớp');
>>>>>>> bd5945b (get list products)
=======

>>>>>>> b0fe977 (place order and make payment)
