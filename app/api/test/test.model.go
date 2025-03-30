package test

import (
	"github.com/jackc/pgx/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
)

type TestService struct {
	storeDB db.Store
	ws      *websocket.WSClientManager
}

type TestController struct {
	service TestServiceInterface
	ws      *websocket.WSClientManager
}

// type TestResponse struct {
// 	ID          int64      `json:"id"`
// 	PetID       int64      `json:"pet_id"`
// 	DoctorID    int64      `json:"doctor_id"`
// 	TestType    string     `json:"test_type"`
// 	Status      string     `json:"status"`
// 	CreatedAt   time.Time  `json:"created_at"`
// 	CompletedAt *time.Time `json:"completed_at,omitempty"`
// }

// type TestResult struct {
// 	Parameters map[string]interface{} `json:"parameters"`
// 	Notes      string                 `json:"notes"`
// 	Files      []string               `json:"files,omitempty"`
// }

// type StatusHistory struct {
// 	Status    string    `json:"status"`
// 	Timestamp time.Time `json:"timestamp"`
// 	UpdatedBy string    `json:"updated_by"`
// }

type Test struct {
	ID             int32   `json:"id"`
	TestID         string  `json:"test_id"`
	CategoryID     string  `json:"category_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	TurnaroundTime string  `json:"turnaround_time"`
	IsActive       bool    `json:"is_active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type TestOrderRequest struct {
	AppointmentID int    `json:"appointment_id" binding:"required"`
	TestIDs       []int  `json:"test_ids" binding:"required"`
	Notes         string `json:"notes"`
}

type TestResponse struct {
	OrderID      int              `json:"order_id"`
	TotalAmount  float64          `json:"total_amount"`
	OrderDate    pgtype.Timestamp `json:"order_date"`
	Status       string           `json:"status"`
	OrderedTests []OrderedTest    `json:"ordered_tests"`
}

type OrderedTest struct {
	TestID   int     `json:"test_id"`
	TestName string  `json:"test_name"`
	Price    float64 `json:"price"`
	Status   string  `json:"status"`
}

type UpdateTestRequest struct {
	TestID         string  `json:"test_id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Price          float64 `json:"price" binding:"required"`
	TurnaroundTime string  `json:"turnaround_time" binding:"required"`
}


// TestCategoryResponse represents a category of tests with its associated tests
type TestCategoryResponse struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Icon        string `json:"icon,omitempty"` // Frontend will handle the actual icon component
    Description string `json:"description"`
    Tests       []Test `json:"tests"`
}
