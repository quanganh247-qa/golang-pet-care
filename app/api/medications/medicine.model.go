package medications

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

type MedicineApi struct {
	controller MedicineControllerInterface
}

type MedicineController struct {
	service MedicineServiceInterface
}

type MedicineService struct {
	storeDB         db.Store
	ws              *websocket.WSClientManager
	taskDistributor worker.TaskDistributor
}

type createMedicineRequest struct {
	MedicineName   string  `json:"medicine_name" validate:"required"`
	Dosage         string  `json:"dosage" validate:"required"`
	Frequency      string  `json:"frequency"`
	Duration       string  `json:"duration"`
	SideEffects    string  `json:"side_effects"`
	Quantity       int64   `json:"quantity"`
	ExpirationDate string  `json:"expiration_date"`
	Description    string  `json:"description"`
	Usage          string  `json:"usage"`
	SupplierID     int64   `json:"supplier_id"`
	UnitPrice      float64 `json:"unit_price"`
	ReorderLevel   int64   `json:"reorder_level"`
}

type createMedicineResponse struct {
	MedicineName   string  `json:"medicine_name"`
	Description    string  `json:"description"`
	Usage          string  `json:"usage"`
	Dosage         string  `json:"dosage"`
	Frequency      string  `json:"frequency"`
	Duration       string  `json:"duration"`
	SideEffects    string  `json:"side_effects"`
	ExpirationDate string  `json:"expiration_date"`
	Quantity       int64   `json:"quantity"`
	SupplierID     int64   `json:"supplier_id"`
	UnitPrice      float64 `json:"unit_price"`
	ReorderLevel   int64   `json:"reorder_level"`
}

type Medication struct {
	ID             int64   `json:"id"`
	MedicineName   string  `json:"medicine_name"`
	Dosage         string  `json:"dosage"`
	Frequency      string  `json:"frequency"`
	Duration       string  `json:"duration"`
	SideEffects    string  `json:"side_effects"`
	Description    string  `json:"description"`
	Usage          string  `json:"usage"`
	ExpirationDate string  `json:"expiration_date"`
	Quantity       int64   `json:"quantity"`
	SupplierName   string  `json:"supplier_name"`
	UnitPrice      float64 `json:"unit_price"`
	ReorderLevel   int64   `json:"reorder_level"`
}

type ListMedicineResponse struct {
	Medications []Medication `json:"medications"`
	Total       int64        `json:"total"`
}

// Models for inventory management
type MedicineTransactionRequest struct {
	MedicineID      int64   `json:"medicine_id" validate:"required"`
	Quantity        int64   `json:"quantity" validate:"required"`
	TransactionType string  `json:"transaction_type" validate:"required,oneof=import export"`
	UnitPrice       float64 `json:"unit_price"`
	SupplierID      int64   `json:"supplier_id"`
	ExpirationDate  string  `json:"expiration_date"`
	Notes           string  `json:"notes"`
	PrescriptionID  int64   `json:"prescription_id"`
	AppointmentID   int64   `json:"appointment_id"`
}

type MedicineTransactionResponse struct {
	ID              int64   `json:"id"`
	MedicineID      int64   `json:"medicine_id"`
	MedicineName    string  `json:"medicine_name"`
	Quantity        int64   `json:"quantity"`
	TransactionType string  `json:"transaction_type"`
	UnitPrice       float64 `json:"unit_price"`
	TotalAmount     float64 `json:"total_amount"`
	TransactionDate string  `json:"transaction_date"`
	SupplierID      int64   `json:"supplier_id"`
	SupplierName    string  `json:"supplier_name"`
	ExpirationDate  string  `json:"expiration_date"`
	Notes           string  `json:"notes"`
	PrescriptionID  int64   `json:"prescription_id"`
	AppointmentID   int64   `json:"appointment_id"`
}

type MedicineSupplierRequest struct {
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	ContactName string `json:"contact_name"`
	Notes       string `json:"notes"`
}

type MedicineSupplierResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	ContactName string `json:"contact_name"`
	Notes       string `json:"notes"`
	CreatedAt   string `json:"created_at"`
}

type ExpiredMedicineNotification struct {
	MedicineID      int64  `json:"medicine_id"`
	MedicineName    string `json:"medicine_name"`
	ExpirationDate  string `json:"expiration_date"`
	DaysUntilExpiry int    `json:"days_until_expiry"`
	Quantity        int64  `json:"quantity"`
}

type LowStockNotification struct {
	MedicineID   int64  `json:"medicine_id"`
	MedicineName string `json:"medicine_name"`
	CurrentStock int64  `json:"current_stock"`
	ReorderLevel int64  `json:"reorder_level"`
}
