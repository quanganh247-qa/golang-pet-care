package disease

import (
	"github.com/lib/pq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
)

type DiseaseApi struct {
	controller DiseaseControllerInterface
}

type DiseaseController struct {
	service DiseaseServiceInterface
}

type DiseaseService struct {
	storeDB db.Store
	es      *elasticsearch.ESService
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
	MedicineID   int     `json:"medicine_id"`
	MedicineName string  `json:"medicine_name"`
	Usage        *string `json:"usage"`
	Dosage       *string `json:"dosage"`
	Frequency    *string `json:"frequency"`
	Duration     *string `json:"duration"`
	SideEffects  *string `json:"side_effects"`
}

type TreatmentPlan struct {
	DiseaseID       int           `json:"disease_id"`
	DiseaseName     string        `json:"disease_name"`
	Description     string        `json:"description"`
	Symptoms        []string      `json:"symptoms"`
	TreatmentPhases []PhaseDetail `json:"treatment_phases"`
}

type PhaseDetail struct {
	PhaseID          int            `json:"phase_id"`
	PhaseNumber      int            `json:"phase_number"`
	PhaseName        string         `json:"phase_name"`
	Duration         string         `json:"duration"`
	PhaseDescription string         `json:"phase_description"`
	PhaseNotes       string         `json:"phase_notes"`
	Medicines        []MedicineInfo `json:"medicines"`
}

type CreateTreatmentRequest struct {
	PetID     int64  `json:"pet_id"`
	DiseaseID int64  `json:"disease_id"`
	StartDate string `json:"start_date" validate:"required, datetime=2006-01-02"` // have lay out for date
	Status    string `json:"status" validate:"required,oneof=ongoing completed paused cancelled"`
	Notes     string `json:"notes"`
}

type Treatment struct {
	ID        int64  `json:"id"`
	PetName   string `json:"pet_name"`
	Disease   string `json:"disease"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status"`
}

type CreateTreatmentResponse struct {
	TreatmentID int64  `json:"treatment_id"`
	Disease     string `json:"disease"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Status      string `json:"status"`
}

type CreateTreatmentPhaseRequest struct {
	TreatmetnID int64  `json:"treatment_id"`
	PhaseName   string `json:"phase_name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	Status      string `json:"status" validate:"required,oneof=pending active completed"`
}

// Assign Medicines to Treatment Phases
type AssignMedicineRequest struct {
	MedicineID int64  `json:"medicine_id"`
	Dosage     string `json:"dosage"`
	Frequency  string `json:"frequency"`
	Duration   string `json:"duration"`
	Notes      string `json:"notes"`
}

type TreatmentPhase struct {
	ID          int64  `json:"id"`
	TreatmentID int64  `json:"treatment_id"`
	PhaseName   string `json:"phase_name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	StartDate   string `json:"start_date"`
	CreatedAt   string `json:"created_at"`
}

type PhaseMedicine struct {
	PhaseID      int64  `json:"phase_id"`
	MedicineName string `json:"medicine_name"`
	Dosage       string `json:"dosage"`
	Frequency    string `json:"frequency"`
	Duration     string `json:"duration"`
	Notes        string `json:"notes"`
	CreatedAt    string `json:"created_at"`
}

type UpdateTreatmentPhaseStatusRequest struct {
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
