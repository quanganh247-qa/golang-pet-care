ALTER TABLE public.verify_emails ALTER COLUMN secret_code TYPE int8 USING secret_code::int8;
