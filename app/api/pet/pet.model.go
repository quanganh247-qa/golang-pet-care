package pet

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"

	"github.com/jackc/pgx/v5/pgtype"
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
	Gender          string  `json:"gender"`
	Healthnotes     string  `json:"healthnotes"`
	LastCheckedDate string  `json:"last_checked_date"`
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

// Add this new struct to your pet.model.go file
type PetLogWithPetInfo struct {
	PetID    int64  `json:"pet_id"`
	PetName  string `json:"pet_name"`
	PetType  string `json:"pet_type"`
	PetBreed string `json:"pet_breed"`
	LogID    int64  `json:"log_id"`
	DateTime string `json:"date_time"`
	Title    string `json:"title"`
	Notes    string `json:"notes"`
}

type UpdatePetLogRequeststruct struct {
	DateTime string `json:"date_time"`
	Title    string `json:"title"`
	Notes    string `json:"notes"`
}

type updatePetRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Breed string `json:"breed"`
	// Age             int     `json:"age"`
	Weight          float64 `json:"weight"`
	Gender          string  `json:"gender"`
	Healthnotes     string  `json:"healthnotes"`
	MicrochipNumber string  `json:"microchip_number"`
	BOD             string  `json:"birth_date"`
}

type updatePetAvatarRequest struct {
	DataImage     []byte `json:"-"`
	OriginalImage string `json:"original_name"`
}

type PetProfileSummary struct {
	Summary string `json:"summary"`
}

// Weight tracking models
// Weight record as stored in database
type PetWeightRecord struct {
	ID         int64            `json:"id"`
	PetID      int64            `json:"pet_id"`
	WeightKg   float64          `json:"weight_kg"`
	RecordedAt pgtype.Timestamp `json:"recorded_at"`
	Notes      pgtype.Text      `json:"notes"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
}

// Request to add a new weight record
type AddWeightRecordRequest struct {
	PetID    int64   `json:"pet_id"`
	WeightKg float64 `json:"weight_kg"`
	WeightLb float64 `json:"weight_lb,omitempty"`
	Notes    string  `json:"notes,omitempty"`
	UnitType string  `json:"unit_type"` // "kg" or "lb"
}

// Response for weight record
type WeightRecordResponse struct {
	ID         int64     `json:"id"`
	PetID      int64     `json:"pet_id"`
	WeightKg   float64   `json:"weight_kg"`
	WeightLb   float64   `json:"weight_lb"`
	RecordedAt time.Time `json:"recorded_at"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// Response for weight history
type PetWeightHistoryResponse struct {
	PetID           int64                  `json:"pet_id"`
	PetName         string                 `json:"pet_name"`
	CurrentWeight   WeightRecordResponse   `json:"current_weight"`
	WeightHistory   []WeightRecordResponse `json:"weight_history"`
	TotalRecords    int64                  `json:"total_records"`
	DefaultUnitType string                 `json:"default_unit_type"` // "kg" or "lb"
}

// Constants for unit conversion
const (
	KgToLbRatio = 2.20462
	LbToKgRatio = 0.453592
)

// Helper function to convert kg to lb
func KgToLb(kg float64) float64 {
	return kg * KgToLbRatio
}

// Helper function to convert lb to kg
func LbToKg(lb float64) float64 {
	return lb * LbToKgRatio
}
