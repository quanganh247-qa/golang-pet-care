-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE users (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	hashed_password varchar NOT NULL,
	full_name varchar NOT NULL,
	email varchar NOT NULL,
	phone_number varchar NULL,
	password_changed_at timestamptz DEFAULT '0001-01-01 06:42:04+06:42:04'::timestamp with time zone NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	is_verified_email bool DEFAULT false NULL,
	removed_at timestamptz NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
);

-- public.verify_emails definition

-- Drop table

-- DROP TABLE public.verify_emails;

CREATE TABLE verify_emails (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	email varchar NOT NULL,
	secret_code varchar NOT NULL,
	is_used bool DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT now() NULL,
	expired_at timestamptz DEFAULT now() + '00:15:00'::interval NULL,
	CONSTRAINT verify_emails_pk PRIMARY KEY (id),
	CONSTRAINT verify_emails_users_fk FOREIGN KEY (username) REFERENCES public.users(username)
);

CREATE TABLE Pet (
    PetID BIGSERIAL PRIMARY KEY,
    UserID BIGINT,
    Name VARCHAR(100) NOT NULL,
    Type VARCHAR(50) NOT NULL, -- Example: 'Chó', 'Mèo', 'Chim', etc.
    Breed VARCHAR(100),
    Age INT,
    Weight DECIMAL(5, 2),
    Gender VARCHAR(10),
    HealthNotes TEXT,
    ProfileImage VARCHAR(255),
    FOREIGN KEY (UserID) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE Vaccination (
    VaccinationID BIGSERIAL PRIMARY KEY,
    PetID BIGINT,
    VaccineName VARCHAR(100) NOT NULL,
    DateAdministered DATE NOT NULL,
    NextDueDate DATE,
    Notes TEXT,
    FOREIGN KEY (PetID) REFERENCES Pet(PetID) ON DELETE CASCADE
);

CREATE TABLE HealthCheck (
    HealthCheckID BIGSERIAL PRIMARY KEY,
    PetID BIGINT,
    Date DATE NOT NULL,
    ClinicName VARCHAR(100),
    DoctorName VARCHAR(100),
    Summary TEXT,
    NextCheckupDate DATE,
    FOREIGN KEY (PetID) REFERENCES Pet(PetID) ON DELETE CASCADE
);

CREATE TABLE FeedingSchedule (
    FeedingScheduleID BIGSERIAL PRIMARY KEY,
    PetID BIGINT,
    MealTime TIME NOT NULL,
    FoodType VARCHAR(100) NOT NULL,
    Quantity DECIMAL(5, 2) NOT NULL,
    Notes TEXT,
    FOREIGN KEY (PetID) REFERENCES Pet(PetID) ON DELETE CASCADE
);

CREATE TABLE Reminder (
    ReminderID BIGSERIAL PRIMARY KEY,
    UserID BIGINT,
    PetID BIGINT,
    ReminderType VARCHAR(50) NOT NULL, -- Example: 'Tiêm phòng', 'Kiểm tra sức khỏe', 'Cho ăn', etc.
    ReminderDate DATE NOT NULL,
    Description TEXT,
    Status VARCHAR(20) CHECK (Status IN ('Pending', 'Completed')),
    FOREIGN KEY (UserID) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (PetID) REFERENCES Pet(PetID) ON DELETE CASCADE
);

CREATE TABLE Service (
    ServiceID BIGSERIAL PRIMARY KEY,
    ServiceType VARCHAR(50) NOT NULL, -- Example: 'Phòng khám thú y', 'Cửa hàng thú cưng', 'Khách sạn thú cưng', etc.
    Name VARCHAR(100) NOT NULL,
    Address VARCHAR(255),
    PhoneNumber VARCHAR(20),
    Website VARCHAR(255),
    Rating DECIMAL(3, 2)
);

CREATE TABLE ServiceReview (
    ReviewID BIGSERIAL PRIMARY KEY,
    ServiceID BIGINT,
    UserID BIGINT,
    Rating DECIMAL(3, 2) NOT NULL CHECK (Rating BETWEEN 1.0 AND 5.0),
    Review TEXT,
    Date DATE NOT NULL,
    FOREIGN KEY (ServiceID) REFERENCES Service(ServiceID) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE PetCommunityPost (
    PostID BIGSERIAL PRIMARY KEY,
    UserID BIGINT,
    PetID BIGINT,
    Content TEXT NOT NULL,
    Image VARCHAR(255),
    PostDate DATE NOT NULL,
    FOREIGN KEY (UserID) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (PetID) REFERENCES Pet(PetID) ON DELETE CASCADE
);

CREATE TABLE PetCommunityComment (
    CommentID BIGSERIAL PRIMARY KEY,
    PostID BIGINT,
    UserID BIGINT,
    Content TEXT NOT NULL,
    CommentDate DATE NOT NULL,
    FOREIGN KEY (PostID) REFERENCES PetCommunityPost(PostID) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES users(id) ON DELETE CASCADE
);
