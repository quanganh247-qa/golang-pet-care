


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