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

// ItemType represents the type of item (test or vaccine)
type ItemType string

const (
	TypeTest    ItemType = "test"
	TypeVaccine ItemType = "vaccine"
)

// Test (now represents both tests and vaccines)
type TestItem struct {
	ID             int32    `json:"id"`
	TestID         string   `json:"test_id"`
	CategoryID     string   `json:"category_id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Price          float64  `json:"price"`
	TurnaroundTime string   `json:"turnaround_time"`
	Type           ItemType `json:"type"` // "test" or "vaccine"
	IsActive       bool     `json:"is_active"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type TestOrderRequest struct {
	AppointmentID int   `json:"appointment_id" binding:"required"`
	ItemIDs       []int `json:"item_ids" binding:"required"` // Renamed from TestIDs to ItemIDs
	// ItemType      ItemType `json:"item_type" binding:"required"` // Type of items being ordered
	Notes string `json:"notes"`
}

type TestResponse struct {
	OrderID      int              `json:"order_id"`
	TotalAmount  float64          `json:"total_amount"`
	OrderDate    pgtype.Timestamp `json:"order_date"`
	Status       string           `json:"status"`
	OrderedItems []OrderedItem    `json:"ordered_items"` // Renamed from OrderedTests
}

type OrderedItem struct {
	ItemID   int      `json:"item_id"`
	ItemName string   `json:"item_name"`
	ItemType ItemType `json:"item_type"`
	Price    float64  `json:"price"`
	Status   string   `json:"status"`
}

type AppointmentWithOrders struct {
	AppointmentID   int               `json:"appointment_id"`
	AppointmentDate string            `json:"appointment_date"`
	PetName         string            `json:"pet_name"`
	Species         string            `json:"species"`
	Orders          []TestOrderDetail `json:"orders"`
}

type TestOrderDetail struct {
	OrderID     int           `json:"order_id"`
	TotalAmount float64       `json:"total_amount"`
	Status      string        `json:"status"`
	OrderDate   string        `json:"order_date"`
	Items       []OrderedItem `json:"items"` // Renamed from Tests
}

type UpdateTestRequest struct {
	TestID         string   `json:"test_id" binding:"required"`
	Name           string   `json:"name" binding:"required"`
	Description    string   `json:"description"`
	Price          float64  `json:"price" binding:"required"`
	TurnaroundTime string   `json:"turnaround_time"`
	Type           ItemType `json:"type" binding:"required"`
	Dosage         string   `json:"dosage,omitempty"`
	Schedule       string   `json:"schedule,omitempty"`
}

// TestCategoryResponse represents a category of tests/vaccines with its associated items
type TestCategoryResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Icon        string     `json:"icon,omitempty"` // Frontend will handle the actual icon component
	Description string     `json:"description"`
	Items       []TestItem `json:"items"` // Renamed from Tests
}

type OrderedItemDetail struct {
	OrderedItemID int      `json:"ordered_item_id"`
	ItemID        string   `json:"item_id"`
	ItemName      string   `json:"item_name"`
	ItemType      ItemType `json:"item_type"`
	CategoryID    string   `json:"category_id"`
	CategoryName  string   `json:"category_name"`
	Price         float64  `json:"price"`
	Status        string   `json:"status"`
	OrderedDate   string   `json:"ordered_date"`
	Notes         string   `json:"notes"`
}

type TestByAppointment struct {
	TestID         string `json:"test_id"`
	TestName       string `json:"test_name"`
	ExpirationDate string `json:"expiration_date"`
	BatchNumber    string `json:"batch_number"`
}

// For backward compatibility
type Test = TestItem
type OrderedTest = OrderedItem
type OrderedTestDetail = OrderedItemDetail

type TestCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
