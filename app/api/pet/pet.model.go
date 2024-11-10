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
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	Breed        string  `json:"breed"`
	Age          int16   `json:"age"`
	Weight       float64 `json:"weight"`
	Gender       string  `json:"gender"`
	Healthnotes  string  `json:"healthnotes"`
	Profileimage string  `json:"profileimage"`
}

type createPetResponse struct {
	Petid    int64   `json:"petid"`
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Breed    string  `json:"breed"`
	Age      int16   `json:"age"`
	Weight   float64 `json:"weight"`
}

type listPetsRequest struct {
	Type   string  `json:"type"`
	Breed  string  `json:"breed"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
}
