package models

// AITreatmentPlan represents an AI-generated treatment plan for a pet
type AITreatmentPlan struct {
	ID              int64            `json:"id"`
	PatientID       int64            `json:"patient_id"`
	Diagnosis       string           `json:"diagnosis"`
	Confidence      float64          `json:"confidence"`
	TreatmentPhases []TreatmentPhase `json:"treatment_phases"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
}

// TreatmentPhase represents a phase in the treatment plan
type TreatmentPhase struct {
	ID          int64           `json:"id"`
	PlanID      int64           `json:"plan_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Duration    int             `json:"duration"` // Duration in days
	Medicines   []PhaseMedicine `json:"medicines"`
	Order       int             `json:"order"`
}

// PhaseMedicine represents a medicine prescribed in a treatment phase
type PhaseMedicine struct {
	ID         int64  `json:"id"`
	PhaseID    int64  `json:"phase_id"`
	MedicineID int64  `json:"medicine_id"`
	Dosage     string `json:"dosage"`
	Frequency  string `json:"frequency"`
	Duration   int    `json:"duration"` // Duration in days
	Notes      string `json:"notes"`
}
