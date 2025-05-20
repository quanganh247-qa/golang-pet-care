-- public.consultations definition

-- Drop table

-- DROP TABLE public.consultations;

CREATE TABLE public.consultations (
	id serial4 NOT NULL,
	appointment_id int8 NULL,
	subjective jsonb NULL,
	objective jsonb NULL,
	assessment jsonb NULL,
	plan int8 NULL,
	created_at timestamp DEFAULT now() NULL,
	CONSTRAINT consultations_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_consultations_appointment_id ON public.consultations USING btree (appointment_id);


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
CREATE INDEX idx_diseases_name ON public.diseases USING btree (name);


-- public.invoices definition

-- Drop table

-- DROP TABLE public.invoices;

CREATE TABLE public.invoices (
	id serial4 NOT NULL,
	invoice_number varchar(50) NOT NULL,
	amount float8 NOT NULL,
	"date" date DEFAULT CURRENT_DATE NOT NULL,
	due_date date NOT NULL,
	status varchar(20) DEFAULT 'unpaid'::character varying NOT NULL,
	description text NULL,
	customer_name varchar(100) NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	"type" varchar NULL,
	appointment_id int8 NULL,
	test_order_id int8 NULL,
	order_id int8 NULL,
	CONSTRAINT invoices_invoice_number_key UNIQUE (invoice_number),
	CONSTRAINT invoices_pkey PRIMARY KEY (id)
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
	created_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT medical_history_pk PRIMARY KEY (id)
);
CREATE INDEX idx_medical_history_diagnosis_date ON public.medical_history USING btree (diagnosis_date);
CREATE INDEX idx_medical_history_medical_record_id ON public.medical_history USING btree (medical_record_id);


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
CREATE INDEX idx_medical_records_pet_id ON public.medical_records USING btree (pet_id);


-- public.medicine_suppliers definition

-- Drop table

-- DROP TABLE public.medicine_suppliers;

CREATE TABLE public.medicine_suppliers (
	id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	email varchar(255) NULL,
	phone varchar(50) NULL,
	address text NULL,
	contact_name varchar(255) NULL,
	notes text NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT medicine_suppliers_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_medicine_suppliers_name ON public.medicine_suppliers USING btree (name);


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
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	expiration_date date NULL,
	quantity int8 NULL,
	unit_price float8 NULL,
	reorder_level int8 DEFAULT 10 NULL,
	supplier_id int8 NULL,
	CONSTRAINT medicines_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_medicines_expiration_date ON public.medicines USING btree (expiration_date);
CREATE INDEX idx_medicines_name ON public.medicines USING btree (name);
CREATE INDEX idx_medicines_reorder_level ON public.medicines USING btree (reorder_level);




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
	expiration_date date NULL,
	reorder_level int4 DEFAULT 10 NULL,
	supplier_id int8 NULL,
	updated_at timestamp NULL,
	CONSTRAINT products_pkey PRIMARY KEY (product_id)
);
CREATE INDEX idx_products_category ON public.products USING btree (category);
CREATE INDEX idx_products_is_available ON public.products USING btree (is_available);
CREATE INDEX idx_products_name ON public.products USING btree (name);


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
	created_at timestamp DEFAULT now() NULL,
	priority int2 DEFAULT 1 NULL,
	removed_at timestamp NULL,
	updated_at timestamp NULL,
	CONSTRAINT services_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_services_category ON public.services USING btree (category);
CREATE INDEX idx_services_name ON public.services USING btree (name);


-- public.smtp_configs definition

-- Drop table

-- DROP TABLE public.smtp_configs;

CREATE TABLE public.smtp_configs (
	id bigserial NOT NULL,
	"name" varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	"password" varchar(255) NOT NULL,
	smtp_host varchar(255) DEFAULT 'smtp.gmail.com'::character varying NOT NULL,
	smtp_port varchar(10) DEFAULT '587'::character varying NOT NULL,
	is_default bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT smtp_configs_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_smtp_configs_email ON public.smtp_configs USING btree (email);
CREATE INDEX idx_smtp_configs_is_default ON public.smtp_configs USING btree (is_default);


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
CREATE INDEX idx_states_state ON public.states USING btree (state);


-- public.test_categories definition

-- Drop table

-- DROP TABLE public.test_categories;

CREATE TABLE public.test_categories (
	id serial4 NOT NULL,
	category_id varchar(50) NOT NULL,
	"name" varchar(100) NOT NULL,
	description text NULL,
	icon_name varchar(50) NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT test_categories_category_id_key UNIQUE (category_id),
	CONSTRAINT test_categories_pkey PRIMARY KEY (id)
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
	created_at timestamp DEFAULT now() NOT NULL,
	is_verified_email bool DEFAULT false NULL,
	removed_at timestamp NULL,
	status varchar(20) NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
);
CREATE INDEX idx_users_full_name ON public.users USING btree (full_name);
CREATE INDEX idx_users_role ON public.users USING btree (role);


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
CREATE INDEX idx_verify_emails_expired_at ON public.verify_emails USING btree (expired_at);
CREATE INDEX idx_verify_emails_username ON public.verify_emails USING btree (username);


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
CREATE INDEX idx_carts_user_id ON public.carts USING btree (user_id);


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
	CONSTRAINT device_tokens_unique UNIQUE (token, username),
	CONSTRAINT fk_device_tokens_username FOREIGN KEY (username) REFERENCES public.users(username) ON DELETE CASCADE
);
CREATE INDEX idx_device_tokens_username ON public.device_tokens USING btree (username);


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
	consultation_fee float8 NULL,
	CONSTRAINT doctors_pkey PRIMARY KEY (id),
	CONSTRAINT fk_doctor_user FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE INDEX idx_doctors_specialization ON public.doctors USING btree (specialization);
CREATE INDEX idx_doctors_user_id ON public.doctors USING btree (user_id);


-- public.examinations definition

-- Drop table

-- DROP TABLE public.examinations;

CREATE TABLE public.examinations (
	id bigserial NOT NULL,
	medical_history_id int8 NOT NULL,
	exam_date timestamp NOT NULL,
	exam_type varchar(100) NOT NULL,
	findings text NOT NULL,
	vet_notes text NULL,
	doctor_id int8 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT examinations_pkey PRIMARY KEY (id),
	CONSTRAINT examinations_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id),
	CONSTRAINT examinations_medical_history_id_fkey FOREIGN KEY (medical_history_id) REFERENCES public.medical_history(id) ON DELETE CASCADE
);
CREATE INDEX idx_examinations_doctor_id ON public.examinations USING btree (doctor_id);
CREATE INDEX idx_examinations_exam_date ON public.examinations USING btree (exam_date);
CREATE INDEX idx_examinations_exam_type ON public.examinations USING btree (exam_type);
CREATE INDEX idx_examinations_medical_history_id ON public.examinations USING btree (medical_history_id);


-- public.invoice_items definition

-- Drop table

-- DROP TABLE public.invoice_items;

CREATE TABLE public.invoice_items (
	id serial4 NOT NULL,
	invoice_id int4 NOT NULL,
	"name" varchar(255) NOT NULL,
	price float8 NOT NULL,
	quantity int4 DEFAULT 1 NOT NULL,
	CONSTRAINT invoice_items_pkey PRIMARY KEY (id),
	CONSTRAINT invoice_items_invoice_id_fkey FOREIGN KEY (invoice_id) REFERENCES public.invoices(id) ON DELETE CASCADE
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
CREATE INDEX idx_notifications_datetime ON public.notifications USING btree (datetime);
CREATE INDEX idx_notifications_is_read ON public.notifications USING btree (is_read);
CREATE INDEX idx_notifications_username ON public.notifications USING btree (username);


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
CREATE INDEX idx_orders_order_date ON public.orders USING btree (order_date);
CREATE INDEX idx_orders_payment_status ON public.orders USING btree (payment_status);
CREATE INDEX idx_orders_user_id ON public.orders USING btree (user_id);


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
CREATE INDEX idx_pets_is_active ON public.pets USING btree (is_active);
CREATE INDEX idx_pets_name ON public.pets USING btree (name);
CREATE INDEX idx_pets_username ON public.pets USING btree (username);


-- public.prescriptions definition

-- Drop table

-- DROP TABLE public.prescriptions;

CREATE TABLE public.prescriptions (
	id bigserial NOT NULL,
	medical_history_id int8 NOT NULL,
	examination_id int8 NOT NULL,
	prescription_date timestamp NOT NULL,
	doctor_id int8 NOT NULL,
	notes text NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT prescriptions_pkey PRIMARY KEY (id),
	CONSTRAINT prescriptions_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id),
	CONSTRAINT prescriptions_examination_id_fkey FOREIGN KEY (examination_id) REFERENCES public.examinations(id) ON DELETE CASCADE,
	CONSTRAINT prescriptions_medical_history_id_fkey FOREIGN KEY (medical_history_id) REFERENCES public.medical_history(id) ON DELETE CASCADE
);
CREATE INDEX idx_prescriptions_doctor_id ON public.prescriptions USING btree (doctor_id);
CREATE INDEX idx_prescriptions_examination_id ON public.prescriptions USING btree (examination_id);
CREATE INDEX idx_prescriptions_medical_history_id ON public.prescriptions USING btree (medical_history_id);
CREATE INDEX idx_prescriptions_prescription_date ON public.prescriptions USING btree (prescription_date);


-- public.product_stock_movements definition

-- Drop table

-- DROP TABLE public.product_stock_movements;
-- +migrate Up
CREATE TYPE movement_type_enum AS ENUM ('import', 'export');

CREATE TABLE public.product_stock_movements (
	movement_id bigserial NOT NULL,
	product_id int8 NOT NULL,
	movement_type movement_type_enum NOT NULL,
	quantity int4 NOT NULL,
	reason text NULL,
	movement_date timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	price numeric(10, 2) DEFAULT 0.00 NOT NULL,
	CONSTRAINT product_stock_movements_pkey PRIMARY KEY (movement_id),
	CONSTRAINT product_stock_movements_quantity_check CHECK ((quantity > 0)),
	CONSTRAINT product_stock_movements_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(product_id) ON DELETE CASCADE
);
CREATE INDEX idx_product_stock_movements_movement_date ON public.product_stock_movements USING btree (movement_date);
CREATE INDEX idx_product_stock_movements_product_id ON public.product_stock_movements USING btree (product_id);


-- public.shifts definition

-- Drop table

-- DROP TABLE public.shifts;

CREATE TABLE public.shifts (
	id bigserial NOT NULL,
	doctor_id int8 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	max_patients int4 DEFAULT 10 NULL,
	"date" date NULL,
	CONSTRAINT shifts_pkey PRIMARY KEY (id),
	CONSTRAINT shifts_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id) ON DELETE CASCADE
);
CREATE INDEX idx_shifts_doctor_id ON public.shifts USING btree (doctor_id);


-- public.test_results definition

-- Drop table

-- DROP TABLE public.test_results;

CREATE TABLE public.test_results (
	id bigserial NOT NULL,
	medical_history_id int8 NOT NULL,
	examination_id int8 NOT NULL,
	test_type varchar(100) NOT NULL,
	test_date timestamp NOT NULL,
	results text NOT NULL,
	interpretation text NULL,
	file_url text NULL,
	doctor_id int8 NOT NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT test_results_pkey PRIMARY KEY (id),
	CONSTRAINT test_results_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id),
	CONSTRAINT test_results_examination_id_fkey FOREIGN KEY (examination_id) REFERENCES public.examinations(id) ON DELETE CASCADE,
	CONSTRAINT test_results_medical_history_id_fkey FOREIGN KEY (medical_history_id) REFERENCES public.medical_history(id) ON DELETE CASCADE
);
CREATE INDEX idx_test_results_doctor_id ON public.test_results USING btree (doctor_id);
CREATE INDEX idx_test_results_examination_id ON public.test_results USING btree (examination_id);
CREATE INDEX idx_test_results_medical_history_id ON public.test_results USING btree (medical_history_id);
CREATE INDEX idx_test_results_test_date ON public.test_results USING btree (test_date);
CREATE INDEX idx_test_results_test_type ON public.test_results USING btree (test_type);


-- public.tests definition

-- Drop table

-- DROP TABLE public.tests;

CREATE TABLE public.tests (
	id serial4 NOT NULL,
	test_id varchar(50) NOT NULL,
	category_id varchar(50) NULL,
	"name" varchar(100) NOT NULL,
	description text NULL,
	price float8 NOT NULL,
	turnaround_time varchar(50) NOT NULL,
	is_active bool DEFAULT true NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	quantity int4 DEFAULT 0 NULL,
	expiration_date date NULL,
	batch_number varchar(50) NULL,
	supplier_id int4 NULL,
	"type" varchar NULL,
	CONSTRAINT tests_pkey PRIMARY KEY (id),
	CONSTRAINT tests_test_id_key UNIQUE (test_id),
	CONSTRAINT tests_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.test_categories(category_id) ON DELETE CASCADE
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
	CONSTRAINT unique_slot UNIQUE (doctor_id, date, start_time),
	CONSTRAINT time_slots_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id) ON DELETE CASCADE,
	CONSTRAINT time_slots_shift_id_fkey FOREIGN KEY (shift_id) REFERENCES public.shifts(id) ON DELETE CASCADE
);
CREATE INDEX idx_time_slots_date ON public.time_slots USING btree (date);
CREATE INDEX idx_time_slots_doctor_id ON public.time_slots USING btree (doctor_id);
CREATE INDEX idx_time_slots_doctor_id_date ON public.time_slots USING btree (doctor_id, date);


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
CREATE INDEX idx_vaccinations_dateadministered ON public.vaccinations USING btree (dateadministered);
CREATE INDEX idx_vaccinations_nextduedate ON public.vaccinations USING btree (nextduedate);
CREATE INDEX idx_vaccinations_petid ON public.vaccinations USING btree (petid);


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
CREATE INDEX idx_cart_items_cart_id ON public.cart_items USING btree (cart_id);
CREATE INDEX idx_cart_items_product_id ON public.cart_items USING btree (product_id);


-- public.files definition

-- Drop table

-- DROP TABLE public.files;

CREATE TABLE public.files (
	id bigserial NOT NULL,
	file_name varchar(255) NOT NULL,
	file_path text NOT NULL,
	file_size int8 NOT NULL,
	file_type varchar(50) NOT NULL,
	uploaded_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	pet_id int8 NULL,
	CONSTRAINT files_pkey PRIMARY KEY (id),
	CONSTRAINT files_pet_id_fkey FOREIGN KEY (pet_id) REFERENCES public.pets(petid) ON DELETE SET NULL
);
CREATE INDEX idx_files_pet_id ON public.files USING btree (pet_id);
CREATE INDEX idx_files_uploaded_at ON public.files USING btree (uploaded_at);


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
CREATE INDEX idx_pet_logs_datetime ON public.pet_logs USING btree (datetime);
CREATE INDEX idx_pet_logs_petid ON public.pet_logs USING btree (petid);


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
CREATE INDEX idx_pet_schedule_pet_id ON public.pet_schedule USING btree (pet_id);
CREATE INDEX idx_pet_schedule_reminder_datetime ON public.pet_schedule USING btree (reminder_datetime);


-- public.pet_treatments definition

-- Drop table

-- DROP TABLE public.pet_treatments;

CREATE TABLE public.pet_treatments (
	id bigserial NOT NULL,
	pet_id int8 NULL,
	start_date date NULL,
	end_date date NULL,
	status varchar(50) NULL,
	description text NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	doctor_id int4 NULL,
	"name" varchar NULL,
	"type" varchar NULL,
	diseases varchar NULL,
	CONSTRAINT pet_treatments_pkey PRIMARY KEY (id),
	CONSTRAINT pet_treatments_pet_id_fkey FOREIGN KEY (pet_id) REFERENCES public.pets(petid)
);
CREATE INDEX idx_pet_treatments_doctor_id ON public.pet_treatments USING btree (doctor_id);
CREATE INDEX idx_pet_treatments_pet_id ON public.pet_treatments USING btree (pet_id);


-- public.pet_weight_history definition

-- Drop table

-- DROP TABLE public.pet_weight_history;

CREATE TABLE public.pet_weight_history (
	id bigserial NOT NULL,
	pet_id int8 NOT NULL,
	weight_kg float8 NOT NULL,
	recorded_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	notes text NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT pet_weight_history_pkey PRIMARY KEY (id),
	CONSTRAINT pet_weight_history_pet_id_fkey FOREIGN KEY (pet_id) REFERENCES public.pets(petid) ON DELETE CASCADE
);
CREATE INDEX idx_pet_weight_history_pet_id ON public.pet_weight_history USING btree (pet_id);
CREATE INDEX idx_pet_weight_history_recorded_at ON public.pet_weight_history USING btree (recorded_at);


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
	updated_at timestamptz NULL,
	is_locked bool DEFAULT false NULL,
	CONSTRAINT treatment_phases_pkey PRIMARY KEY (id),
	CONSTRAINT treatment_phases_treatment_id_fkey FOREIGN KEY (treatment_id) REFERENCES public.pet_treatments(id)
);
CREATE INDEX idx_treatment_phases_treatment_id ON public.treatment_phases USING btree (treatment_id);


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
CREATE INDEX idx_phase_medicines_medicine_id ON public.phase_medicines USING btree (medicine_id);
CREATE INDEX idx_phase_medicines_phase_id ON public.phase_medicines USING btree (phase_id);


-- public.appointments definition

-- Drop table

-- DROP TABLE public.appointments;

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
	updated_at timestamp NULL,
	CONSTRAINT appointments_pkey PRIMARY KEY (appointment_id)
);
CREATE INDEX idx_appointments_date ON public.appointments USING btree (date);
CREATE INDEX idx_appointments_doctor_id ON public.appointments USING btree (doctor_id);
CREATE INDEX idx_appointments_petid ON public.appointments USING btree (petid);
CREATE INDEX idx_appointments_service_id ON public.appointments USING btree (service_id);
CREATE INDEX idx_appointments_state_id ON public.appointments USING btree (state_id);
CREATE INDEX idx_appointments_time_slot_id ON public.appointments USING btree (time_slot_id);
CREATE INDEX idx_appointments_username ON public.appointments USING btree (username);


-- public.medicine_transactions definition

-- Drop table

-- DROP TABLE public.medicine_transactions;

CREATE TABLE public.medicine_transactions (
	id bigserial NOT NULL,
	medicine_id int8 NOT NULL,
	quantity int8 NOT NULL,
	transaction_type varchar(20) NOT NULL,
	unit_price float8 NULL,
	total_amount float8 NULL,
	transaction_date timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	supplier_id int8 NULL,
	expiration_date date NULL,
	notes text NULL,
	prescription_id int8 NULL,
	appointment_id int8 NULL,
	created_by varchar(255) NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT medicine_transactions_pkey PRIMARY KEY (id),
	CONSTRAINT medicine_transactions_transaction_type_check CHECK (((transaction_type)::text = ANY ((ARRAY['import'::character varying, 'export'::character varying])::text[])))
);
CREATE INDEX idx_medicine_transactions_medicine_id ON public.medicine_transactions USING btree (medicine_id);
CREATE INDEX idx_medicine_transactions_supplier_id ON public.medicine_transactions USING btree (supplier_id);
CREATE INDEX idx_medicine_transactions_transaction_date ON public.medicine_transactions USING btree (transaction_date);
CREATE INDEX idx_medicine_transactions_transaction_type ON public.medicine_transactions USING btree (transaction_type);


-- public.ordered_tests definition

-- Drop table

-- DROP TABLE public.ordered_tests;

CREATE TABLE public.ordered_tests (
	id serial4 NOT NULL,
	order_id int4 NULL,
	test_id int4 NULL,
	price_at_order float8 NOT NULL,
	status varchar(20) DEFAULT 'pending'::character varying NULL,
	results text NULL,
	results_date timestamptz NULL,
	technician_notes text NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	doctor_id int4 NULL,
	next_due_date date NULL,
	CONSTRAINT ordered_tests_pkey PRIMARY KEY (id)
);


-- public.payments definition

-- Drop table

-- DROP TABLE public.payments;

CREATE TABLE public.payments (
	id serial4 NOT NULL,
	amount float8 NOT NULL,
	payment_method varchar(50) NOT NULL,
	payment_status varchar(50) DEFAULT 'pending'::character varying NOT NULL,
	order_id int4 NULL,
	test_order_id int4 NULL,
	transaction_id varchar(255) NULL,
	payment_details jsonb NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	appointment_id int8 NULL,
	CONSTRAINT payments_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_payments_created_at ON public.payments USING btree (created_at);
CREATE INDEX idx_payments_order_id ON public.payments USING btree (order_id);
CREATE INDEX idx_payments_status ON public.payments USING btree (payment_status);
CREATE INDEX idx_payments_test_order_id ON public.payments USING btree (test_order_id);


-- public.rooms definition

-- Drop table

-- DROP TABLE public.rooms;

CREATE TABLE public.rooms (
	id bigserial NOT NULL,
	"name" varchar(100) NOT NULL,
	"type" varchar(50) NOT NULL,
	status varchar(20) DEFAULT 'available'::character varying NULL,
	current_appointment_id int8 NULL,
	available_at timestamp NULL,
	CONSTRAINT rooms_pkey PRIMARY KEY (id)
);


-- public.test_orders definition

-- Drop table

-- DROP TABLE public.test_orders;

CREATE TABLE public.test_orders (
	order_id serial4 NOT NULL,
	appointment_id int4 NULL,
	order_date timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	status varchar(20) DEFAULT 'pending'::character varying NULL,
	total_amount float8 NULL,
	notes text NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT test_orders_pkey PRIMARY KEY (order_id)
);


-- public.appointments foreign keys

ALTER TABLE public.appointments ADD CONSTRAINT appointment_pet_fk FOREIGN KEY (petid) REFERENCES public.pets(petid);
ALTER TABLE public.appointments ADD CONSTRAINT appointment_service_fk FOREIGN KEY (service_id) REFERENCES public.services(id);
ALTER TABLE public.appointments ADD CONSTRAINT appointments_room_id_fkey FOREIGN KEY (room_id) REFERENCES public.rooms(id);
ALTER TABLE public.appointments ADD CONSTRAINT appointments_state_id_fkey FOREIGN KEY (state_id) REFERENCES public.states(id);


-- public.medicine_transactions foreign keys

ALTER TABLE public.medicine_transactions ADD CONSTRAINT medicine_transactions_appointment_id_fkey FOREIGN KEY (appointment_id) REFERENCES public.appointments(appointment_id);
ALTER TABLE public.medicine_transactions ADD CONSTRAINT medicine_transactions_medicine_id_fkey FOREIGN KEY (medicine_id) REFERENCES public.medicines(id);
ALTER TABLE public.medicine_transactions ADD CONSTRAINT medicine_transactions_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.medicine_suppliers(id);


-- public.ordered_tests foreign keys

ALTER TABLE public.ordered_tests ADD CONSTRAINT ordered_tests_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.test_orders(order_id) ON DELETE CASCADE;
ALTER TABLE public.ordered_tests ADD CONSTRAINT ordered_tests_test_id_fkey FOREIGN KEY (test_id) REFERENCES public.tests(id) ON DELETE CASCADE;


-- public.payments foreign keys

ALTER TABLE public.payments ADD CONSTRAINT payments_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE SET NULL;
ALTER TABLE public.payments ADD CONSTRAINT payments_test_order_id_fkey FOREIGN KEY (test_order_id) REFERENCES public.test_orders(order_id) ON DELETE SET NULL;


-- public.rooms foreign keys

ALTER TABLE public.rooms ADD CONSTRAINT rooms_current_appointment_id_fkey FOREIGN KEY (current_appointment_id) REFERENCES public.appointments(appointment_id);


-- public.test_orders foreign keys

ALTER TABLE public.test_orders ADD CONSTRAINT test_orders_appointment_id_fkey FOREIGN KEY (appointment_id) REFERENCES public.appointments(appointment_id) ON DELETE CASCADE;