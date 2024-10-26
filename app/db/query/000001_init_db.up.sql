CREATE TABLE users (
  id BIGSERIAL NOT NULL,
  username varchar NOT NULL,
  hashed_password varchar NOT NULL,
  full_name varchar NOT NULL,
  email varchar NOT NULL,
  phone_number varchar,
  address varchar,
  avatar varchar(255),
  role VARCHAR(20) CHECK (role IN ('Admin', 'User')),
  created_at timestamptz NOT NULL DEFAULT (now()),
  is_verified_email bool DEFAULT false,
  removed_at timestamptz,
  PRIMARY KEY (id)
);

CREATE TABLE verify_emails (
  id BIGSERIAL NOT NULL,
  username varchar NOT NULL,
  email varchar NOT NULL,
  secret_code varchar NOT NULL,
  is_used bool NOT NULL DEFAULT false,
  created_at timestamptz DEFAULT (now()),
  expired_at timestamptz DEFAULT (now()+'00:15:00'::interval),
  PRIMARY KEY (id)
);

-- public.pet definition

-- Drop table

-- DROP TABLE public.pet;

CREATE TABLE pet (
	petid bigserial NOT NULL,
	"name" varchar(100) NOT NULL,
	"type" varchar(50) NOT NULL,
	breed varchar(100) NULL,
	age int4 NULL,
	gender varchar(10) NULL,
	healthnotes text NULL,
	profileimage varchar(255) NULL,
	weight float8 NULL,
	username varchar NOT NULL,
	CONSTRAINT pet_pkey PRIMARY KEY (petid),
	CONSTRAINT pet_users_fk FOREIGN KEY (username) REFERENCES users(username)
);

CREATE TABLE Vaccination (
  VaccinationID BIGSERIAL PRIMARY KEY,
  PetID BIGINT,
  VaccineName VARCHAR(100) NOT NULL,
  DateAdministered DATE NOT NULL,
  NextDueDate DATE,
  Notes TEXT
);

CREATE TABLE FeedingSchedule (
  FeedingScheduleID BIGSERIAL PRIMARY KEY,
  PetID BIGINT,
  MealTime TIME NOT NULL,
  FoodType VARCHAR(100) NOT NULL,
  Quantity DECIMAL(5,2) NOT NULL,
  Notes TEXT
);

CREATE TABLE ServiceType (
  TypeID BIGSERIAL PRIMARY KEY,
  ServiceTypeName varchar
);

CREATE TABLE Service (
  ServiceID BIGSERIAL PRIMARY KEY,
  TypeID BIGINT,
  Name varchar,
  Price DOUBLE PRECISION
);

CREATE TABLE Checkout (
  CheckoutID BIGSERIAL PRIMARY KEY,
  PetID BIGSERIAL,
  DoctorID BIGSERIAL,
  Date timestamp,
  Total_Amount DOUBLE PRECISION,
  Note varchar
);

CREATE TABLE Appointment (
  AppointmentID BIGSERIAL PRIMARY KEY,
  PetID BIGINT,
  DoctorID BIGINT,
  ServiceID BIGINT,
  Date timestamp,
  Status VARCHAR(20) CHECK (Status IN ('Pending', 'Completed', 'Cancelled'))
);

CREATE TABLE Checkout_Service (
  Checkout_Service_ID BIGSERIAL PRIMARY KEY,
  CheckoutID BIGSERIAL,
  ServiceID BIGSERIAL
);

CREATE UNIQUE INDEX users_username_key ON users (username);

ALTER TABLE verify_emails ADD CONSTRAINT verify_emails_users_fk FOREIGN KEY (username) REFERENCES users (username);

ALTER TABLE Pet ADD FOREIGN KEY (username) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE Vaccination ADD FOREIGN KEY (PetID) REFERENCES Pet (PetID) ON DELETE CASCADE;

ALTER TABLE FeedingSchedule ADD FOREIGN KEY (PetID) REFERENCES Pet (PetID) ON DELETE CASCADE;

ALTER TABLE Service ADD FOREIGN KEY (TypeID) REFERENCES ServiceType (TypeID) ON DELETE CASCADE;

ALTER TABLE Appointment ADD FOREIGN KEY (PetID) REFERENCES Pet (PetID) ON DELETE CASCADE;

ALTER TABLE Appointment ADD FOREIGN KEY (DoctorID) REFERENCES users (id) ON DELETE CASCADE;

ALTER TABLE Checkout ADD FOREIGN KEY (DoctorID) REFERENCES users (id) ON DELETE CASCADE;