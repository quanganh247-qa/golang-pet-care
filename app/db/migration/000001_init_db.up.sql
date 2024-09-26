-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE users (
    id bigserial  PRIMARY KEY,
	username varchar NOT NULL,
	hashed_password varchar NOT NULL,
	full_name varchar NOT NULL,
	email varchar NOT NULL,
	password_changed_at timestamptz DEFAULT '0001-01-01 06:42:04+06:42:04'::timestamp with time zone NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	is_verified_email bool DEFAULT false NULL
);

-- CREATE TABLE "sessions" (
--   "id" uuid PRIMARY KEY,
--   "username" varchar NOT NULL,
--   "refresh_token" varchar NOT NULL,
--   "user_agent" varchar NOT NULL,
--   "client_ip" varchar NOT NULL,
--   "is_blocked" boolean NOT NULL DEFAULT false,
--   "expires_at" timestamptz NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

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

-- Projects table
CREATE TABLE projects (
  	id bigserial NOT NULL,
	"name" varchar(100) NOT NULL,
	description text NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	username varchar NOT NULL,
	CONSTRAINT projects_pkey PRIMARY KEY (id),
	CONSTRAINT projects_users_fk FOREIGN KEY (username) REFERENCES public.users(username)
);

-- Pages table (updated to store JSON)
CREATE TABLE pages (
    id bigserial PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    content JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (project_id, slug)
);

-- Assets table
CREATE TABLE assets (
    id bigserial PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    name VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Components table (for reusable GrapesJS components)
CREATE TABLE components (
    id bigserial PRIMARY KEY,
    project_id INTEGER REFERENCES projects(id),
    name VARCHAR(100) NOT NULL,
    content JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
