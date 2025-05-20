package medical_records

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicalRecordServiceInterface interface {
	CreateMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error)
	CreateMedicalHistory(ctx *gin.Context, req *MedicalHistoryRequest, recordID int64) (*MedicalHistoryResponse, error)
	GetMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error)
	ListMedicalHistory(ctx *gin.Context, recordID int64, pagination *util.Pagination) ([]MedicalHistoryResponse, error)
	GetMedicalHistoryByID(ctx *gin.Context, medicalHistoryID int64) (*MedicalHistoryResponse, error)
	GetMedicalHistoryByPetID(ctx *gin.Context, petID int64) ([]MedicalHistoryResponse, error)

	// // Updated methods for examinations
	// CreateExamination(ctx *gin.Context, req *ExaminationRequest) (*ExaminationResponse, error)
	// GetExamination(ctx *gin.Context, examinationID int64) (*ExaminationResponse, error)
	// ListExaminationsByMedicalHistory(ctx *gin.Context, medicalHistoryID int64, pagination *util.Pagination) ([]ExaminationResponse, error)
	// ListExaminationsByPet(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]ExaminationResponse, error)

	// // Updated methods for prescriptions
	// CreatePrescription(ctx *gin.Context, req *PrescriptionRequest) (*PrescriptionResponse, error)
	// GetPrescription(ctx *gin.Context, prescriptionID int64) (*PrescriptionResponse, error)
	// ListPrescriptionsByMedicalHistory(ctx *gin.Context, medicalHistoryID int64, pagination *util.Pagination) ([]PrescriptionResponse, error)
	// ListPrescriptionsByPet(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PrescriptionResponse, error)

	// // Updated methods for test results
	// CreateTestResult(ctx *gin.Context, req *TestResultRequest) (*TestResultResponse, error)
	// GetTestResult(ctx *gin.Context, testResultID int64) (*TestResultResponse, error)
	// ListTestResultsByMedicalHistory(ctx *gin.Context, medicalHistoryID int64, pagination *util.Pagination) ([]TestResultResponse, error)
	// ListTestResultsByPet(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]TestResultResponse, error)

	// Comprehensive medical history
	// GetCompleteMedicalHistory(ctx *gin.Context, petID int64) (*MedicalHistorySummary, error)
	GetAllSoapNotesByPetID(ctx *gin.Context, petID int64) ([]SoapNoteResponse, error)

	// New method for appointment visit summary
	GetAppointmentVisitSummary(ctx *gin.Context, appointmentID int64) (*AppointmentVisitSummaryResponse, error)
}

// Quản lý Bệnh án
func (s *MedicalRecordService) CreateMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error) {

	var rec db.MedicalRecord
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		rec, err = q.CreateMedicalRecord(ctx, pgtype.Int8{Int64: petID, Valid: true})
		if err != nil {
			log.Printf("failed to create medical record: %v", err)
			return fmt.Errorf("failed to create medical record: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("failed to create medical record: %v", err)
		return nil, fmt.Errorf("failed to create medical record: %w", err)
	}

	return &MedicalRecordResponse{
		ID:        rec.ID,
		PetID:     rec.PetID.Int64,
		CreatedAt: rec.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: rec.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

// medical history
func (s *MedicalRecordService) CreateMedicalHistory(ctx *gin.Context, req *MedicalHistoryRequest, recordID int64) (*MedicalHistoryResponse, error) {
	var rec db.MedicalHistory
	var err error

	layout := "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, req.DiagnosisDate)

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {

		rec, err = q.CreateMedicalHistory(ctx, db.CreateMedicalHistoryParams{
			MedicalRecordID: pgtype.Int8{Int64: recordID, Valid: true},
			Condition:       pgtype.Text{String: req.Condition, Valid: true},
			DiagnosisDate:   pgtype.Timestamp{Time: t, Valid: true},
			Notes:           pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			log.Printf("failed to create medical history: %v", err)
			return fmt.Errorf("failed to create medical history: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Printf("failed to create medical history: %v", err)
		return nil, fmt.Errorf("failed to create medical history: %w", err)
	}

	return &MedicalHistoryResponse{
		ID:              rec.ID,
		MedicalRecordID: rec.MedicalRecordID.Int64,
		Condition:       rec.Condition.String,
		DiagnosisDate:   rec.DiagnosisDate.Time.Format("2006-01-02 15:04:05"),
		Notes:           rec.Notes.String,
		CreatedAt:       rec.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:       rec.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *MedicalRecordService) GetMedicalRecord(ctx *gin.Context, petID int64) (*MedicalRecordResponse, error) {
	// First verify pet exists and is active
	pet, err := s.storeDB.GetPetByID(ctx, petID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify pet: %w", err)
	}
	if !pet.IsActive.Bool {
		return nil, fmt.Errorf("pet is not active")
	}

	record, err := s.storeDB.GetMedicalRecord(ctx, pgtype.Int8{Int64: petID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get medical record: %w", err)
	}

	return &MedicalRecordResponse{
		ID:        record.ID,
		PetID:     record.PetID.Int64,
		CreatedAt: record.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt: record.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *MedicalRecordService) ListMedicalHistory(ctx *gin.Context, recordID int64, pagination *util.Pagination) ([]MedicalHistoryResponse, error) {
	offset := (pagination.Page - 1) * pagination.PageSize
	medicalHistories, err := s.storeDB.GetMedicalHistory(ctx, db.GetMedicalHistoryParams{
		MedicalRecordID: pgtype.Int8{Int64: recordID, Valid: true},
		Limit:           int32(pagination.PageSize),
		Offset:          int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get medical history: %w", err)
	}

	var medicalHistoryResponse []MedicalHistoryResponse
	for _, medicalHistory := range medicalHistories {
		medicalHistoryResponse = append(medicalHistoryResponse, MedicalHistoryResponse{
			ID:              medicalHistory.ID,
			MedicalRecordID: medicalHistory.MedicalRecordID.Int64,
		})
	}

	return medicalHistoryResponse, nil
}

func (s *MedicalRecordService) GetMedicalHistoryByID(ctx *gin.Context, medicalHistoryID int64) (*MedicalHistoryResponse, error) {
	medicalHistory, err := s.storeDB.GetMedicalHistoryByID(ctx, medicalHistoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical history: %w", err)
	}

	return &MedicalHistoryResponse{
		ID:              medicalHistory.ID,
		MedicalRecordID: medicalHistory.MedicalRecordID.Int64,
	}, nil
}

func (s *MedicalRecordService) GetMedicalHistoryByPetID(ctx *gin.Context, petID int64) ([]MedicalHistoryResponse, error) {
	medicalHistories, err := s.storeDB.GetMedicalHistoryByPetID(ctx, pgtype.Int8{Int64: petID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get medical history: %w", err)
	}

	var medicalHistoryResponse []MedicalHistoryResponse
	for _, medicalHistory := range medicalHistories {
		medicalHistoryResponse = append(medicalHistoryResponse, MedicalHistoryResponse{
			ID:              medicalHistory.ID,
			MedicalRecordID: medicalHistory.MedicalRecordID.Int64,
			Condition:       medicalHistory.Condition.String,
			DiagnosisDate:   medicalHistory.DiagnosisDate.Time.Format("2006-01-02 15:04:05"),
			Notes:           medicalHistory.Notes.String,
			CreatedAt:       medicalHistory.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			UpdatedAt:       medicalHistory.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return medicalHistoryResponse, nil
}

// Examination-related service methods

// // CreateExamination creates a new examination record for a pet's medical history
// func (s *MedicalRecordService) CreateExamination(ctx *gin.Context, req *ExaminationRequest) (*ExaminationResponse, error) {
// 	var examination ExaminationResponse

// 	// Parse the examination date from string to time.Time
// 	layout := "2006-01-02 15:04:05"
// 	examDate, err := time.Parse(layout, req.ExamDate)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid examination date format: %w", err)
// 	}

// 	// Begin transaction
// 	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		// Check if medical history entry exists
// 		medicalHistory, err := q.GetMedicalHistoryById(ctx, req.MedicalHistoryID)
// 		if err != nil {
// 			return fmt.Errorf("failed to get medical history: %w", err)
// 		}

// 		// Check if doctor exists
// 		doctor, err := q.GetDoctor(ctx, req.DoctorID)
// 		if err != nil {
// 			return fmt.Errorf("failed to get doctor: %w", err)
// 		}

// 		// Create examination record
// 		result, err := q.CreateExamination(ctx, db.CreateExaminationParams{
// 			MedicalHistoryID: req.MedicalHistoryID,
// 			ExamDate:         pgtype.Timestamp{Time: examDate, Valid: true},
// 			ExamType:         req.ExamType,
// 			Findings:         req.Findings,
// 			VetNotes:         pgtype.Text{String: req.VetNotes, Valid: true},
// 			DoctorID:         req.DoctorID,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to create examination: %w", err)
// 		}

// 		// Populate response
// 		examination = ExaminationResponse{
// 			ID:               result.ID,
// 			MedicalHistoryID: medicalHistory.ID,
// 			ExamDate:         result.ExamDate.Time.Format(layout),
// 			ExamType:         result.ExamType,
// 			Findings:         result.Findings,
// 			VetNotes:         result.VetNotes.String,
// 			DoctorID:         doctor.ID,
// 			DoctorName:       doctor.Name,
// 			CreatedAt:        result.CreatedAt.Time.Format(layout),
// 			UpdatedAt:        result.UpdatedAt.Time.Format(layout),
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &examination, nil
// }

// GetExamination retrieves an examination by ID
// func (s *MedicalRecordService) GetExamination(ctx *gin.Context, examinationID int64) (*ExaminationResponse, error) {
// 	examination, err := s.storeDB.GetExaminationByID(ctx, examinationID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get examination: %w", err)
// 	}

// 	// Get doctor info
// 	doctor, err := s.storeDB.GetDoctor(ctx, examination.DoctorID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get doctor: %w", err)
// 	}

// 	layout := "2006-01-02 15:04:05"
// 	return &ExaminationResponse{
// 		ID:               examination.ID,
// 		MedicalHistoryID: examination.MedicalHistoryID,
// 		ExamDate:         examination.ExamDate.Time.Format(layout),
// 		ExamType:         examination.ExamType,
// 		Findings:         examination.Findings,
// 		VetNotes:         examination.VetNotes.String,
// 		DoctorID:         examination.DoctorID,
// 		DoctorName:       doctor.Name,
// 		CreatedAt:        examination.CreatedAt.Time.Format(layout),
// 		UpdatedAt:        examination.UpdatedAt.Time.Format(layout),
// 	}, nil
// }

// ListExaminationsByMedicalHistory lists all examinations for a medical history with pagination
// func (s *MedicalRecordService) ListExaminationsByMedicalHistory(ctx *gin.Context, medicalHistoryID int64, pagination *util.Pagination) ([]ExaminationResponse, error) {
// 	offset := (pagination.Page - 1) * pagination.PageSize

// 	examinations, err := s.storeDB.ListExaminationsByMedicalHistory(ctx, db.ListExaminationsByMedicalHistoryParams{
// 		MedicalHistoryID: medicalHistoryID,
// 		Limit:            int32(pagination.PageSize),
// 		Offset:           int32(offset),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list examinations: %w", err)
// 	}

// 	var results []ExaminationResponse
// 	layout := "2006-01-02 15:04:05"

// 	for _, examination := range examinations {
// 		// Get doctor info
// 		doctor, err := s.storeDB.GetDoctor(ctx, examination.DoctorID)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get doctor: %w", err)
// 		}

// 		results = append(results, ExaminationResponse{
// 			ID:               examination.ID,
// 			MedicalHistoryID: examination.MedicalHistoryID,
// 			ExamDate:         examination.ExamDate.Time.Format(layout),
// 			ExamType:         examination.ExamType,
// 			Findings:         examination.Findings,
// 			VetNotes:         examination.VetNotes.String,
// 			DoctorID:         examination.DoctorID,
// 			DoctorName:       doctor.Name,
// 			CreatedAt:        examination.CreatedAt.Time.Format(layout),
// 			UpdatedAt:        examination.UpdatedAt.Time.Format(layout),
// 		})
// 	}

// 	return results, nil
// }

// // ListExaminationsByPet lists all examinations for a pet with pagination
// func (s *MedicalRecordService) ListExaminationsByPet(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]ExaminationResponse, error) {
// 	// Get medical record for pet
// 	medicalRecord, err := s.storeDB.GetMedicalRecord(ctx, pgtype.Int8{Int64: petID, Valid: true})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get medical record: %w", err)
// 	}

// 	// Use the ListExaminationsByMedicalHistory function
// 	return s.ListExaminationsByMedicalHistory(ctx, medicalRecord.ID, pagination)
// }

// CreatePrescription creates a new prescription record
// func (s *MedicalRecordService) CreatePrescription(ctx *gin.Context, req *PrescriptionRequest) (*PrescriptionResponse, error) {
// 	var prescriptionResponse PrescriptionResponse

// 	// Parse the prescription date
// 	layout := "2006-01-02 15:04:05"
// 	prescriptionDate, err := time.Parse(layout, req.PrescriptionDate)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid prescription date format: %w", err)
// 	}

// 	// Begin transaction
// 	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		// Verify that medical history exists
// 		_, err := q.GetMedicalHistoryById(ctx, req.MedicalHistoryID)
// 		if err != nil {
// 			return fmt.Errorf("failed to get medical history: %w", err)
// 		}

// 		// Verify that examination exists
// 		examination, err := q.GetExaminationByID(ctx, req.ExaminationID)
// 		if err != nil {
// 			return fmt.Errorf("failed to get examination: %w", err)
// 		}

// 		// Verify that examination belongs to this medical history
// 		if examination.MedicalHistoryID != req.MedicalHistoryID {
// 			return fmt.Errorf("examination does not belong to the specified medical history")
// 		}

// 		// Verify that doctor exists
// 		doctor, err := q.GetDoctor(ctx, req.DoctorID)
// 		if err != nil {
// 			return fmt.Errorf("failed to get doctor: %w", err)
// 		}

// 		// Create prescription
// 		prescription, err := q.CreatePrescription(ctx, db.CreatePrescriptionParams{
// 			MedicalHistoryID: req.MedicalHistoryID,
// 			ExaminationID:    req.ExaminationID,
// 			PrescriptionDate: pgtype.Timestamp{Time: prescriptionDate, Valid: true},
// 			DoctorID:         req.DoctorID,
// 			Notes:            pgtype.Text{String: req.Notes, Valid: true},
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to create prescription: %w", err)
// 		}

// 		// Add medications to the prescription
// 		var medicationResponses []PrescriptionMedicationResponse
// 		for _, medication := range req.Medications {
// 			// Verify medicine exists
// 			medicine, err := q.GetMedicineByID(ctx, medication.MedicineID)
// 			if err != nil {
// 				return fmt.Errorf("failed to get medicine (ID: %d): %w", medication.MedicineID, err)
// 			}

// 			// Create prescription medication entry
// 			prescriptionMed, err := q.CreatePrescriptionMedication(ctx, db.CreatePrescriptionMedicationParams{
// 				PrescriptionID: prescription.ID,
// 				MedicineID:     medication.MedicineID,
// 				Dosage:         medication.Dosage,
// 				Frequency:      medication.Frequency,
// 				Duration:       medication.Duration,
// 				Instructions:   pgtype.Text{String: medication.Instructions, Valid: true},
// 			})
// 			if err != nil {
// 				return fmt.Errorf("failed to add medicine to prescription: %w", err)
// 			}

// 			medicationResponses = append(medicationResponses, PrescriptionMedicationResponse{
// 				ID:             prescriptionMed.ID,
// 				PrescriptionID: prescription.ID,
// 				MedicineID:     medicine.ID,
// 				MedicineName:   medicine.Name,
// 				Dosage:         prescriptionMed.Dosage,
// 				Frequency:      prescriptionMed.Frequency,
// 				Duration:       prescriptionMed.Duration,
// 				Instructions:   prescriptionMed.Instructions.String,
// 			})
// 		}

// 		// Populate response
// 		prescriptionResponse = PrescriptionResponse{
// 			ID:               prescription.ID,
// 			MedicalHistoryID: prescription.MedicalHistoryID,
// 			ExaminationID:    prescription.ExaminationID,
// 			PrescriptionDate: prescription.PrescriptionDate.Time.Format(layout),
// 			DoctorID:         prescription.DoctorID,
// 			DoctorName:       doctor.Name,
// 			Notes:            prescription.Notes.String,
// 			Medications:      medicationResponses,
// 			CreatedAt:        prescription.CreatedAt.Time.Format(layout),
// 			UpdatedAt:        prescription.UpdatedAt.Time.Format(layout),
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &prescriptionResponse, nil
// }

// GetPrescription retrieves a prescription by ID with all medications
// func (s *MedicalRecordService) GetPrescription(ctx *gin.Context, prescriptionID int64) (*PrescriptionResponse, error) {
// 	prescription, err := s.storeDB.GetPrescriptionByID(ctx, prescriptionID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get prescription: %w", err)
// 	}

// 	// Get doctor info
// 	doctor, err := s.storeDB.GetDoctor(ctx, prescription.DoctorID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get doctor: %w", err)
// 	}

// 	// Get medications for this prescription
// 	medications, err := s.storeDB.ListPrescriptionMedications(ctx, prescriptionID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get prescription medications: %w", err)
// 	}

// 	var medicationResponses []PrescriptionMedicationResponse
// 	for _, med := range medications {
// 		// Get medicine details
// 		medicine, err := s.storeDB.GetMedicineByID(ctx, med.MedicineID)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get medicine details: %w", err)
// 		}

// 		medicationResponses = append(medicationResponses, PrescriptionMedicationResponse{
// 			ID:             med.ID,
// 			PrescriptionID: prescriptionID,
// 			MedicineID:     med.MedicineID,
// 			MedicineName:   medicine.Name,
// 			Dosage:         med.Dosage,
// 			Frequency:      med.Frequency,
// 			Duration:       med.Duration,
// 			Instructions:   med.Instructions.String,
// 		})
// 	}

// 	layout := "2006-01-02 15:04:05"
// 	return &PrescriptionResponse{
// 		ID:               prescription.ID,
// 		MedicalHistoryID: prescription.MedicalHistoryID,
// 		ExaminationID:    prescription.ExaminationID,
// 		PrescriptionDate: prescription.PrescriptionDate.Time.Format(layout),
// 		DoctorID:         prescription.DoctorID,
// 		DoctorName:       doctor.Name,
// 		Notes:            prescription.Notes.String,
// 		Medications:      medicationResponses,
// 		CreatedAt:        prescription.CreatedAt.Time.Format(layout),
// 		UpdatedAt:        prescription.UpdatedAt.Time.Format(layout),
// 	}, nil
// }

// // ListPrescriptionsByMedicalHistory lists all prescriptions for a medical history with pagination
// func (s *MedicalRecordService) ListPrescriptionsByMedicalHistory(ctx *gin.Context, medicalHistoryID int64, pagination *util.Pagination) ([]PrescriptionResponse, error) {
// 	offset := (pagination.Page - 1) * pagination.PageSize

// 	prescriptions, err := s.storeDB.ListPrescriptionsByMedicalHistory(ctx, db.ListPrescriptionsByMedicalHistoryParams{
// 		MedicalHistoryID: medicalHistoryID,
// 		Limit:            int32(pagination.PageSize),
// 		Offset:           int32(offset),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list prescriptions: %w", err)
// 	}

// 	var results []PrescriptionResponse
// 	// layout := "2006-01-02 15:04:05"

// 	for _, prescription := range prescriptions {
// 		// Get prescription with full details
// 		prescriptionDetail, err := s.GetPrescription(ctx, prescription.ID)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get prescription details: %w", err)
// 		}

// 		results = append(results, *prescriptionDetail)
// 	}

// 	return results, nil
// }

// // ListPrescriptionsByPet lists all prescriptions for a pet with pagination
// func (s *MedicalRecordService) ListPrescriptionsByPet(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PrescriptionResponse, error) {
// 	// Get medical histories for pet
// 	histories, err := s.GetMedicalHistoryByPetID(ctx, petID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get medical histories: %w", err)
// 	}

// 	if len(histories) == 0 {
// 		return []PrescriptionResponse{}, nil
// 	}

// 	// Get prescriptions for each medical history and combine them
// 	var allPrescriptions []PrescriptionResponse
// 	for _, history := range histories {
// 		prescriptions, err := s.ListPrescriptionsByMedicalHistory(ctx, history.ID, pagination)
// 		if err != nil {
// 			// Log error but continue
// 			log.Printf("Error fetching prescriptions for history ID %d: %v", history.ID, err)
// 			continue
// 		}
// 		allPrescriptions = append(allPrescriptions, prescriptions...)
// 	}

// 	return allPrescriptions, nil
// }

// Test Result-related service methods

// // CreateTestResult creates a new test result record
// func (s *MedicalRecordService) CreateTestResult(ctx *gin.Context, req *TestResultRequest) (*TestResultResponse, error) {
// 	var testResultResponse TestResultResponse

// 	// Parse the test date
// 	layout := "2006-01-02 15:04:05"
// 	testDate, err := time.Parse(layout, req.TestDate)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid test date format: %w", err)
// 	}

// 	// Begin transaction
// 	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		// Verify that medical history exists
// 		_, err := q.GetMedicalHistoryById(ctx, req.MedicalHistoryID)
// 		if err != nil {
// 			return fmt.Errorf("failed to get medical history: %w", err)
// 		}

// 		// Verify that examination exists
// 		examination, err := q.GetExaminationByID(ctx, req.ExaminationID)
// 		if err != nil {
// 			return fmt.Errorf("failed to get examination: %w", err)
// 		}

// 		// Verify that examination belongs to this medical history
// 		if examination.MedicalHistoryID != req.MedicalHistoryID {
// 			return fmt.Errorf("examination does not belong to the specified medical history")
// 		}

// 		// Verify that doctor exists
// 		doctor, err := q.GetDoctor(ctx, req.DoctorID)
// 		if err != nil {
// 			return fmt.Errorf("failed to get doctor: %w", err)
// 		}

// 		// Create test result
// 		testResult, err := q.CreateTestResult(ctx, db.CreateTestResultParams{
// 			MedicalHistoryID: req.MedicalHistoryID,
// 			ExaminationID:    req.ExaminationID,
// 			TestType:         req.TestType,
// 			TestDate:         pgtype.Timestamp{Time: testDate, Valid: true},
// 			Results:          req.Results,
// 			Interpretation:   pgtype.Text{String: req.Interpretation, Valid: true},
// 			FileUrl:          pgtype.Text{String: req.FileURL, Valid: true},
// 			DoctorID:         req.DoctorID,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to create test result: %w", err)
// 		}

// 		// Populate response
// 		testResultResponse = TestResultResponse{
// 			ID:               testResult.ID,
// 			MedicalHistoryID: testResult.MedicalHistoryID,
// 			ExaminationID:    testResult.ExaminationID,
// 			TestType:         testResult.TestType,
// 			TestDate:         testResult.TestDate.Time.Format(layout),
// 			Results:          testResult.Results,
// 			Interpretation:   testResult.Interpretation.String,
// 			FileURL:          testResult.FileUrl.String,
// 			DoctorID:         testResult.DoctorID,
// 			DoctorName:       doctor.Name,
// 			CreatedAt:        testResult.CreatedAt.Time.Format(layout),
// 			UpdatedAt:        testResult.UpdatedAt.Time.Format(layout),
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &testResultResponse, nil
// }

// GetTestResult retrieves a test result by ID
// func (s *MedicalRecordService) GetTestResult(ctx *gin.Context, testResultID int64) (*TestResultResponse, error) {
// 	testResult, err := s.storeDB.GetTestResultByID(ctx, testResultID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get test result: %w", err)
// 	}

// 	// Get doctor info
// 	doctor, err := s.storeDB.GetDoctor(ctx, testResult.DoctorID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get doctor: %w", err)
// 	}

// 	layout := "2006-01-02 15:04:05"
// 	return &TestResultResponse{
// 		ID:               testResult.ID,
// 		MedicalHistoryID: testResult.MedicalHistoryID,
// 		ExaminationID:    testResult.ExaminationID,
// 		TestType:         testResult.TestType,
// 		TestDate:         testResult.TestDate.Time.Format(layout),
// 		Results:          testResult.Results,
// 		Interpretation:   testResult.Interpretation.String,
// 		FileURL:          testResult.FileUrl.String,
// 		DoctorID:         testResult.DoctorID,
// 		DoctorName:       doctor.Name,
// 		CreatedAt:        testResult.CreatedAt.Time.Format(layout),
// 		UpdatedAt:        testResult.UpdatedAt.Time.Format(layout),
// 	}, nil
// }

// ListTestResultsByMedicalHistory lists all test results for a medical history with pagination
// func (s *MedicalRecordService) ListTestResultsByMedicalHistory(ctx *gin.Context, medicalHistoryID int64, pagination *util.Pagination) ([]TestResultResponse, error) {
// 	offset := (pagination.Page - 1) * pagination.PageSize

// 	testResults, err := s.storeDB.ListTestResultsByMedicalHistory(ctx, db.ListTestResultsByMedicalHistoryParams{
// 		MedicalHistoryID: medicalHistoryID,
// 		Limit:            int32(pagination.PageSize),
// 		Offset:           int32(offset),
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list test results: %w", err)
// 	}

// 	var results []TestResultResponse
// 	layout := "2006-01-02 15:04:05"

// 	for _, testResult := range testResults {
// 		// Get doctor info
// 		doctor, err := s.storeDB.GetDoctor(ctx, testResult.DoctorID)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get doctor: %w", err)
// 		}

// 		results = append(results, TestResultResponse{
// 			ID:               testResult.ID,
// 			MedicalHistoryID: testResult.MedicalHistoryID,
// 			ExaminationID:    testResult.ExaminationID,
// 			TestType:         testResult.TestType,
// 			TestDate:         testResult.TestDate.Time.Format(layout),
// 			Results:          testResult.Results,
// 			Interpretation:   testResult.Interpretation.String,
// 			FileURL:          testResult.FileUrl.String,
// 			DoctorID:         testResult.DoctorID,
// 			DoctorName:       doctor.Name,
// 			CreatedAt:        testResult.CreatedAt.Time.Format(layout),
// 			UpdatedAt:        testResult.UpdatedAt.Time.Format(layout),
// 		})
// 	}

// 	return results, nil
// }

// // ListTestResultsByPet lists all test results for a pet with pagination
// func (s *MedicalRecordService) ListTestResultsByPet(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]TestResultResponse, error) {
// 	// Get medical histories for pet
// 	histories, err := s.GetMedicalHistoryByPetID(ctx, petID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get medical histories: %w", err)
// 	}

// 	if len(histories) == 0 {
// 		return []TestResultResponse{}, nil
// 	}

// 	// Get test results for each medical history and combine them
// 	var allTestResults []TestResultResponse
// 	for _, history := range histories {
// 		testResults, err := s.ListTestResultsByMedicalHistory(ctx, history.ID, pagination)
// 		if err != nil {
// 			// Log error but continue
// 			log.Printf("Error fetching test results for history ID %d: %v", history.ID, err)
// 			continue
// 		}
// 		allTestResults = append(allTestResults, testResults...)
// 	}

// 	return allTestResults, nil
// }

// // Comprehensive medical history
// // GetCompleteMedicalHistory combines all medical history data for a pet in a single response
// func (s *MedicalRecordService) GetCompleteMedicalHistory(ctx *gin.Context, petID int64) (*MedicalHistorySummary, error) {
// 	// Get or create medical record
// 	medicalRecord, err := s.GetMedicalRecord(ctx, petID)
// 	if err != nil {
// 		// If medical record doesn't exist, try to create one
// 		medicalRecord, err = s.CreateMedicalRecord(ctx, petID)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get or create medical record: %w", err)
// 		}
// 	}

// 	// Use unlimited pagination for comprehensive view
// 	unlimitedPagination := &util.Pagination{
// 		Page:     1,
// 		PageSize: 1000, // Large number to get all records
// 	}

// 	// Get all medical history entries
// 	conditions, err := s.GetMedicalHistoryByPetID(ctx, petID)
// 	if err != nil {
// 		log.Printf("Error fetching medical conditions: %v", err)
// 		// Don't return error, just provide empty list
// 		conditions = []MedicalHistoryResponse{}
// 	}

// 	// Get all examinations
// 	examinations, err := s.ListExaminationsByPet(ctx, petID, unlimitedPagination)
// 	if err != nil {
// 		log.Printf("Error fetching examinations: %v", err)
// 		// Don't return error, just provide empty list
// 		examinations = []ExaminationResponse{}
// 	}

// 	// Get all prescriptions
// 	prescriptions, err := s.ListPrescriptionsByPet(ctx, petID, unlimitedPagination)
// 	if err != nil {
// 		log.Printf("Error fetching prescriptions: %v", err)
// 		// Don't return error, just provide empty list
// 		prescriptions = []PrescriptionResponse{}
// 	}

// 	// Get all test results
// 	testResults, err := s.ListTestResultsByPet(ctx, petID, unlimitedPagination)
// 	if err != nil {
// 		log.Printf("Error fetching test results: %v", err)
// 		// Don't return error, just provide empty list
// 		testResults = []TestResultResponse{}
// 	}

// 	// Construct comprehensive medical history
// 	summary := &MedicalHistorySummary{
// 		MedicalRecord: *medicalRecord,
// 		Examinations:  examinations,
// 		Prescriptions: prescriptions,
// 		TestResults:   testResults,
// 		Conditions:    conditions,
// 	}

// 	return summary, nil
// }

// Get all soap notes by pet
func (s *MedicalRecordService) GetAllSoapNotesByPetID(ctx *gin.Context, petID int64) ([]SoapNoteResponse, error) {
	// Get medical histories for pet
	histories, err := s.storeDB.ListSOAPNotes(ctx, pgtype.Int8{Int64: petID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get medical histories: %w", err)
	}

	if len(histories) == 0 {
		return []SoapNoteResponse{}, nil
	}

	// Get SOAP notes for each medical history and combine them
	var allSoapNotes []SoapNoteResponse
	for _, history := range histories {
		allSoapNotes = append(allSoapNotes, SoapNoteResponse{
			ID:               int64(history.ID),
			PetID:            history.PetID.Int64,
			Subjective:       string(history.Subjective),
			Objective:        string(history.Objective), // Convert []byte to string
			Assessment:       string(history.Assessment),
			Plan:             fmt.Sprint(history.Plan.Int64), // Convert int64 to string
			DoctorID:         history.DoctorID,
			DoctorName:       history.DoctorName,
			ConsultationDate: history.ConsultationDate.Time.Format("2006-01-02 15:04:05"),
			CreatedAt:        history.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	return allSoapNotes, nil
}

// GetAppointmentVisitSummary tổng hợp tất cả thông tin của một cuộc hẹn/lần khám
func (s *MedicalRecordService) GetAppointmentVisitSummary(ctx *gin.Context, appointmentID int64) (*AppointmentVisitSummaryResponse, error) {
	// 1. Lấy thông tin chi tiết của cuộc hẹn
	appointmentDetail, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy thông tin cuộc hẹn: %w", err)
	}

	// 2. Tạo đối tượng phản hồi
	summary := &AppointmentVisitSummaryResponse{
		AppointmentID:     appointmentID,
		VisitDate:         appointmentDetail.Date.Time,
		AppointmentStatus: appointmentDetail.StateName.String,
		Reason:            appointmentDetail.AppointmentReason.String,
	}

	// 3. Điền thông tin thú cưng
	summary.PetInfo = PetBasicInfo{
		PetID: appointmentDetail.PetID.Int64,
		Name:  appointmentDetail.PetName.String,
		Breed: appointmentDetail.PetBreed.String,
	}

	// Lấy thêm thông tin chi tiết về thú cưng
	if appointmentDetail.Petid.Valid {
		pet, err := s.storeDB.GetPetByID(ctx, appointmentDetail.Petid.Int64)
		if err == nil {
			// Bổ sung thêm thông tin
			summary.PetInfo.Type = pet.Type
			summary.PetInfo.Age = pet.Age.Int32
			summary.PetInfo.Gender = pet.Gender.String
			summary.PetInfo.Weight = pet.Weight.Float64
		}
	}

	// 4. Điền thông tin chủ sở hữu
	summary.OwnerInfo = OwnerBasicInfo{
		Username: appointmentDetail.Username.String,
		FullName: appointmentDetail.OwnerName.String,
		Phone:    appointmentDetail.OwnerPhone.String,
		Email:    appointmentDetail.OwnerEmail.String,
	}

	// 5. Điền thông tin bác sĩ - sử dụng DoctorName từ doctor_id
	// Lấy thông tin bác sĩ từ bảng doctors dựa trên doctor_id
	doctor, err := s.getDoctorInfo(ctx, appointmentDetail.DoctorID.Int64)
	if err == nil {
		summary.DoctorInfo = DoctorBasicInfo{
			DoctorID: appointmentDetail.DoctorID.Int64,
			Name:     doctor,
		}
	} else {
		summary.DoctorInfo = DoctorBasicInfo{
			DoctorID: appointmentDetail.DoctorID.Int64,
			Name:     "Unknown Doctor", // Giá trị mặc định
		}
	}

	// 6. Điền thông tin dịch vụ
	service, err := s.getServiceInfo(ctx, appointmentDetail.ServiceID.Int64)
	if err == nil {
		summary.ServiceInfo = ServiceBasicInfo{
			ServiceID: appointmentDetail.ServiceID.Int64,
			Name:      service.Name,
			Duration:  int32(service.Duration), // Chuyển đổi kiểu dữ liệu
			Price:     service.Price,
		}
	} else {
		summary.ServiceInfo = ServiceBasicInfo{
			ServiceID: appointmentDetail.ServiceID.Int64,
			Name:      "Unknown Service",
		}
	}

	// 7. Lấy thông tin SOAP Notes nếu có
	if appointmentDetail.PetID.Valid {
		soapNote, err := s.storeDB.GetSOAPByAppointmentID(ctx, pgtype.Int8{Int64: appointmentID, Valid: true})
		if err == nil {
			summary.SOAPNote = &SOAPNoteInfo{
				ID:         int32(soapNote.ID),
				Subjective: string(soapNote.Subjective),
				Assessment: string(soapNote.Assessment),
				CreatedAt:  soapNote.CreatedAt.Time,
			}

			// Xử lý dữ liệu Objective từ dạng []byte
			if soapNote.Objective != nil {
				summary.SOAPNote.Objective = string(soapNote.Objective)
			}

			// Xử lý Plan (có thể là ID trỏ đến bảng khác hoặc văn bản)
			if soapNote.Plan.Valid {
				summary.SOAPNote.Plan = fmt.Sprintf("%d", soapNote.Plan.Int64)
			}
		}
	}

	// 8. Lấy thông tin điều trị
	if appointmentDetail.PetID.Valid {
		// Lấy các treatment liên quan đến pet và được tạo vào ngày khám
		arg := db.GetTreatmentsByPetParams{
			PetID:  pgtype.Int8{Int64: appointmentDetail.PetID.Int64, Valid: true},
			Limit:  10,
			Offset: 0,
		}

		treatments, err := s.storeDB.GetTreatmentsByPet(ctx, arg)
		if err == nil {
			for _, treatment := range treatments {
				// Kiểm tra xem treatment có được tạo trong khoảng thời gian của cuộc hẹn
				// hoặc có cùng ngày với cuộc hẹn hay không
				if treatment.CreatedAt.Time.Format("2006-01-02") == appointmentDetail.Date.Time.Format("2006-01-02") {
					treatmentInfo := TreatmentBasicInfo{
						ID:          treatment.ID,
						Name:        treatment.Name.String,
						Type:        treatment.Type.String,
						Status:      treatment.Status.String,
						StartDate:   treatment.StartDate.Time,
						EndDate:     treatment.EndDate.Time,
						Description: treatment.Description.String,
					}
					summary.Treatments = append(summary.Treatments, treatmentInfo)

				}
			}
		}
	}

	// 9. Lấy thông tin đơn thuốc
	if appointmentDetail.PetID.Valid {
		// Lấy tất cả đơn thuốc của thú cưng được kê vào ngày khám
		arg := db.ListMedicinesByPetParams{
			PetID:  pgtype.Int8{Int64: appointmentDetail.PetID.Int64, Valid: true},
			Status: pgtype.Text{String: "active", Valid: true},
			Limit:  10,
			Offset: 0,
		}

		medications, err := s.storeDB.ListMedicinesByPet(ctx, arg)
		if err == nil {
			for _, med := range medications {
				// Kiểm tra xem thuốc có được kê vào ngày của cuộc hẹn không
				if med.CreatedAt.Time.Format("2006-01-02") == appointmentDetail.Date.Time.Format("2006-01-02") {
					prescriptionInfo := PrescriptionBasicInfo{
						ID:           med.ID,
						MedicineName: med.Name,
						Dosage:       med.Dosage.String,
						Frequency:    med.Frequency.String,
						Duration:     med.Duration.String,
						Quantity:     int32(med.Quantity.Int64),
						IssuedDate:   med.CreatedAt.Time,
					}
					summary.Prescriptions = append(summary.Prescriptions, prescriptionInfo)
				}
			}
		}
	}

	// 10. Lấy thông tin kết quả xét nghiệm - giả định hàm này tồn tại
	if appointmentDetail.PetID.Valid {
		// Sử dụng các hàm hiện có hoặc sử dụng truy vấn SQL riêng
		testResults, err := s.getTestResultsByPet(ctx, appointmentDetail.PetID.Int64, appointmentDetail.Date.Time)
		if err == nil {
			summary.TestResults = testResults
		}
	}

	// 11. Tìm cuộc hẹn tiếp theo
	// Tìm cuộc hẹn có ngày lớn hơn cuộc hẹn hiện tại và cùng pet
	if appointmentDetail.PetID.Valid {
		// Để đơn giản, chúng ta sẽ giả định lấy tất cả cuộc hẹn và lọc thủ công
		appointments, err := s.storeDB.ListAllAppointments(ctx)
		if err == nil {
			var nextAppointment *db.Appointment
			for _, appt := range appointments {
				if appt.Petid.Int64 == appointmentDetail.PetID.Int64 &&
					appt.Date.Time.After(appointmentDetail.Date.Time) &&
					(nextAppointment == nil || appt.Date.Time.Before(nextAppointment.Date.Time)) {
					temp := appt // Tạo bản sao để tránh vấn đề với biến vòng lặp
					nextAppointment = &temp
				}
			}

			if nextAppointment != nil {
				// Lấy thêm thông tin về dịch vụ và bác sĩ cho cuộc hẹn tiếp theo
				nextServiceName := s.getNextServiceName(ctx, nextAppointment.ServiceID.Int64)
				nextDoctorName := s.getNextDoctorName(ctx, nextAppointment.DoctorID.Int64)

				summary.NextAppointment = &NextAppointmentInfo{
					AppointmentID: nextAppointment.AppointmentID,
					Date:          nextAppointment.Date.Time,
					ServiceName:   nextServiceName,
					DoctorName:    nextDoctorName,
				}
			}
		}
	}

	// // 12. Lấy dấu hiệu sinh tồn (nếu có)
	// // Nếu có bảng dữ liệu riêng cho dấu hiệu sinh tồn, chúng ta sẽ truy vấn ở đây
	// // Hiện tại chúng ta sẽ sử dụng cân nặng từ thông tin pet
	// if appointmentDetail.PetID.Valid {
	// 	pet, err := s.storeDB.GetPetByID(ctx, appointmentDetail.PetID.Int64)
	// 	if err == nil && pet.Weight.Valid {
	// 		summary.VitalSigns = &VitalSignsInfo{
	// 			Weight: pet.Weight.Float64,
	// 			// Các dấu hiệu sinh tồn khác có thể được thêm vào nếu có dữ liệu
	// 		}
	// 	}
	// }

	return summary, nil
}

// Helper functions to abstract complex data retrieval
func (s *MedicalRecordService) getDoctorInfo(ctx *gin.Context, doctorID int64) (string, error) {
	// TODO: Triển khai thực tế khi có DB schema
	// Giả lập response
	return "Dr. John Doe", nil
}

func (s *MedicalRecordService) getServiceInfo(ctx *gin.Context, serviceID int64) (*struct {
	Name     string
	Duration int
	Price    float64
}, error) {
	// TODO: Triển khai thực tế khi có DB schema
	// Giả lập response
	return &struct {
		Name     string
		Duration int
		Price    float64
	}{
		Name:     "Regular Checkup",
		Duration: 30,
		Price:    50.00,
	}, nil
}

func (s *MedicalRecordService) getTestResultsByPet(ctx *gin.Context, petID int64, visitDate time.Time) ([]TestResultBasicInfo, error) {
	// TODO: Triển khai thực tế khi có DB schema
	// Giả lập response
	results := []TestResultBasicInfo{
		{
			ID:       1,
			TestName: "Blood Test",
			Result:   "Normal",
			Status:   "Completed",
			TestDate: visitDate,
		},
	}
	return results, nil
}

func (s *MedicalRecordService) getNextServiceName(ctx *gin.Context, serviceID int64) string {
	// TODO: Triển khai thực tế khi có DB schema
	return "Follow-up Visit"
}

func (s *MedicalRecordService) getNextDoctorName(ctx *gin.Context, doctorID int64) string {
	// TODO: Triển khai thực tế khi có DB schema
	return "Dr. Jane Smith"
}
