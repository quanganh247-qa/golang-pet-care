-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE users (
	id bigserial NOT NULL,
	username varchar NOT NULL,
	hashed_password varchar NOT NULL,
	full_name varchar NOT NULL,
	email varchar NOT NULL,
	password_changed_at timestamptz DEFAULT '0001-01-01 06:42:04+06:42:04'::timestamp with time zone NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	is_verified_email bool DEFAULT false NULL,
	plan_type varchar NULL,
	removed_at timestamptz NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id),
	CONSTRAINT users_username_key UNIQUE (username)
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

-- public.categories definition

-- Drop table

-- DROP TABLE public.categories;

CREATE TABLE categories (
	id bigserial NOT NULL,
	category_name varchar NULL,
	CONSTRAINT categories_pk PRIMARY KEY (id),
	CONSTRAINT categories_unique UNIQUE (category_name)
);

-- Pages table (updated to store JSON)
CREATE TABLE pages (
	id bigserial NOT NULL,
	project_id int4 NULL,
	"name" varchar(100) NOT NULL,
	slug varchar(100) NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	"content" text NULL,
	category_name varchar NULL,
	component_code varchar NULL,
	removed_at timestamptz NULL,
	CONSTRAINT pages_pkey PRIMARY KEY (id),
	CONSTRAINT pages_project_id_slug_key UNIQUE (project_id, slug),
	CONSTRAINT pages_categories_fk FOREIGN KEY (category_name) REFERENCES categories(category_name),
	CONSTRAINT pages_components_fk FOREIGN KEY (component_code) REFERENCES components(component_code),
	CONSTRAINT pages_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- public.pages_historic definition

-- Drop table

-- DROP TABLE public.pages_historic;

CREATE TABLE pages_historic (
	id bigserial NOT NULL,
	page_id bigserial NOT NULL,
	page_version int4 NULL,
	created_at timestamptz NULL,
	removed_at timestamptz NULL,
	"content" text NULL,
	"name" varchar NULL,
	slug varchar NULL,
	category_name varchar NULL,
	CONSTRAINT pages_historic_pk PRIMARY KEY (id),
	CONSTRAINT pages_historic_pages_fk FOREIGN KEY (id) REFERENCES pages(id)
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
	id bigserial NOT NULL,
	project_id int4 NOT NULL,
	"name" varchar(100) NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
	component_code varchar NULL,
	description varchar NULL,
	"content" text NULL,
	removed_at timestamptz NULL,
	CONSTRAINT components_pkey PRIMARY KEY (id),
	CONSTRAINT components_unique UNIQUE (component_code),
	CONSTRAINT components_project_id_fkey FOREIGN KEY (project_id) REFERENCES projects(id)
);


-- public.subscriptions definition

CREATE TABLE subscriptions (
	id bigserial NOT NULL,
	username varchar NULL,
	plan_type varchar NULL,
	price numeric NULL,
	start_date timestamptz NULL,
	end_date timestamptz NULL,
	CONSTRAINT subscriptions_pk PRIMARY KEY (id),
	CONSTRAINT subscriptions_users_fk FOREIGN KEY (username) REFERENCES users(username)
);

-- public.domains definition

-- Drop table

-- DROP TABLE public.domains;

CREATE TABLE domains (
	id bigserial NOT NULL,
	username varchar NULL,
	page_id bigserial NOT NULL,
	domain_name varchar NULL,
	status varchar NULL,
	created_at timestamptz NULL,
	CONSTRAINT domains_pk PRIMARY KEY (id),
	CONSTRAINT domains_unique_1 UNIQUE (domain_name),
	CONSTRAINT domains_pages_fk FOREIGN KEY (id) REFERENCES pages(id),
	CONSTRAINT domains_users_fk FOREIGN KEY (username) REFERENCES users(username)
);

