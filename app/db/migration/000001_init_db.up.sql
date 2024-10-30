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
  VaccinationID BIGSERIAL PRIMARY KEY,
  PetID BIGINT,
  VaccineName VARCHAR(100) NOT NULL,
  DateAdministered timestamp NOT NULL,
  NextDueDate timestamp,
  VaccineProvider VARCHAR(100),
  BatchNumber VARCHAR(50),
  Notes TEXT
);

CREATE TABLE FeedingSchedule (
  FeedingScheduleID BIGSERIAL PRIMARY KEY,
  PetID BIGINT,
  MealTime TIME NOT NULL,
  FoodType VARCHAR(100) NOT NULL,
  Quantity DECIMAL(5,2) NOT NULL,
  Frequency VARCHAR(50) NOT NULL,
  LastFed timestamp,
  Notes TEXT,
  IsActive BOOLEAN DEFAULT true
);

CREATE TABLE ActivityLog (
  LogID BIGSERIAL PRIMARY KEY,
  PetID BIGINT,
  ActivityType VARCHAR(50) NOT NULL,
  StartTime timestamp NOT NULL,
  Duration INTERVAL,
  Notes TEXT
);

CREATE TABLE Reminders (
  ReminderID BIGSERIAL PRIMARY KEY,
  PetID BIGINT,
  Title VARCHAR(100) NOT NULL,
  Description TEXT,
  DueDate timestamp NOT NULL,
  RepeatInterval VARCHAR(50),
  IsCompleted BOOLEAN DEFAULT false,
  NotificationSent BOOLEAN DEFAULT false
);

CREATE TABLE ServiceType (
  TypeID BIGSERIAL PRIMARY KEY,
  ServiceTypeName varchar NOT NULL,
  Description TEXT,
  IconURL VARCHAR(255)
);

CREATE TABLE Service (
  ServiceID BIGSERIAL PRIMARY KEY,
  TypeID BIGINT,
  Name varchar NOT NULL,
  Price float8,
  Duration INTERVAL,
  Description TEXT,
  IsAvailable BOOLEAN DEFAULT true
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
  CheckoutService_ID BIGSERIAL PRIMARY KEY,
  CheckoutID BIGINT,
  ServiceID BIGINT,
  Quantity INT DEFAULT 1,
  UnitPrice float8,
  Subtotal float8
);

CREATE TABLE PetServiceLocation (
  LocationID BIGSERIAL PRIMARY KEY,
  Name VARCHAR(100) NOT NULL,
  Address VARCHAR(255) NOT NULL,
  Latitude DECIMAL(10,8),
  Longitude DECIMAL(11,8),
  ContactNumber VARCHAR(20),
  BusinessHours VARCHAR(100),
  ServiceTypes TEXT[],
  Rating DECIMAL(3,2),
  IsVerified BOOLEAN DEFAULT false
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

ALTER TABLE Vaccination ADD CONSTRAINT vaccination_pet_fk FOREIGN KEY (PetID) REFERENCES Pet (petid);

ALTER TABLE FeedingSchedule ADD CONSTRAINT feeding_pet_fk FOREIGN KEY (PetID) REFERENCES Pet (petid);

ALTER TABLE ActivityLog ADD CONSTRAINT activity_pet_fk FOREIGN KEY (PetID) REFERENCES Pet (petid);

ALTER TABLE Reminders ADD CONSTRAINT reminder_pet_fk FOREIGN KEY (PetID) REFERENCES Pet (petid);

ALTER TABLE Service ADD CONSTRAINT service_type_fk FOREIGN KEY (TypeID) REFERENCES ServiceType (TypeID);

ALTER TABLE Appointment ADD CONSTRAINT appointment_pet_fk FOREIGN KEY (petid) REFERENCES Pet (petid);

ALTER TABLE Appointment ADD CONSTRAINT appointment_service_fk FOREIGN KEY (service_id) REFERENCES Service (ServiceID);

-- ALTER TABLE Checkout ADD CONSTRAINT checkout_pet_fk FOREIGN KEY (PetID) REFERENCES Pet (petid);

ALTER TABLE CheckoutService ADD CONSTRAINT cs_checkout_fk FOREIGN KEY (CheckoutID) REFERENCES Checkout (checkout_id);

ALTER TABLE CheckoutService ADD CONSTRAINT cs_service_fk FOREIGN KEY (ServiceID) REFERENCES Service (ServiceID);

ALTER TABLE Doctors ADD CONSTRAINT fk_doctor_user FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE DoctorSchedules ADD CONSTRAINT fk_schedule_doctor FOREIGN KEY (doctor_id) REFERENCES Doctors (id);

-- ALTER TABLE DoctorTimeOff ADD CONSTRAINT fk_timeoff_doctor FOREIGN KEY (doctor_id) REFERENCES Doctors (id);

-- ALTER TABLE Appointment ADD CONSTRAINT fk_appointment_timeslot FOREIGN KEY (time_slot_id) REFERENCES TimeSlots (id);
