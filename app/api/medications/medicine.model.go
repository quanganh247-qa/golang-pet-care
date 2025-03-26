package medications

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
)

type MedicineApi struct {
	controller MedicineControllerInterface
}

type MedicineController struct {
	service MedicineServiceInterface
}

type MedicineService struct {
	storeDB db.Store
	es      *elasticsearch.ESService
}

type createMedicineRequest struct {
	MedicineName   string `json:"medicine_name" validate:"required"`
	Dosage         string `json:"dosage" validate:"required"`
	Frequency      string `json:"frequency"`
	Duration       string `json:"duration"`
	SideEffects    string `json:"side_effects"`
	Quantity       int64  `json:"quantity"`
	ExpirationDate string `json:"expiration_date"`
	Description    string `json:"description"`
	Usage          string `json:"usage"`
}

type createMedicineResponse struct {
	MedicineName   string `json:"medicine_name"`
	Description    string `json:"description"`
	Usage          string `json:"usage"`
	Dosage         string `json:"dosage"`
	Frequency      string `json:"frequency"`
	Duration       string `json:"duration"`
	SideEffects    string `json:"side_effects"`
	ExpirationDate string `json:"expiration_date"`
	Quantity       int64  `json:"quantity"`
}
