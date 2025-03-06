package medical_records

import (
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"errors"
=======
>>>>>>> 3bf345d (happy new year)
=======
	"errors"
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> 3bf345d (happy new year)
=======
	"errors"
>>>>>>> e859654 (Elastic search)
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> e859654 (Elastic search)
var (
	ErrMedicalRecordNotFound = errors.New("medical record not found")
	ErrInvalidDiagnosisDate  = errors.New("invalid diagnosis date format")
	ErrInvalidPetID          = errors.New("invalid pet ID")
)

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	Condition     string `json:"condition" binding:"required"`
	DiagnosisDate string `json:"diagnosis_date" binding:"required,datetime=2006-01-02 15:04:05"`
	Treatment     int64  `json:"treatment" binding:"required"`
=======
	Condition     string `json:"condition"`
	DiagnosisDate string `json:"diagnosis_date"`
	Treatment     int64  `json:"treatment"`
>>>>>>> 3bf345d (happy new year)
=======
	Condition     string `json:"condition" binding:"required"`
	DiagnosisDate string `json:"diagnosis_date" binding:"required,datetime=2006-01-02 15:04:05"`
	Treatment     int64  `json:"treatment" binding:"required"`
>>>>>>> e859654 (Elastic search)
=======
	Condition     string `json:"condition"`
	DiagnosisDate string `json:"diagnosis_date"`
	Treatment     int64  `json:"treatment"`
>>>>>>> 3bf345d (happy new year)
=======
	Condition     string `json:"condition" binding:"required"`
	DiagnosisDate string `json:"diagnosis_date" binding:"required,datetime=2006-01-02 15:04:05"`
	Treatment     int64  `json:"treatment" binding:"required"`
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	ID              int64     `json:"id"`
	MedicalRecordID int64     `json:"medical_record_id"`
=======
	ID              string    `json:"id"`
	MedicalRecordID string    `json:"medical_record_id"`
>>>>>>> 3bf345d (happy new year)
=======
	ID              int64     `json:"id"`
	MedicalRecordID int64     `json:"medical_record_id"`
>>>>>>> e859654 (Elastic search)
=======
	ID              string    `json:"id"`
	MedicalRecordID string    `json:"medical_record_id"`
>>>>>>> 3bf345d (happy new year)
=======
	ID              int64     `json:"id"`
	MedicalRecordID int64     `json:"medical_record_id"`
>>>>>>> e859654 (Elastic search)
	Allergen        string    `json:"allergen"`
	Severity        string    `json:"severity"`
	Reaction        string    `json:"reaction"`
	Notes           string    `json:"notes"`
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> e859654 (Elastic search)
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

type AllergyRequest struct {
	Allergen string `json:"allergen" binding:"required"`
	Severity string `json:"severity" binding:"required,oneof=mild moderate severe"`
	Reaction string `json:"reaction" binding:"required"`
	Notes    string `json:"notes"`
<<<<<<< HEAD
<<<<<<< HEAD
=======
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
=======
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
}
