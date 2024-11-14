ALTER TABLE public.pet_schedule ALTER COLUMN duration TYPE varchar USING duration::varchar;
