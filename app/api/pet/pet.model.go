package pet

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type PetApi struct {
	controller PetControllerInterface
}

type PetController struct {
	service PetServiceInterface
}

type PetService struct {
	storeDB db.Store
}

type createPetRequest struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Breed           string  `json:"breed"`
	Age             int16   `json:"age"`
	Weight          float64 `json:"weight"`
	Gender          string  `json:"gender"`
	Healthnotes     string  `json:"healthnotes"`
	BOD             string  `json:"birth_date"`
	MicrochipNumber string  `json:"microchip_number"`
	DataImage       []byte  `json:"-"`
	OriginalImage   string  `json:"original_name"`
}

type createPetResponse struct {
	Petid           int64   `json:"petid"`
	Username        string  `json:"username"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Breed           string  `json:"breed"`
	Age             int16   `json:"age"`
	BOD             string  `json:"birth_date"`
	Weight          float64 `json:"weight"`
	DataImage       []byte  `json:"data_image"`
	OriginalImage   string  `json:"original_name"`
	MicrochipNumber string  `json:"microchip_number"`
}

type listPetsRequest struct {
	Type   string  `json:"type"`
	Breed  string  `json:"breed"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
}
