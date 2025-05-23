package disease

import (
	"time"

	"github.com/lib/pq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/llm"
)

var (
	hospitalLogo    = "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"
	hospitalName    = "Aewan Clinic"
	hospitalAddress = ""
	hospitalPhone   = "0123-456-789"
)

type Prescription struct {
	ID              string         `json:"id"`
	HospitalLogo    string         `json:"hospital_logo"`
	HospitalName    string         `json:"hospital_name"`
	HospitalAddress string         `json:"hospital_address"`
	HospitalPhone   string         `json:"hospital_phone"`
	PatientName     string         `json:"patient_name"`
	PatientGender   string         `json:"patient_gender"`
	PatientAge      int            `json:"patient_age"`
	Diagnosis       string         `json:"diagnosis"`
	Medicines       []MedicineInfo `json:"medicines"`
	Notes           string         `json:"notes"`
	PrescribedDate  time.Time      `json:"prescribed_date"`
	DoctorName      string         `json:"doctor_name"`
	DoctorTitle     string         `json:"doctor_title"`
}

type DiseaseApi struct {
	controller DiseaseControllerInterface
}

type DiseaseController struct {
	service    DiseaseServiceInterface
	llmService llm.AIService
}

type DiseaseService struct {
	storeDB db.Store
}

// DiseaseMedicineInfo holds the disease and associated medicine information.
type DiseaseMedicineInfo struct {
	DiseaseID          int            `json:"disease_id"`
	DiseaseName        string         `json:"disease_name"`
	DiseaseDescription *string        `json:"disease_description"`
	Symptoms           pq.StringArray `json:"symptoms"`
	Medicines          []MedicineInfo `json:"medicines"`
}

// MedicineInfo holds medicine-related details for a specific disease.
type MedicineInfo struct {
	MedicineID   int    `json:"medicine_id"`
	MedicineName string `json:"medicine_name"`
	Usage        string `json:"usage"`
	Dosage       string `json:"dosage"`
	Frequency    string `json:"frequency"`
	Duration     string `json:"duration"`
	SideEffects  string `json:"side_effects"`
	Notes        string `json:"notes"`
}

type CreateTreatmentRequest struct {
	PetID     int64  `json:"pet_id"`
	Name      string `json:"name" validate:"required"`
	DoctorID  int64  `json:"doctor_id"`
	Type      string `json:"type" validate:"required"`
	Diseases  string `json:"diseases"`
	StartDate string `json:"start_date" validate:"required, datetime=2006-01-02"` // have lay out for date
	Status    string `json:"status" validate:"required"`
	Notes     string `json:"notes"`
}

type Treatment struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Type        string           `json:"type"`
	Disease     string           `json:"disease"`
	StartDate   string           `json:"start_date"`
	EndDate     string           `json:"end_date"`
	Status      string           `json:"status"`
	Description string           `json:"description"`
	DoctorName  string           `json:"doctor_name"`
	Phases      []TreatmentPhase `json:"phases"`
}

type CreateTreatmentResponse struct {
	TreatmentID int64  `json:"treatment_id"`
	Disease     string `json:"disease"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Status      string `json:"status"`
	DoctorName  string `json:"doctor_name"`
}

type CreateTreatmentPhaseRequest struct {
	PhaseName   string `json:"phase_name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	Status      string `json:"status" validate:"required,oneof=pending active completed"`
}

// Assign Medicines to Treatment Phases
type AssignMedicineRequest struct {
	MedicineID    int64  `json:"medicine_id"`
	Dosage        string `json:"dosage"`
	Frequency     string `json:"frequency"`
	Duration      string `json:"duration"`
	Quantity      int64  `json:"quantity"`
	Notes         string `json:"notes"`
	AppointmentID int64  `json:"appointment_id"`
}

type TreatmentPhase struct {
	ID          int64           `json:"id"`
	TreatmentID int64           `json:"treatment_id"`
	PhaseName   string          `json:"phase_name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	StartDate   string          `json:"start_date"`
	CreatedAt   string          `json:"created_at"`
	Medications []PhaseMedicine `json:"medications"`
}

type PhaseMedicine struct {
	PhaseID      int64  `json:"phase_id"`
	MedicineID   int64  `json:"medicine_id"`
	MedicineName string `json:"medicine_name"`
	Dosage       string `json:"dosage"`
	Frequency    string `json:"frequency"`
	Duration     string `json:"duration"`
	Quantity     int64  `json:"quantity"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"created_at"`
}

type UpdateTreatmentPhaseStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending active completed"`
}

type UpdateTreatmentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending active completed"`
}

type TreatmentProgressDetail struct {
	PhaseName    string `json:"phase_name"`
	Status       string `json:"status"`
	StartDate    string `json:"start_date"`
	NumMedicines int32  `json:"num_medicines"`
}

type CreateDiseaseRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Symptoms    []string `json:"symptoms"`
}

type PrescriptionResponse struct {
	PrescriptionID  int64  `json:"prescription_id"`
	PrescriptionURL string `json:"prescription_url"`
}

type PetAllergy struct {
	ID     int64  `json:"id"`
	PetID  int64  `json:"pet_id"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

type CreateAllergyRequest struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

// Add this if it doesn't exist
type MedicineResponse struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Usage          string `json:"usage,omitempty"`
	Dosage         string `json:"dosage,omitempty"`
	Frequency      string `json:"frequency,omitempty"`
	Duration       string `json:"duration,omitempty"`
	SideEffects    string `json:"side_effects,omitempty"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty"`
	Quantity       int64  `json:"quantity,omitempty"`
}

type SymptomAnalysisRequest struct {
	PatientID int64  `json:"patient_id"`
	Symptoms  string `json:"symptoms"`
}
