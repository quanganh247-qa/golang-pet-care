package vaccination

import "time"

type createVaccinationRequest struct {
	PetID            int64     `json:"pet_id"`
	VaccineName      string    `json:"vaccine_name"`
	DateAdministered time.Time `json:"date_administered"`
	NextDueDate      time.Time `json:"next_due_date,omitempty"`
	VaccineProvider  string    `json:"vaccine_provider,omitempty"`
	BatchNumber      string    `json:"batch_number,omitempty"`
	Notes            string    `json:"notes,omitempty"`
}

type updateVaccinationRequest struct {
	VaccinationID    int64     `json:"vaccination_id"`
	VaccineName      string    `json:"vaccine_name"`
	DateAdministered time.Time `json:"date_administered"`
	NextDueDate      time.Time `json:"next_due_date,omitempty"`
	VaccineProvider  string    `json:"vaccine_provider,omitempty"`
	BatchNumber      string    `json:"batch_number,omitempty"`
	Notes            string    `json:"notes,omitempty"`
}

type VaccinationResponse struct {
	VaccinationID    int64     `json:"vaccination_id"`
	PetID            int64     `json:"pet_id"`
	VaccineName      string    `json:"vaccine_name"`
	DateAdministered time.Time `json:"date_administered"`
	NextDueDate      time.Time `json:"next_due_date,omitempty"`
	VaccineProvider  string    `json:"vaccine_provider,omitempty"`
	BatchNumber      string    `json:"batch_number,omitempty"`
	Notes            string    `json:"notes,omitempty"`
	DaysRemaining    int       `json:"days_remaining,omitempty"` // How many days remain until the vaccine is due
}
