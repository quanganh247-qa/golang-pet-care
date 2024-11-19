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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 67140c6 (updated create pet)
=======
>>>>>>> 67140c6 (updated create pet)
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Breed           string  `json:"breed"`
	Age             int16   `json:"age"`
	Weight          float64 `json:"weight"`
	Gender          string  `json:"gender"`
	Healthnotes     string  `json:"healthnotes"`
	BOD             string  `json:"birth_date"`
	MicrochipNumber string  `json:"microchip_number"`
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 7a9ad08 (updated pet api)
	ProfileImage    string  `json:"profileimage"`
	DataImage       []byte  `json:"-"`
	OriginalImage   string  `json:"original_name"`
}

type CreatePetResponse struct {
<<<<<<< HEAD
	Petid           int64   `json:"petid"`
	Username        string  `json:"username"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Breed           string  `json:"breed"`
	Age             int16   `json:"age"`
	BOD             string  `json:"birth_date"`
	Weight          float64 `json:"weight"`
	DataImage       []byte  `json:"data_image"`
<<<<<<< HEAD
	OriginalImage   string  `json:"original_name"`
	MicrochipNumber string  `json:"microchip_number"`
=======
=======
>>>>>>> 9d28896 (image pet)
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Breed            string  `json:"breed"`
	Age              int16   `json:"age"`
	Weight           float64 `json:"weight"`
	Gender           string  `json:"gender"`
	Healthnotes      string  `json:"healthnotes"`
	BOD              string  `json:"BOD"`
	MicrophoneNumber string  `json:"microphone"`
	DataImage        []byte  `json:"-"`
	OriginalImage    string  `json:"original_name"`
}

type createPetResponse struct {
	Petid         int64   `json:"petid"`
	Username      string  `json:"username"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Breed         string  `json:"breed"`
	Age           int16   `json:"age"`
	Weight        float64 `json:"weight"`
	DataImage     []byte  `json:"-"`
	OriginalImage string  `json:"original_name"`
<<<<<<< HEAD
>>>>>>> 9d28896 (image pet)
=======
=======
	ProfileImage    string  `json:"profileimage"`
>>>>>>> 7a9ad08 (updated pet api)
	DataImage       []byte  `json:"-"`
	OriginalImage   string  `json:"original_name"`
}

type createPetResponse struct {
=======
>>>>>>> 98e9e45 (ratelimit and recovery function)
=======
	DataImage       []byte  `json:"-"`
	OriginalImage   string  `json:"original_name"`
}

type createPetResponse struct {
>>>>>>> 67140c6 (updated create pet)
	Petid           int64   `json:"petid"`
	Username        string  `json:"username"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Breed           string  `json:"breed"`
	Age             int16   `json:"age"`
	BOD             string  `json:"birth_date"`
	Weight          float64 `json:"weight"`
<<<<<<< HEAD
	DataImage       []byte  `json:"data_image"`
	OriginalImage   string  `json:"original_name"`
	MicrochipNumber string  `json:"microchip_number"`
>>>>>>> 67140c6 (updated create pet)
=======
>>>>>>> 9d28896 (image pet)
=======
	DataImage       []byte  `json:"-"`
=======
>>>>>>> 0637caf (upadted get lít)
	OriginalImage   string  `json:"original_name"`
	MicrochipNumber string  `json:"microchip_number"`
>>>>>>> 67140c6 (updated create pet)
}

type listPetsRequest struct {
	Type   string  `json:"type"`
	Breed  string  `json:"breed"`
	Age    int     `json:"age"`
	Weight float64 `json:"weight"`
<<<<<<< HEAD
<<<<<<< HEAD
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

type PetProfileSummary struct {
	Summary string `json:"summary"`
=======
>>>>>>> c73e2dc (pagination function)
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

type PetProfileSummary struct {
	Summary string `json:"summary"`
=======
>>>>>>> c73e2dc (pagination function)
}

type PetLog struct {
	PetID    int64  `json:"pet_id"`
	LogID    int64  `json:"log_id"`
	DateTime string `json:"date_time"`
	Title    string `json:"title"`
	Notes    string `json:"notes"`
}
