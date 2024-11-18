CREATE TABLE pet_logs (
    log_id BIGSERIAL PRIMARY KEY,
	petid int8 NOT NULL,
	datetime timestamp NULL,
	title varchar NULL,
	notes text NULL,
	CONSTRAINT newtable_pet_fk FOREIGN KEY (petid) REFERENCES pet(petid)
);
