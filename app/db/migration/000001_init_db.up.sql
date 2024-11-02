CREATE TABLE users (
  id BIGSERIAL NOT NULL,
  username varchar NOT NULL UNIQUE,
  hashed_password varchar NOT NULL,
  full_name varchar NOT NULL,
  email varchar NOT NULL,
  phone_number varchar,
  address varchar,
  avatar varchar(255),
  role VARCHAR(20),
  created_at timestamp NOT NULL DEFAULT (now()),
  is_verified_email bool DEFAULT false,
  removed_at timestamp,
  PRIMARY KEY (id)
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
  profileimage varchar(255),
  weight float8,
  birth_date date,
  username varchar NOT NULL,
  microchip_number varchar(50),
  last_checkup_date date,
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

CREATE TABLE FeedingSchedule (
  feedingScheduleID BIGSERIAL PRIMARY KEY,
  petID BIGINT,
  mealTime TIME NOT NULL,
  foodType VARCHAR(100) NOT NULL,
  quantity DECIMAL(5,2) NOT NULL,
  frequency VARCHAR(50) NOT NULL,
  lastFed timestamp,
  notes TEXT,
  isActive BOOLEAN DEFAULT true
);

CREATE TABLE ActivityLog (
  logID BIGSERIAL PRIMARY KEY,
  petID BIGINT,
  activityType VARCHAR(50) NOT NULL,
  startTime timestamp NOT NULL,
  duration INTERVAL,
  notes TEXT
);

CREATE TABLE Reminders (
  reminderID BIGSERIAL PRIMARY KEY,
  petID BIGINT,
  title VARCHAR(100) NOT NULL,
  description TEXT,
  dueDate timestamp NOT NULL,
  repeatInterval VARCHAR(50),
  isCompleted BOOLEAN DEFAULT false,
  notificationSent BOOLEAN DEFAULT false
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

CREATE TABLE PetServiceLocation (
  locationID BIGSERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  address VARCHAR(255) NOT NULL,
  latitude DECIMAL(10,8),
  longitude DECIMAL(11,8),
  contactNumber VARCHAR(20),
  businessHours VARCHAR(100),
  serviceTypes TEXT[],
  rating DECIMAL(3,2),
  isVerified BOOLEAN DEFAULT false
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

ALTER TABLE Pet ADD CONSTRAINT pet_users_fk FOREIGN KEY (username) REFERENCES users (username);

ALTER TABLE Vaccination ADD CONSTRAINT vaccination_pet_fk FOREIGN KEY (petID) REFERENCES Pet (petid);

ALTER TABLE FeedingSchedule ADD CONSTRAINT feeding_pet_fk FOREIGN KEY (petID) REFERENCES Pet (petid);

ALTER TABLE ActivityLog ADD CONSTRAINT activity_pet_fk FOREIGN KEY (petID) REFERENCES Pet (petid);

ALTER TABLE Reminders ADD CONSTRAINT reminder_pet_fk FOREIGN KEY (petID) REFERENCES Pet (petid);

ALTER TABLE Service ADD CONSTRAINT service_type_fk FOREIGN KEY (typeID) REFERENCES ServiceType (typeID);

ALTER TABLE Appointment ADD CONSTRAINT appointment_pet_fk FOREIGN KEY (petid) REFERENCES Pet (petid);

ALTER TABLE Appointment ADD CONSTRAINT appointment_service_fk FOREIGN KEY (service_id) REFERENCES Service (serviceID);

-- ALTER TABLE Checkout ADD CONSTRAINT checkout_pet_fk FOREIGN KEY (PetID) REFERENCES Pet (petid);

ALTER TABLE CheckoutService ADD CONSTRAINT cs_checkout_fk FOREIGN KEY (checkoutID) REFERENCES Checkout (checkout_id);

ALTER TABLE CheckoutService ADD CONSTRAINT cs_service_fk FOREIGN KEY (serviceID) REFERENCES Service (serviceID);

ALTER TABLE Doctors ADD CONSTRAINT fk_doctor_user FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE DoctorSchedules ADD CONSTRAINT fk_schedule_doctor FOREIGN KEY (doctor_id) REFERENCES Doctors (id);

-- ALTER TABLE DoctorTimeOff ADD CONSTRAINT fk_timeoff_doctor FOREIGN KEY (doctor_id) REFERENCES Doctors (id);

-- ALTER TABLE Appointment ADD CONSTRAINT fk_appointment_timeslot FOREIGN KEY (time_slot_id) REFERENCES TimeSlots (id);


-- public.token_info definition

-- Drop table

-- DROP TABLE public.token_info;

CREATE TABLE token_info (
	id bigserial NOT NULL,
	user_name varchar NOT NULL,
	access_token text NOT NULL,
	token_type varchar NOT NULL,
	refresh_token text NULL,
	expiry timestamptz NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	updated_at timestamptz DEFAULT now() NULL,
	CONSTRAINT token_info_pk PRIMARY KEY (id),
	CONSTRAINT token_info_unique UNIQUE (user_name)
);