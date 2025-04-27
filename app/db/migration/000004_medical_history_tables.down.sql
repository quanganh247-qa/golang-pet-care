-- Drop function first
DROP FUNCTION IF EXISTS get_pet_medical_history_summary;

-- Drop indexes
DROP INDEX IF EXISTS idx_test_results_test_type;
DROP INDEX IF EXISTS idx_test_results_test_date;
DROP INDEX IF EXISTS idx_test_results_doctor_id;
DROP INDEX IF EXISTS idx_test_results_examination_id;
DROP INDEX IF EXISTS idx_test_results_medical_history_id;

DROP INDEX IF EXISTS idx_prescription_medications_medicine_id;
DROP INDEX IF EXISTS idx_prescription_medications_prescription_id;

DROP INDEX IF EXISTS idx_prescriptions_prescription_date;
DROP INDEX IF EXISTS idx_prescriptions_doctor_id;
DROP INDEX IF EXISTS idx_prescriptions_examination_id;
DROP INDEX IF EXISTS idx_prescriptions_medical_history_id;

DROP INDEX IF EXISTS idx_examinations_exam_type;
DROP INDEX IF EXISTS idx_examinations_exam_date;
DROP INDEX IF EXISTS idx_examinations_doctor_id;
DROP INDEX IF EXISTS idx_examinations_medical_history_id;

-- Drop tables in reverse order of creation to respect foreign key constraints
DROP TABLE IF EXISTS public.test_results;
DROP TABLE IF EXISTS public.prescription_medications;
DROP TABLE IF EXISTS public.prescriptions;
DROP TABLE IF EXISTS public.examinations;