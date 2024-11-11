CREATE TABLE Medications (
  medication_id BIGSERIAL NOT NULL,
  pet_id BIGINT NOT NULL,
  medication_name varchar(100) NOT NULL,
  dosage varchar(50) NOT NULL,
  frequency varchar(50) NOT NULL,
  start_date timestamp NOT NULL,
  end_date timestamp,
  notes text,
  PRIMARY KEY (medication_id),
  FOREIGN KEY (pet_id) REFERENCES Pet (petid) ON DELETE CASCADE
);