package disease

import (
	"github.com/lib/pq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type DiceaseApi struct {
	controller DiceaseControllerInterface
}

type DiceaseController struct {
	service DiceaseServiceInterface
}

type DiceaseService struct {
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

// DiseaseMedicineInfo holds the disease and associated medicine information.
