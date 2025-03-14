CREATE OR REPLACE FUNCTION create_time_slots_from_shift()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO public.time_slots (doctor_id, date, start_time, end_time, max_patients, booked_patients, created_at)
    SELECT
        NEW.doctor_id,
        DATE(NEW.start_time),
        (NEW.start_time + (n * INTERVAL '30 minutes'))::time,
        (NEW.start_time + ((n + 1) * INTERVAL '30 minutes'))::time,
        2,
        0,
        CURRENT_TIMESTAMP
    FROM generate_series(0, 
        EXTRACT(EPOCH FROM (NEW.end_time - NEW.start_time)) / 1800 - 1) AS n
    WHERE NEW.start_time < NEW.end_time;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER shift_insert_trigger
AFTER INSERT ON public.shifts
FOR EACH ROW EXECUTE FUNCTION create_time_slots_from_shift();