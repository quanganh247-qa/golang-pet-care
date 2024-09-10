-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	username varchar NOT NULL,
	hashed_password varchar NOT NULL,
	full_name varchar NOT NULL,
	email varchar NOT NULL,
	password_changed_at timestamptz DEFAULT '0001-01-01 06:42:04+06:42:04'::timestamp with time zone NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	is_verified_email bool DEFAULT false NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (username)
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

CREATE TABLE public.verify_emails (
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
