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
	Treatment       int64  `json:"treatment"`
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
