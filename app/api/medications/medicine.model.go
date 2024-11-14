package medications

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
=======
>>>>>>> 79a3bcc (medicine api)
)
=======
// type MedicineApi struct {
// 	controller MedicineControllerInterface
// }
>>>>>>> 6c35562 (dicease and treatment plan)
=======
import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
>>>>>>> a415f25 (new data)
=======
import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
)
>>>>>>> e859654 (Elastic search)
=======
import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)
>>>>>>> 79a3bcc (medicine api)
=======
// type MedicineApi struct {
// 	controller MedicineControllerInterface
// }
>>>>>>> 6c35562 (dicease and treatment plan)

// type MedicineController struct {
// 	service MedicineServiceInterface
// }

// type MedicineService struct {
// 	storeDB db.Store
// }

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
type MedicineService struct {
	storeDB db.Store
<<<<<<< HEAD
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
=======
=======
type MedicineService struct {
	storeDB db.Store
>>>>>>> 79a3bcc (medicine api)
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
<<<<<<< HEAD
>>>>>>> 79a3bcc (medicine api)
}
=======
// type createMedicineRequest struct {
// 	PetID        int64  `json:"pet_id"`
// 	MedicineName string `json:"medicine_name" validate:"required"`
// 	Dosage       string `json:"dosage" validate:"required"`
// 	Frequency    string `json:"frequency"`
// 	StartDate    string `json:"start_date"`
// 	EndDate      string `json:"end_date"`
// 	Notes        string `json:"notes"`
// }

=======
// type createMedicineRequest struct {
// 	PetID        int64  `json:"pet_id"`
// 	MedicineName string `json:"medicine_name" validate:"required"`
// 	Dosage       string `json:"dosage" validate:"required"`
// 	Frequency    string `json:"frequency"`
// 	StartDate    string `json:"start_date"`
// 	EndDate      string `json:"end_date"`
// 	Notes        string `json:"notes"`
// }

>>>>>>> 6c35562 (dicease and treatment plan)
// type createMedicineResponse struct {
// 	MedicineName string `json:"medicine_name"`
// 	Dosage       string `json:"dosage"`
// 	Frequency    string `json:"frequency"`
// 	StartDate    string `json:"start_date"`
// 	EndDate      string `json:"end_date"`
// 	Notes        string `json:"notes"`
// }
<<<<<<< HEAD
>>>>>>> 6c35562 (dicease and treatment plan)
=======
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
>>>>>>> a415f25 (new data)
=======
}
>>>>>>> 79a3bcc (medicine api)
=======
>>>>>>> 6c35562 (dicease and treatment plan)
