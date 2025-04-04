package medications

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicineServiceInterface interface {
	CreateMedicine(ctx *gin.Context, username string, req createMedicineRequest) (*createMedicineResponse, error)
	GetMedicineByID(ctx *gin.Context, medicineid int64) (*createMedicineResponse, error)
	ListMedicines(ctx *gin.Context, pagination *util.Pagination, petID int64) ([]Medication, error)
	UpdateMedicine(ctx *gin.Context, medicineid int64, req createMedicineRequest) error
	GetAllMedicines(ctx *gin.Context) ([]Medication, error)

	// Inventory management
	CreateMedicineTransaction(ctx *gin.Context, username string, req MedicineTransactionRequest) (*MedicineTransactionResponse, error)
	GetMedicineTransactions(ctx *gin.Context, pagination *util.Pagination) ([]MedicineTransactionResponse, error)
	GetMedicineTransactionsByMedicineID(ctx *gin.Context, medicineID int64, pagination *util.Pagination) ([]MedicineTransactionResponse, error)

	// Supplier management
	CreateSupplier(ctx *gin.Context, req MedicineSupplierRequest) (*MedicineSupplierResponse, error)
	GetSupplierByID(ctx *gin.Context, supplierID int64) (*MedicineSupplierResponse, error)
	GetAllSuppliers(ctx *gin.Context, pagination *util.Pagination) ([]MedicineSupplierResponse, error)
	UpdateSupplier(ctx *gin.Context, supplierID int64, req MedicineSupplierRequest) error

	// Expiration and stock alerts
	GetExpiringMedicines(ctx *gin.Context, days int) ([]ExpiredMedicineNotification, error)
	GetLowStockMedicines(ctx *gin.Context) ([]LowStockNotification, error)

	// WebSocket handler
	HandleWebSocket(ctx *gin.Context)
}

func (s *MedicineService) CreateMedicine(ctx *gin.Context, username string, req createMedicineRequest) (*createMedicineResponse, error) {
	var expirationDate time.Time
	var err error

	if req.ExpirationDate != "" {
		expirationDate, err = time.Parse("2006-01-02", req.ExpirationDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expiration date: %w", err)
		}
	}

	arg := db.CreateMedicineParams{
		Name:           req.MedicineName,
		Description:    pgtype.Text{String: req.Description, Valid: true},
		Usage:          pgtype.Text{String: req.Usage, Valid: true},
		Dosage:         pgtype.Text{String: req.Dosage, Valid: true},
		Frequency:      pgtype.Text{String: req.Frequency, Valid: true},
		Duration:       pgtype.Text{String: req.Duration, Valid: true},
		SideEffects:    pgtype.Text{String: req.SideEffects, Valid: true},
		Quantity:       pgtype.Int8{Int64: req.Quantity, Valid: true},
		ExpirationDate: pgtype.Date{Time: expirationDate, Valid: req.ExpirationDate != ""},
		UnitPrice:      pgtype.Float8{Float64: req.UnitPrice, Valid: req.UnitPrice > 0},
		ReorderLevel:   pgtype.Int8{Int64: req.ReorderLevel, Valid: req.ReorderLevel > 0},
	}

	medicine, err := s.storeDB.CreateMedicine(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create medicine: %w", err)
	}

	return &createMedicineResponse{
		MedicineName:   medicine.Name,
		Description:    medicine.Description.String,
		Usage:          medicine.Usage.String,
		Dosage:         medicine.Dosage.String,
		Frequency:      medicine.Frequency.String,
		Duration:       medicine.Duration.String,
		SideEffects:    medicine.SideEffects.String,
		Quantity:       medicine.Quantity.Int64,
		ExpirationDate: formatDate(medicine.ExpirationDate.Time),
		SupplierID:     req.SupplierID,
		UnitPrice:      medicine.UnitPrice.Float64,
		ReorderLevel:   medicine.ReorderLevel.Int64,
	}, nil
}

func (s *MedicineService) GetMedicineByID(ctx *gin.Context, medicineid int64) (*createMedicineResponse, error) {
	medicine, err := s.storeDB.GetMedicineByID(ctx, medicineid)
	if err != nil {
		return nil, fmt.Errorf("failed to get medicine: %w", err)
	}

	return &createMedicineResponse{
		MedicineName:   medicine.Name,
		Description:    medicine.Description.String,
		Usage:          medicine.Usage.String,
		Dosage:         medicine.Dosage.String,
		Frequency:      medicine.Frequency.String,
		Duration:       medicine.Duration.String,
		SideEffects:    medicine.SideEffects.String,
		Quantity:       medicine.Quantity.Int64,
		ExpirationDate: formatDate(medicine.ExpirationDate.Time),
		UnitPrice:      medicine.UnitPrice.Float64,
		ReorderLevel:   medicine.ReorderLevel.Int64,
	}, nil
}

func (s *MedicineService) ListMedicines(ctx *gin.Context, pagination *util.Pagination, petID int64) ([]Medication, error) {
	arg := db.ListMedicinesByPetParams{
		PetID:  pgtype.Int8{Int64: petID, Valid: true},
		Status: pgtype.Text{String: "active", Valid: true},
		Limit:  int32(pagination.PageSize),
		Offset: int32((pagination.Page - 1) * pagination.PageSize),
	}

	results, err := s.storeDB.ListMedicinesByPet(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to list medicines: %w", err)
	}

	var medicines []Medication
	for _, r := range results {
		medicines = append(medicines, Medication{
			ID:             r.ID,
			MedicineName:   r.Name,
			Dosage:         r.Dosage.String,
			Frequency:      r.Frequency.String,
			Duration:       r.Duration.String,
			SideEffects:    r.SideEffects.String,
			Description:    r.Description.String,
			Usage:          r.Usage.String,
			ExpirationDate: formatDate(r.ExpirationDate.Time),
			Quantity:       r.Quantity.Int64,
			UnitPrice:      r.UnitPrice.Float64,
			ReorderLevel:   r.ReorderLevel.Int64,
		})
	}

	return medicines, nil
}

func (s *MedicineService) UpdateMedicine(ctx *gin.Context, medicineid int64, req createMedicineRequest) error {
	var expirationDate time.Time
	var err error

	if req.ExpirationDate != "" {
		expirationDate, err = time.Parse("2006-01-02", req.ExpirationDate)
		if err != nil {
			return fmt.Errorf("failed to parse expiration date: %w", err)
		}
	}

	arg := db.UpdateMedicineParams{
		ID:             medicineid,
		Name:           req.MedicineName,
		Description:    pgtype.Text{String: req.Description, Valid: true},
		Usage:          pgtype.Text{String: req.Usage, Valid: true},
		Dosage:         pgtype.Text{String: req.Dosage, Valid: true},
		Frequency:      pgtype.Text{String: req.Frequency, Valid: true},
		Duration:       pgtype.Text{String: req.Duration, Valid: true},
		SideEffects:    pgtype.Text{String: req.SideEffects, Valid: true},
		Quantity:       pgtype.Int8{Int64: req.Quantity, Valid: true},
		ExpirationDate: pgtype.Date{Time: expirationDate, Valid: req.ExpirationDate != ""},
		UnitPrice:      pgtype.Float8{Float64: req.UnitPrice, Valid: req.UnitPrice > 0},
		ReorderLevel:   pgtype.Int8{Int64: req.ReorderLevel, Valid: req.ReorderLevel > 0},
	}

	err = s.storeDB.UpdateMedicine(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update medicine: %w", err)
	}

	return nil
}

func (s *MedicineService) GetAllMedicines(ctx *gin.Context) ([]Medication, error) {
	meds, err := s.storeDB.GetAllMedicines(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all medicines: %w", err)
	}

	var medsInfo []Medication
	for _, med := range meds {
		medsInfo = append(medsInfo, Medication{
			ID:             med.ID,
			MedicineName:   med.Name,
			Dosage:         med.Dosage.String,
			Frequency:      med.Frequency.String,
			Duration:       med.Duration.String,
			SideEffects:    med.SideEffects.String,
			Description:    med.Description.String,
			Usage:          med.Usage.String,
			ExpirationDate: formatDate(med.ExpirationDate.Time),
			Quantity:       med.Quantity.Int64,
			SupplierName:   med.SupplierName,
			UnitPrice:      med.UnitPrice.Float64,
			ReorderLevel:   med.ReorderLevel.Int64,
		})
	}

	return medsInfo, nil
}

// Inventory Management Implementation
func (s *MedicineService) CreateMedicineTransaction(ctx *gin.Context, username string, req MedicineTransactionRequest) (*MedicineTransactionResponse, error) {
	var err error
	var expirationDate time.Time

	if req.ExpirationDate != "" {
		expirationDate, err = time.Parse("2006-01-02", req.ExpirationDate)
		if err != nil {
			return nil, fmt.Errorf("failed to parse expiration date: %w", err)
		}
	}

	// Calculate total amount
	totalAmount := req.UnitPrice * float64(req.Quantity)

	// Begin a transaction
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Create the transaction record
		txParam := db.CreateMedicineTransactionParams{
			MedicineID:      req.MedicineID,
			Quantity:        req.Quantity,
			TransactionType: req.TransactionType,
			UnitPrice:       pgtype.Float8{Float64: req.UnitPrice, Valid: req.UnitPrice > 0},
			TotalAmount:     pgtype.Float8{Float64: totalAmount, Valid: true},
			SupplierID:      pgtype.Int8{Int64: req.SupplierID, Valid: req.SupplierID > 0},
			ExpirationDate:  pgtype.Date{Time: expirationDate, Valid: req.ExpirationDate != ""},
			Notes:           pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
			PrescriptionID:  pgtype.Int8{Int64: req.PrescriptionID, Valid: req.PrescriptionID > 0},
			AppointmentID:   pgtype.Int8{Int64: req.AppointmentID, Valid: req.AppointmentID > 0},
			CreatedBy:       pgtype.Text{String: username, Valid: true},
		}

		_, err := q.CreateMedicineTransaction(ctx, txParam)
		if err != nil {
			return fmt.Errorf("failed to create medicine transaction: %w", err)
		}

		// Update medicine quantity
		var qtyChange int64
		if req.TransactionType == "import" {
			qtyChange = req.Quantity
		} else {
			qtyChange = -req.Quantity
		}

		err = q.UpdateMedicineQuantity(ctx, db.UpdateMedicineQuantityParams{
			ID:       req.MedicineID,
			Quantity: pgtype.Int8{Int64: qtyChange, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to update medicine quantity: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Get the medicine and supplier details for the response
	medicine, err := s.storeDB.GetMedicineByID(ctx, req.MedicineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medicine details: %w", err)
	}

	var supplierName string
	if req.SupplierID > 0 {
		supplier, err := s.storeDB.GetSupplierByID(ctx, req.SupplierID)
		if err == nil {
			supplierName = supplier.Name
		}
	}

	response := &MedicineTransactionResponse{
		MedicineID:      req.MedicineID,
		MedicineName:    medicine.Name,
		Quantity:        req.Quantity,
		TransactionType: req.TransactionType,
		UnitPrice:       req.UnitPrice,
		TotalAmount:     totalAmount,
		TransactionDate: time.Now().Format(time.RFC3339),
		SupplierID:      req.SupplierID,
		SupplierName:    supplierName,
		ExpirationDate:  req.ExpirationDate,
		Notes:           req.Notes,
		PrescriptionID:  req.PrescriptionID,
		AppointmentID:   req.AppointmentID,
	}

	// Check for low stock after transaction and notify if needed
	if req.TransactionType == "export" {
		// Check if the medicine is now below reorder level
		s.checkAndNotifyLowStock(ctx, req.MedicineID)
	}

	return response, nil
}

func (s *MedicineService) GetMedicineTransactions(ctx *gin.Context, pagination *util.Pagination) ([]MedicineTransactionResponse, error) {
	arg := db.GetMedicineTransactionsParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32((pagination.Page - 1) * pagination.PageSize),
	}

	transactions, err := s.storeDB.GetMedicineTransactions(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get medicine transactions: %w", err)
	}

	var response []MedicineTransactionResponse
	for _, tx := range transactions {
		response = append(response, MedicineTransactionResponse{
			ID:              tx.ID,
			MedicineID:      tx.MedicineID,
			MedicineName:    tx.MedicineName,
			Quantity:        tx.Quantity,
			TransactionType: tx.TransactionType,
			UnitPrice:       tx.UnitPrice.Float64,
			TotalAmount:     tx.TotalAmount.Float64,
			TransactionDate: tx.TransactionDate.Time.Format(time.RFC3339),
			SupplierID:      tx.SupplierID.Int64,
			SupplierName:    tx.SupplierName,
			ExpirationDate:  formatDate(tx.ExpirationDate.Time),
			Notes:           tx.Notes.String,
			PrescriptionID:  tx.PrescriptionID.Int64,
			AppointmentID:   tx.AppointmentID.Int64,
		})
	}

	return response, nil
}

func (s *MedicineService) GetMedicineTransactionsByMedicineID(ctx *gin.Context, medicineID int64, pagination *util.Pagination) ([]MedicineTransactionResponse, error) {
	arg := db.GetMedicineTransactionsByMedicineIDParams{
		MedicineID: medicineID,
		Limit:      int32(pagination.PageSize),
		Offset:     int32((pagination.Page - 1) * pagination.PageSize),
	}

	transactions, err := s.storeDB.GetMedicineTransactionsByMedicineID(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get medicine transactions: %w", err)
	}

	var response []MedicineTransactionResponse
	for _, tx := range transactions {
		response = append(response, MedicineTransactionResponse{
			ID:              tx.ID,
			MedicineID:      tx.MedicineID,
			MedicineName:    tx.MedicineName,
			Quantity:        tx.Quantity,
			TransactionType: tx.TransactionType,
			UnitPrice:       tx.UnitPrice.Float64,
			TotalAmount:     tx.TotalAmount.Float64,
			TransactionDate: tx.TransactionDate.Time.Format(time.RFC3339),
			SupplierID:      tx.SupplierID.Int64,
			SupplierName:    tx.SupplierName,
			ExpirationDate:  formatDate(tx.ExpirationDate.Time),
			Notes:           tx.Notes.String,
			PrescriptionID:  tx.PrescriptionID.Int64,
			AppointmentID:   tx.AppointmentID.Int64,
		})
	}

	return response, nil
}

// Supplier Management Implementation

func (s *MedicineService) CreateSupplier(ctx *gin.Context, req MedicineSupplierRequest) (*MedicineSupplierResponse, error) {
	arg := db.CreateSupplierParams{
		Name:        req.Name,
		Email:       pgtype.Text{String: req.Email, Valid: req.Email != ""},
		Phone:       pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
		Address:     pgtype.Text{String: req.Address, Valid: req.Address != ""},
		ContactName: pgtype.Text{String: req.ContactName, Valid: req.ContactName != ""},
		Notes:       pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	}

	supplier, err := s.storeDB.CreateSupplier(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create supplier: %w", err)
	}

	return &MedicineSupplierResponse{
		ID:          supplier.ID,
		Name:        supplier.Name,
		Email:       supplier.Email.String,
		Phone:       supplier.Phone.String,
		Address:     supplier.Address.String,
		ContactName: supplier.ContactName.String,
		Notes:       supplier.Notes.String,
		CreatedAt:   supplier.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (s *MedicineService) GetSupplierByID(ctx *gin.Context, supplierID int64) (*MedicineSupplierResponse, error) {
	supplier, err := s.storeDB.GetSupplierByID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier: %w", err)
	}

	return &MedicineSupplierResponse{
		ID:          supplier.ID,
		Name:        supplier.Name,
		Email:       supplier.Email.String,
		Phone:       supplier.Phone.String,
		Address:     supplier.Address.String,
		ContactName: supplier.ContactName.String,
		Notes:       supplier.Notes.String,
		CreatedAt:   supplier.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (s *MedicineService) GetAllSuppliers(ctx *gin.Context, pagination *util.Pagination) ([]MedicineSupplierResponse, error) {
	arg := db.GetAllSuppliersParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32((pagination.Page - 1) * pagination.PageSize),
	}

	suppliers, err := s.storeDB.GetAllSuppliers(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get suppliers: %w", err)
	}

	var response []MedicineSupplierResponse
	for _, supplier := range suppliers {
		response = append(response, MedicineSupplierResponse{
			ID:          supplier.ID,
			Name:        supplier.Name,
			Email:       supplier.Email.String,
			Phone:       supplier.Phone.String,
			Address:     supplier.Address.String,
			ContactName: supplier.ContactName.String,
			Notes:       supplier.Notes.String,
			CreatedAt:   supplier.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *MedicineService) UpdateSupplier(ctx *gin.Context, supplierID int64, req MedicineSupplierRequest) error {
	arg := db.UpdateSupplierParams{
		ID:          supplierID,
		Name:        req.Name,
		Email:       pgtype.Text{String: req.Email, Valid: req.Email != ""},
		Phone:       pgtype.Text{String: req.Phone, Valid: req.Phone != ""},
		Address:     pgtype.Text{String: req.Address, Valid: req.Address != ""},
		ContactName: pgtype.Text{String: req.ContactName, Valid: req.ContactName != ""},
		Notes:       pgtype.Text{String: req.Notes, Valid: req.Notes != ""},
	}

	err := s.storeDB.UpdateSupplier(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update supplier: %w", err)
	}

	return nil
}

// Expiration and Stock Alerts Implementation
func (s *MedicineService) GetExpiringMedicines(ctx *gin.Context, days int) ([]ExpiredMedicineNotification, error) {

	expiryDate := time.Now().AddDate(0, 0, days)
	expiryDate, err := time.Parse("2006-01-02", expiryDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	// Create pgtype.Date properly
	daysDate := pgtype.Date{
		Time:  expiryDate,
		Valid: true,
	}

	expiringMedicines, err := s.storeDB.GetExpiringMedicines(ctx, daysDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring medicines: %w", err)
	}

	var notifications []ExpiredMedicineNotification
	for _, medicine := range expiringMedicines {
		notifications = append(notifications, ExpiredMedicineNotification{
			MedicineID:      medicine.ID,
			MedicineName:    medicine.Name,
			ExpirationDate:  formatDate(medicine.ExpirationDate.Time),
			DaysUntilExpiry: int(medicine.DaysUntilExpiry),
			Quantity:        medicine.Quantity.Int64,
		})
	}

	return notifications, nil
}

func (s *MedicineService) GetLowStockMedicines(ctx *gin.Context) ([]LowStockNotification, error) {
	lowStockMedicines, err := s.storeDB.GetLowStockMedicines(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock medicines: %w", err)
	}

	var notifications []LowStockNotification
	for _, medicine := range lowStockMedicines {
		notifications = append(notifications, LowStockNotification{
			MedicineID:   medicine.ID,
			MedicineName: medicine.Name,
			CurrentStock: medicine.CurrentStock.Int64,
			ReorderLevel: medicine.ReorderLevel.Int64,
		})
	}

	return notifications, nil
}

// WebSocket Handler

func (s *MedicineService) HandleWebSocket(ctx *gin.Context) {
	s.ws.HandleWebSocket(ctx)
}

// Helper functions

func (s *MedicineService) checkAndNotifyLowStock(ctx *gin.Context, medicineID int64) {
	// Get medicine details
	medicine, err := s.storeDB.GetMedicineByID(ctx, medicineID)
	if err != nil {
		// Log error but don't interrupt the transaction
		fmt.Printf("Error checking stock level: %v\n", err)
		return
	}

	// Check if the medicine is now below reorder level
	if medicine.Quantity.Int64 < medicine.ReorderLevel.Int64 {
		notification := LowStockNotification{
			MedicineID:   medicineID,
			MedicineName: medicine.Name,
			CurrentStock: medicine.Quantity.Int64,
			ReorderLevel: medicine.ReorderLevel.Int64,
		}

		// Send WebSocket notification
		s.sendLowStockNotification(notification)
	}
}

func (s *MedicineService) sendLowStockNotification(notification LowStockNotification) {
	// Send notification via WebSocket
	wsMessage := websocket.WebSocketMessage{
		Type: "low_stock_alert",
		Data: notification,
	}

	s.ws.BroadcastToAll(wsMessage)
}

func (s *MedicineService) CheckExpiringMedicines() {
	// This would be scheduled to run daily
	// Get medicines expiring in the next 30 days
	ctx := context.Background()

	// Create a new gin context from the background context
	ginCtx := &gin.Context{}
	ginCtx.Request = &http.Request{}
	ginCtx.Request = ginCtx.Request.WithContext(ctx)

	expiringMedicines, err := s.GetExpiringMedicines(ginCtx, 30)
	if err != nil {
		fmt.Printf("Error checking expiring medicines: %v\n", err)
		return
	}

	// Send notification for each expiring medicine
	for _, medicine := range expiringMedicines {
		wsMessage := websocket.WebSocketMessage{
			Type: "expiration_alert",
			Data: medicine,
		}

		s.ws.BroadcastToAll(wsMessage)
	}
}

// Utility function to format date
func formatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02")
}
