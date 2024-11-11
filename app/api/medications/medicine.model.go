package medications

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type MedicineApi struct {
	controller MedicineControllerInterface
}

type MedicineController struct {
	service MedicineServiceInterface
}

type MedicineService struct {
	storeDB db.Store
}

type createMedicineRequest struct {
	PetID        int64  `json:"pet_id"`
	MedicineName string `json:"medicine_name" validate:"required"`
	Dosage       string `json:"dosage" validate:"required"`
	Frequency    string `json:"frequency"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Notes        string `json:"notes"`
}

type createMedicineResponse struct {
	MedicineName string `json:"medicine_name"`
	Dosage       string `json:"dosage"`
	Frequency    string `json:"frequency"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Notes        string `json:"notes"`
}
