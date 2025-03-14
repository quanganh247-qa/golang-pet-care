

-- Indexing for public.allergies
CREATE INDEX idx_allergies_medical_record_id ON public.allergies (medical_record_id);

-- Indexing for public.checkouts
CREATE INDEX idx_checkouts_petid ON public.checkouts (petid);
CREATE INDEX idx_checkouts_doctor_id ON public.checkouts (doctor_id);
CREATE INDEX idx_checkouts_date ON public.checkouts (date);

-- Indexing for public.departments
CREATE INDEX idx_departments_name ON public.departments (name);

-- Indexing for public.diseases
CREATE INDEX idx_diseases_name ON public.diseases (name);

-- Indexing for public.medical_history
CREATE INDEX idx_medical_history_medical_record_id ON public.medical_history (medical_record_id);
CREATE INDEX idx_medical_history_diagnosis_date ON public.medical_history (diagnosis_date);

-- Indexing for public.medical_records
CREATE INDEX idx_medical_records_pet_id ON public.medical_records (pet_id);

-- Indexing for public.medicines
CREATE INDEX idx_medicines_name ON public.medicines (name);
CREATE INDEX idx_medicines_expiration_date ON public.medicines (expiration_date);

-- Indexing for public.products
CREATE INDEX idx_products_name ON public.products (name);
CREATE INDEX idx_products_category ON public.products (category);
CREATE INDEX idx_products_is_available ON public.products (is_available);

-- Indexing for public.services
CREATE INDEX idx_services_name ON public.services (name);
CREATE INDEX idx_services_category ON public.services (category);

-- Indexing for public.states
CREATE INDEX idx_states_state ON public.states (state);

-- Indexing for public.users
CREATE INDEX idx_users_full_name ON public.users (full_name);
CREATE INDEX idx_users_role ON public.users (role);

-- Indexing for public.verify_emails
CREATE INDEX idx_verify_emails_username ON public.verify_emails (username);
CREATE INDEX idx_verify_emails_expired_at ON public.verify_emails (expired_at);

-- Indexing for public.carts
CREATE INDEX idx_carts_user_id ON public.carts (user_id);

-- Indexing for public.checkout_services
CREATE INDEX idx_checkout_services_checkoutid ON public.checkout_services (checkoutid);
CREATE INDEX idx_checkout_services_serviceid ON public.checkout_services (serviceid);

-- Indexing for public.device_tokens
CREATE INDEX idx_device_tokens_username ON public.device_tokens (username);

-- Indexing for public.doctors
CREATE INDEX idx_doctors_user_id ON public.doctors (user_id);
CREATE INDEX idx_doctors_specialization ON public.doctors (specialization);

-- Indexing for public.files
CREATE INDEX idx_files_user_id ON public.files (user_id);
CREATE INDEX idx_files_uploaded_at ON public.files (uploaded_at);

-- Indexing for public.notifications
CREATE INDEX idx_notifications_username ON public.notifications (username);
CREATE INDEX idx_notifications_datetime ON public.notifications (datetime);
CREATE INDEX idx_notifications_is_read ON public.notifications (is_read);

-- Indexing for public.orders
CREATE INDEX idx_orders_user_id ON public.orders (user_id);
CREATE INDEX idx_orders_order_date ON public.orders (order_date);
CREATE INDEX idx_orders_payment_status ON public.orders (payment_status);

-- Indexing for public.pets
CREATE INDEX idx_pets_username ON public.pets (username);
CREATE INDEX idx_pets_name ON public.pets (name);
CREATE INDEX idx_pets_is_active ON public.pets (is_active);

-- Indexing for public.time_slots
CREATE INDEX idx_time_slots_doctor_id ON public.time_slots (doctor_id);
CREATE INDEX idx_time_slots_date ON public.time_slots (date);
CREATE INDEX idx_time_slots_doctor_id_date ON public.time_slots (doctor_id, date);

-- Indexing for public.vaccinations
CREATE INDEX idx_vaccinations_petid ON public.vaccinations (petid);
CREATE INDEX idx_vaccinations_dateadministered ON public.vaccinations (dateadministered);
CREATE INDEX idx_vaccinations_nextduedate ON public.vaccinations (nextduedate);

-- Indexing for public.appointments
CREATE INDEX idx_appointments_petid ON public.appointments (petid);
CREATE INDEX idx_appointments_username ON public.appointments (username);
CREATE INDEX idx_appointments_doctor_id ON public.appointments (doctor_id);
CREATE INDEX idx_appointments_service_id ON public.appointments (service_id);
CREATE INDEX idx_appointments_date ON public.appointments (date);
CREATE INDEX idx_appointments_time_slot_id ON public.appointments (time_slot_id);
CREATE INDEX idx_appointments_state_id ON public.appointments (state_id);

-- Indexing for public.cart_items
CREATE INDEX idx_cart_items_cart_id ON public.cart_items (cart_id);
CREATE INDEX idx_cart_items_product_id ON public.cart_items (product_id);

-- Indexing for public.consultations
CREATE INDEX idx_consultations_appointment_id ON public.consultations (appointment_id);

-- Indexing for public.pet_logs
CREATE INDEX idx_pet_logs_petid ON public.pet_logs (petid);
CREATE INDEX idx_pet_logs_datetime ON public.pet_logs (datetime);

-- Indexing for public.pet_schedule
CREATE INDEX idx_pet_schedule_pet_id ON public.pet_schedule (pet_id);
CREATE INDEX idx_pet_schedule_reminder_datetime ON public.pet_schedule (reminder_datetime);

-- Indexing for public.pet_treatments
CREATE INDEX idx_pet_treatments_pet_id ON public.pet_treatments (pet_id);
CREATE INDEX idx_pet_treatments_disease_id ON public.pet_treatments (disease_id);
CREATE INDEX idx_pet_treatments_doctor_id ON public.pet_treatments (doctor_id);

-- Indexing for public.treatment_phases
CREATE INDEX idx_treatment_phases_treatment_id ON public.treatment_phases (treatment_id);

-- Indexing for public.phase_medicines
CREATE INDEX idx_phase_medicines_phase_id ON public.phase_medicines (phase_id);
CREATE INDEX idx_phase_medicines_medicine_id ON public.phase_medicines (medicine_id);

-- Indexing for public.clinics
CREATE INDEX idx_clinics_name ON public.clinics (name);

-- Indexing for public.shifts
CREATE INDEX idx_shifts_doctor_id ON public.shifts (doctor_id);
CREATE INDEX idx_shifts_start_time ON public.shifts (start_time);
CREATE INDEX idx_shifts_end_time ON public.shifts (end_time);