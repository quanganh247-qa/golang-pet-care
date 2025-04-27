package medical_records

import (
	"errors"
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

var (
	ErrMedicalRecordNotFound = errors.New("medical record not found")
	ErrInvalidDiagnosisDate  = errors.New("invalid diagnosis date format")
	ErrInvalidPetID          = errors.New("invalid pet ID")
	ErrPrescriptionNotFound  = errors.New("prescription not found")
	ErrTestResultNotFound    = errors.New("test result not found")
)

type MedicalRecordApi struct {
	controller MedicalRecordControllerInterface
}

type MedicalRecordController struct {
	service MedicalRecordServiceInterface
}

type MedicalRecordService struct {
	storeDB db.Store
}
type MedicalRecord struct {
	ID        string    `json:"id"`
	PetID     string    `json:"pet_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MedicalRecordResponse struct {
	ID        int64  `json:"id"`
	PetID     int64  `json:"pet_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MedicalHistoryRequest struct {
	Condition     string `json:"condition" binding:"required"`
	DiagnosisDate string `json:"diagnosis_date" binding:"required,datetime=2006-01-02 15:04:05"`
	Notes         string `json:"notes"`
}

type MedicalHistoryResponse struct {
	ID              int64  `json:"id"`
	MedicalRecordID int64  `json:"medical_record_id"`
	Condition       string `json:"condition"`
	DiagnosisDate   string `json:"diagnosis_date"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type Allergy struct {
	ID              int64  `json:"id"`
	MedicalRecordID int64  `json:"medical_record_id"`
	Allergen        string `json:"allergen"`
	Severity        string `json:"severity"`
	Reaction        string `json:"reaction"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type AllergyRequest struct {
	Allergen string `json:"allergen" binding:"required"`
	Severity string `json:"severity" binding:"required,oneof=mild moderate severe"`
	Reaction string `json:"reaction" binding:"required"`
	Notes    string `json:"notes"`
}

// New structures for complete medical history

// Examination represents a medical examination for a pet
type ExaminationRequest struct {
	MedicalHistoryID int64  `json:"medical_history_id" binding:"required"`
	ExamDate         string `json:"exam_date" binding:"required,datetime=2006-01-02 15:04:05"`
	ExamType         string `json:"exam_type" binding:"required"`
	Findings         string `json:"findings" binding:"required"`
	VetNotes         string `json:"vet_notes"`
	DoctorID         int64  `json:"doctor_id" binding:"required"`
}

type ExaminationResponse struct {
	ID               int64  `json:"id"`
	MedicalHistoryID int64  `json:"medical_history_id"`
	ExamDate         string `json:"exam_date"`
	ExamType         string `json:"exam_type"`
	Findings         string `json:"findings"`
	VetNotes         string `json:"vet_notes"`
	DoctorID         int64  `json:"doctor_id"`
	DoctorName       string `json:"doctor_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// Prescription represents a medication prescription for a pet
type PrescriptionRequest struct {
	MedicalHistoryID int64                  `json:"medical_history_id" binding:"required"`
	ExaminationID    int64                  `json:"examination_id" binding:"required"`
	PrescriptionDate string                 `json:"prescription_date" binding:"required,datetime=2006-01-02 15:04:05"`
	DoctorID         int64                  `json:"doctor_id" binding:"required"`
	Notes            string                 `json:"notes"`
	Medications      []PrescribedMedication `json:"medications" binding:"required"`
}

type PrescribedMedication struct {
	MedicineID   int64  `json:"medicine_id" binding:"required"`
	Dosage       string `json:"dosage" binding:"required"`
	Frequency    string `json:"frequency" binding:"required"`
	Duration     string `json:"duration" binding:"required"`
	Instructions string `json:"instructions"`
}

type PrescriptionResponse struct {
	ID               int64                            `json:"id"`
	MedicalHistoryID int64                            `json:"medical_history_id"`
	ExaminationID    int64                            `json:"examination_id"`
	PrescriptionDate string                           `json:"prescription_date"`
	DoctorID         int64                            `json:"doctor_id"`
	DoctorName       string                           `json:"doctor_name"`
	Notes            string                           `json:"notes"`
	Medications      []PrescriptionMedicationResponse `json:"medications"`
	CreatedAt        string                           `json:"created_at"`
	UpdatedAt        string                           `json:"updated_at"`
}

type PrescriptionMedicationResponse struct {
	ID             int64  `json:"id"`
	PrescriptionID int64  `json:"prescription_id"`
	MedicineID     int64  `json:"medicine_id"`
	MedicineName   string `json:"medicine_name"`
	Dosage         string `json:"dosage"`
	Frequency      string `json:"frequency"`
	Duration       string `json:"duration"`
	Instructions   string `json:"instructions"`
}

// TestResult represents test results for a pet
type TestResultRequest struct {
	MedicalHistoryID int64  `json:"medical_history_id" binding:"required"`
	ExaminationID    int64  `json:"examination_id" binding:"required"`
	TestType         string `json:"test_type" binding:"required"`
	TestDate         string `json:"test_date" binding:"required,datetime=2006-01-02 15:04:05"`
	Results          string `json:"results" binding:"required"`
	Interpretation   string `json:"interpretation"`
	FileURL          string `json:"file_url"`
	DoctorID         int64  `json:"doctor_id" binding:"required"`
}

type TestResultResponse struct {
	ID               int64  `json:"id"`
	MedicalHistoryID int64  `json:"medical_history_id"`
	ExaminationID    int64  `json:"examination_id"`
	TestType         string `json:"test_type"`
	TestDate         string `json:"test_date"`
	Results          string `json:"results"`
	Interpretation   string `json:"interpretation"`
	FileURL          string `json:"file_url"`
	DoctorID         int64  `json:"doctor_id"`
	DoctorName       string `json:"doctor_name"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// MedicalHistorySummary represents a comprehensive view of medical history
type MedicalHistorySummary struct {
	MedicalRecord MedicalRecordResponse    `json:"medical_record"`
	Examinations  []ExaminationResponse    `json:"examinations"`
	Prescriptions []PrescriptionResponse   `json:"prescriptions"`
	TestResults   []TestResultResponse     `json:"test_results"`
	Conditions    []MedicalHistoryResponse `json:"conditions"`
	Allergies     []Allergy                `json:"allergies"`
}
