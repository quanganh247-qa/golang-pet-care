-- Add new tables for comprehensive pet medical history

-- Examinations table to store medical examination details
CREATE TABLE public.examinations (
    id bigserial NOT NULL,
    medical_history_id int8 NOT NULL,
    exam_date timestamp NOT NULL,
    exam_type varchar(100) NOT NULL,
    findings text NOT NULL,
    vet_notes text NULL,
    doctor_id int8 NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT examinations_pkey PRIMARY KEY (id),
    CONSTRAINT examinations_medical_history_id_fkey FOREIGN KEY (medical_history_id) REFERENCES public.medical_history(id) ON DELETE CASCADE,
    CONSTRAINT examinations_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id)
);

-- Prescriptions table to store medication prescriptions
CREATE TABLE public.prescriptions (
    id bigserial NOT NULL,
    medical_history_id int8 NOT NULL,
    examination_id int8 NOT NULL,
    prescription_date timestamp NOT NULL,
    doctor_id int8 NOT NULL,
    notes text NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT prescriptions_pkey PRIMARY KEY (id),
    CONSTRAINT prescriptions_medical_history_id_fkey FOREIGN KEY (medical_history_id) REFERENCES public.medical_history(id) ON DELETE CASCADE,
    CONSTRAINT prescriptions_examination_id_fkey FOREIGN KEY (examination_id) REFERENCES public.examinations(id) ON DELETE CASCADE,
    CONSTRAINT prescriptions_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id)
);

-- Prescription medications junction table to track medications in prescriptions
CREATE TABLE public.prescription_medications (
    id bigserial NOT NULL,
    prescription_id int8 NOT NULL,
    medicine_id int8 NOT NULL,
    dosage text NOT NULL,
    frequency text NOT NULL,
    duration text NOT NULL,
    instructions text NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT prescription_medications_pkey PRIMARY KEY (id),
    CONSTRAINT prescription_medications_prescription_id_fkey FOREIGN KEY (prescription_id) REFERENCES public.prescriptions(id) ON DELETE CASCADE,
    CONSTRAINT prescription_medications_medicine_id_fkey FOREIGN KEY (medicine_id) REFERENCES public.medicines(id)
);

-- Test results table to store lab test results
CREATE TABLE public.test_results (
    id bigserial NOT NULL,
    medical_history_id int8 NOT NULL,
    examination_id int8 NOT NULL,
    test_type varchar(100) NOT NULL,
    test_date timestamp NOT NULL,
    results text NOT NULL,
    interpretation text NULL,
    file_url text NULL,
    doctor_id int8 NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT test_results_pkey PRIMARY KEY (id),
    CONSTRAINT test_results_medical_history_id_fkey FOREIGN KEY (medical_history_id) REFERENCES public.medical_history(id) ON DELETE CASCADE,
    CONSTRAINT test_results_examination_id_fkey FOREIGN KEY (examination_id) REFERENCES public.examinations(id) ON DELETE CASCADE,
    CONSTRAINT test_results_doctor_id_fkey FOREIGN KEY (doctor_id) REFERENCES public.doctors(id)
);

-- Create indexes for improved query performance
CREATE INDEX idx_examinations_medical_history_id ON public.examinations (medical_history_id);
CREATE INDEX idx_examinations_doctor_id ON public.examinations (doctor_id);
CREATE INDEX idx_examinations_exam_date ON public.examinations (exam_date);
CREATE INDEX idx_examinations_exam_type ON public.examinations (exam_type);

CREATE INDEX idx_prescriptions_medical_history_id ON public.prescriptions (medical_history_id);
CREATE INDEX idx_prescriptions_examination_id ON public.prescriptions (examination_id);
CREATE INDEX idx_prescriptions_doctor_id ON public.prescriptions (doctor_id);
CREATE INDEX idx_prescriptions_prescription_date ON public.prescriptions (prescription_date);

CREATE INDEX idx_prescription_medications_prescription_id ON public.prescription_medications (prescription_id);
CREATE INDEX idx_prescription_medications_medicine_id ON public.prescription_medications (medicine_id);

CREATE INDEX idx_test_results_medical_history_id ON public.test_results (medical_history_id);
CREATE INDEX idx_test_results_examination_id ON public.test_results (examination_id);
CREATE INDEX idx_test_results_doctor_id ON public.test_results (doctor_id);
CREATE INDEX idx_test_results_test_date ON public.test_results (test_date);
CREATE INDEX idx_test_results_test_type ON public.test_results (test_type);

-- Function to get complete medical history summary for a pet
CREATE OR REPLACE FUNCTION get_pet_medical_history_summary(p_pet_id bigint)
RETURNS json
LANGUAGE plpgsql
AS $$
DECLARE
    v_medical_record_id bigint;
    result json;
BEGIN
    -- Get the medical record ID for this pet
    SELECT id INTO v_medical_record_id 
    FROM medical_records 
    WHERE pet_id = p_pet_id;
    
    IF v_medical_record_id IS NULL THEN
        RETURN json_build_object(
            'success', false,
            'message', 'No medical record found for this pet'
        );
    END IF;
    
    -- Build the comprehensive medical history
    SELECT json_build_object(
        'medical_record', (
            SELECT json_build_object(
                'id', mr.id,
                'pet_id', mr.pet_id,
                'created_at', mr.created_at,
                'updated_at', mr.updated_at
            )
            FROM medical_records mr
            WHERE mr.id = v_medical_record_id
        ),
        'conditions', (
            SELECT json_agg(
                json_build_object(
                    'id', mh.id,
                    'condition', mh.condition,
                    'diagnosis_date', mh.diagnosis_date,
                    'notes', mh.notes,
                    'examinations', (
                        SELECT json_agg(
                            json_build_object(
                                'id', e.id,
                                'exam_date', e.exam_date,
                                'exam_type', e.exam_type,
                                'findings', e.findings,
                                'vet_notes', e.vet_notes,
                                'doctor_id', e.doctor_id,
                                'doctor_name', (SELECT u.full_name FROM doctors d JOIN users u ON d.user_id = u.id WHERE d.id = e.doctor_id),
                                'created_at', e.created_at,
                                'updated_at', e.updated_at
                            )
                        )
                        FROM examinations e
                        WHERE e.medical_history_id = mh.id
                        ORDER BY e.exam_date DESC
                    ),
                    'prescriptions', (
                        SELECT json_agg(
                            json_build_object(
                                'id', p.id,
                                'examination_id', p.examination_id,
                                'prescription_date', p.prescription_date,
                                'doctor_id', p.doctor_id,
                                'doctor_name', (SELECT u.full_name FROM doctors d JOIN users u ON d.user_id = u.id WHERE d.id = p.doctor_id),
                                'notes', p.notes,
                                'medications', (
                                    SELECT json_agg(
                                        json_build_object(
                                            'id', pm.id,
                                            'medicine_id', pm.medicine_id,
                                            'medicine_name', (SELECT m.name FROM medicines m WHERE m.id = pm.medicine_id),
                                            'dosage', pm.dosage,
                                            'frequency', pm.frequency,
                                            'duration', pm.duration,
                                            'instructions', pm.instructions
                                        )
                                    )
                                    FROM prescription_medications pm
                                    WHERE pm.prescription_id = p.id
                                )
                            )
                        )
                        FROM prescriptions p
                        WHERE p.medical_history_id = mh.id
                        ORDER BY p.prescription_date DESC
                    ),
                    'test_results', (
                        SELECT json_agg(
                            json_build_object(
                                'id', tr.id,
                                'examination_id', tr.examination_id,
                                'test_type', tr.test_type,
                                'test_date', tr.test_date,
                                'results', tr.results,
                                'interpretation', tr.interpretation,
                                'file_url', tr.file_url,
                                'doctor_id', tr.doctor_id,
                                'doctor_name', (SELECT u.full_name FROM doctors d JOIN users u ON d.user_id = u.id WHERE d.id = tr.doctor_id)
                            )
                        )
                        FROM test_results tr
                        WHERE tr.medical_history_id = mh.id
                        ORDER BY tr.test_date DESC
                    ),
                    'created_at', mh.created_at,
                    'updated_at', mh.updated_at
                )
            )
            FROM medical_history mh
            WHERE mh.medical_record_id = v_medical_record_id
            ORDER BY mh.diagnosis_date DESC
        ),
        'allergies', (
            SELECT json_agg(
                json_build_object(
                    'id', pa.id,
                    'type', pa.type,
                    'detail', pa.detail
                )
            )
            FROM pet_allergies pa
            WHERE pa.pet_id = p_pet_id
        )
    ) INTO result;
    
    RETURN json_build_object(
        'success', true,
        'data', result
    );
END;
$$;

