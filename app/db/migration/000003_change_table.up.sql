
CREATE TABLE Doctors (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  specialization VARCHAR(100),
  years_of_experience INT,
  education TEXT,
  certificate_number VARCHAR(50),
  bio TEXT,
  consultation_fee FLOAT8
);


CREATE TABLE Departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE TimeSlots (
    id BIGSERIAL PRIMARY KEY,
    doctor_id INT NOT NULL,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    max_patients INT NULL,
    booked_patients INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (doctor_id) REFERENCES Doctors(id),
    CONSTRAINT unique_slot UNIQUE (doctor_id, date, start_time)
);

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


INSERT INTO medicines (name, description, dosage, frequency, duration, side_effects)
VALUES ('Paracetamol', 'For fever and pain relief', '1 tablet', '3 times a day', '5 days', 'None'),
       ('Amoxicillin', 'For bacterial infections', '1 tablet', '2 times a day', '7 days', 'Diarrhea, vomiting'),
       ('Ibuprofen', 'For pain relief', '1 tablet', '2 times a day', '3 days', 'Stomach pain, heartburn'),
       ('Cetirizine', 'For allergies', '1 tablet', '1 time a day', '5 days', 'Drowsiness, dry mouth'),
       ('Diazepam', 'For anxiety', '1 tablet', '1 time a day', '7 days', 'Drowsiness, dizziness'),
       ('Prednisone', 'For inflammation', '1 tablet', '1 time a day', '5 days', 'Weight gain, high blood pressure'),
       ('Cephalexin', 'For bacterial infections', '1 tablet', '2 times a day', '7 days', 'Diarrhea, vomiting'),
       ('Doxycycline', 'For bacterial infections', '1 tablet', '1 time a day', '7 days', 'Diarrhea, vomiting'),
       ('Clindamycin', 'For bacterial infections', '1 tablet', '3 times a day', '7 days', 'Diarrhea, vomiting'),
       ('Metronidazole', 'For bacterial infections', '1 tablet', '2 times a day', '7 days', 'Diarrhea, vomiting'),
       ('Fluconazole', 'For fungal infections', '1 tablet', '1 time a day', '7 days', 'Nausea, headache'),
       ('Ketoconazole', 'For fungal infections', '1 tablet', '1 time a day', '7 days', 'Nausea, headache'),
       ('Itraconazole', 'For fungal infections', '1 tablet', '1 time a day', '7 days', 'Nausea, headache'),
       ('Terbinafine', 'For fungal infections', '1 tablet', '1 time a day', '7 days', 'Nausea, headache'),
       ('Miconazole', 'For fungal infections', '1 tablet', '1 time a day', '7 days', 'Nausea, headache'),
       ('Acyclovir', 'For viral infections', '1 tablet', '5 times a day', '7 days', 'Nausea, headache'),
       ('Valacyclovir', 'For viral infections', '1 tablet', '2 times a day', '7 days', 'Nausea, headache'),
       ('Famciclovir', 'For viral infections', '1 tablet', '2 times a day', '7 days', 'Nausea, headache'),
       ('Oseltamivir', 'For viral infections', '1 tablet', '2 times a day', '5 days', 'Nausea, headache'),
       ('Zanamivir', 'For viral infections', '1 tablet', '2 times a day', '5 days', 'Nausea, headache');
       
       