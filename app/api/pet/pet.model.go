package pet

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

type PetApi struct {
	controller PetControllerInterface
}

type PetController struct {
	service PetServiceInterface
}

type PetService struct {
	storeDB db.Store
	redis   *redis.ClientType
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
	ProfileImage    string  `json:"profileimage"`
	DataImage       []byte  `json:"-"`
	OriginalImage   string  `json:"original_name"`
}

type CreatePetResponse struct {
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

type PetLog struct {
	PetID    int64  `json:"pet_id"`
	LogID    int64  `json:"log_id"`
	DateTime string `json:"date_time"`
	Title    string `json:"title"`
	Notes    string `json:"notes"`
}

type updatePetRequest struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Breed           string  `json:"breed"`
	Age             int     `json:"age"`
	Weight          float64 `json:"weight"`
	Gender          string  `json:"gender"`
	Healthnotes     string  `json:"healthnotes"`
	BOD             string  `json:"birth_date"`
	MicrochipNumber string  `json:"microchip_number"`
}

type updatePetAvatarRequest struct {
	DataImage     []byte `json:"-"`
	OriginalImage string `json:"original_name"`
}
