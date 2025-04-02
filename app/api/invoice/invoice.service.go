package invoice

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func NewInvoiceService(store db.Store) *InvoiceService {
	return &InvoiceService{
		storeDB: store,
	}
}

func (s *InvoiceService) CreateInvoice(ctx *gin.Context, req CreateInvoiceRequest) (*InvoiceResponse, error) {
	var invoice db.Invoice
	var err error

	// Use transaction to ensure all operations succeed or fail together
	err = s.storeDB.ExecWithTransaction(ctx, func(queries *db.Queries) error {
		// Set default date if not provided
		date := req.Date
		if date.IsZero() {
			date = time.Now()
		}

		// Set default status if not provided
		status := req.Status
		if status == "" {
			status = "unpaid"
		}

		// Create invoice
		invoice, err = queries.CreateInvoice(ctx, db.CreateInvoiceParams{
			InvoiceNumber: req.InvoiceNumber,
			Amount:        req.Amount,
			Date:          pgtype.Date{Time: date, Valid: true},
			DueDate:       pgtype.Date{Time: req.DueDate, Valid: true},
			Status:        status,
			Description:   pgtype.Text{String: req.Description, Valid: req.Description != ""},
			CustomerName:  pgtype.Text{String: req.CustomerName, Valid: req.CustomerName != ""},
		})
		if err != nil {
			return fmt.Errorf("failed to create invoice: %w", err)
		}

		// Create invoice items
		for _, item := range req.Items {
			quantity := item.Quantity
			if quantity == 0 {
				quantity = 1
			}

			_, err = queries.CreateInvoiceItem(ctx, db.CreateInvoiceItemParams{
				InvoiceID: invoice.ID,
				Name:      item.Name,
				Price:     item.Price,
				Quantity:  quantity,
			})
			if err != nil {
				return fmt.Errorf("failed to create invoice item: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &InvoiceResponse{
		ID:            invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
		Amount:        invoice.Amount,
		Date:          invoice.Date.Time,
		DueDate:       invoice.DueDate.Time,
		Status:        invoice.Status,
		Description:   invoice.Description.String,
		CustomerName:  invoice.CustomerName.String,
		CreatedAt:     invoice.CreatedAt.Time,
	}, nil
}

func (s *InvoiceService) GetInvoiceByID(ctx *gin.Context, id int32) (*InvoiceWithItemsResponse, error) {
	// Get invoice
	invoice, err := s.storeDB.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Get invoice items
	items, err := s.storeDB.GetInvoiceItems(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}

	// Convert items to response format
	var itemResponses []InvoiceItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, InvoiceItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}

	return &InvoiceWithItemsResponse{
		ID:            invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
		Amount:        invoice.Amount,
		Date:          invoice.Date.Time,
		DueDate:       invoice.DueDate.Time,
		Status:        invoice.Status,
		Description:   invoice.Description.String,
		CustomerName:  invoice.CustomerName.String,
		CreatedAt:     invoice.CreatedAt.Time,
		Items:         itemResponses,
	}, nil
}

func (s *InvoiceService) GetInvoiceByNumber(ctx *gin.Context, invoiceNumber string) (*InvoiceWithItemsResponse, error) {
	// Get invoice
	invoice, err := s.storeDB.GetInvoiceByNumber(ctx, invoiceNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Get invoice items
	items, err := s.storeDB.GetInvoiceItems(ctx, invoice.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}

	// Convert items to response format
	var itemResponses []InvoiceItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, InvoiceItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}

	return &InvoiceWithItemsResponse{
		ID:            invoice.ID,
		InvoiceNumber: invoice.InvoiceNumber,
		Amount:        invoice.Amount,
		Date:          invoice.Date.Time,
		DueDate:       invoice.DueDate.Time,
		Status:        invoice.Status,
		Description:   invoice.Description.String,
		CustomerName:  invoice.CustomerName.String,
		CreatedAt:     invoice.CreatedAt.Time,
		Items:         itemResponses,
	}, nil
}

func (s *InvoiceService) ListInvoices(ctx *gin.Context, pagination *util.Pagination) (*ListInvoicesResponse, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	// Get invoices
	invoices, err := s.storeDB.ListInvoices(ctx, db.ListInvoicesParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}

	// Get total count
	total, err := s.storeDB.CountInvoices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count invoices: %w", err)
	}

	// Convert to response format with items
	var invoiceResponses []InvoiceWithItemsResponse
	for _, invoice := range invoices {
		// Get items for each invoice
		items, err := s.storeDB.GetInvoiceItems(ctx, invoice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get items for invoice %d: %w", invoice.ID, err)
		}

		// Convert items to response format
		var itemResponses []InvoiceItemResponse
		for _, item := range items {
			itemResponses = append(itemResponses, InvoiceItemResponse{
				ID:       item.ID,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			})
		}

		invoiceResponses = append(invoiceResponses, InvoiceWithItemsResponse{
			ID:            invoice.ID,
			InvoiceNumber: invoice.InvoiceNumber,
			Amount:        invoice.Amount,
			Date:          invoice.Date.Time,
			DueDate:       invoice.DueDate.Time,
			Status:        invoice.Status,
			Description:   invoice.Description.String,
			CustomerName:  invoice.CustomerName.String,
			CreatedAt:     invoice.CreatedAt.Time,
			Items:         itemResponses,
		})
	}

	return &ListInvoicesResponse{
		Invoices: invoiceResponses,
		Total:    total,
	}, nil
}

func (s *InvoiceService) UpdateInvoiceStatus(ctx *gin.Context, id int32, status string) error {
	err := s.storeDB.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("failed to update invoice status: %w", err)
	}
	return nil
}

func (s *InvoiceService) DeleteInvoice(ctx *gin.Context, id int32) error {
	err := s.storeDB.DeleteInvoice(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}
	return nil
}

// Add these functions to your InvoiceService struct

func (s *InvoiceService) AddInvoiceItem(ctx *gin.Context, invoiceID int32, item InvoiceItem) (*InvoiceItemResponse, error) {
	// First check if invoice exists
	_, err := s.storeDB.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	// Set default quantity if not provided
	quantity := item.Quantity
	if quantity == 0 {
		quantity = 1
	}

	// Create invoice item
	newItem, err := s.storeDB.CreateInvoiceItem(ctx, db.CreateInvoiceItemParams{
		InvoiceID: invoiceID,
		Name:      item.Name,
		Price:     item.Price,
		Quantity:  quantity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice item: %w", err)
	}

	// Update invoice total amount
	err = s.storeDB.UpdateInvoiceAmount(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice amount: %w", err)
	}

	return &InvoiceItemResponse{
		ID:       newItem.ID,
		Name:     newItem.Name,
		Price:    newItem.Price,
		Quantity: newItem.Quantity,
	}, nil
}

func (s *InvoiceService) GetInvoiceItems(ctx *gin.Context, invoiceID int32) ([]InvoiceItemResponse, error) {
	// Get invoice items
	items, err := s.storeDB.GetInvoiceItems(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice items: %w", err)
	}

	// Convert to response format
	var itemResponses []InvoiceItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, InvoiceItemResponse{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}

	return itemResponses, nil
}

func (s *InvoiceService) UpdateInvoiceItem(ctx *gin.Context, itemID int32, req UpdateInvoiceItemRequest) (*InvoiceItemResponse, error) {
	// Get the item to find its invoice ID
	item, err := s.storeDB.GetInvoiceItemByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("invoice item not found: %w", err)
	}

	// Prepare update parameters
	var name *string
	var price *float64
	var quantity *int32

	if req.Name != "" {
		name = &req.Name
	}

	if req.Price > 0 {
		price = &req.Price
	}

	if req.Quantity > 0 {
		quantity = &req.Quantity
	}

	// Update the item
	updatedItem, err := s.storeDB.UpdateInvoiceItem(ctx, db.UpdateInvoiceItemParams{
		ID:       itemID,
		Name:     *name,
		Price:    *price,
		Quantity: *quantity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice item: %w", err)
	}

	// Update invoice total amount
	err = s.storeDB.UpdateInvoiceAmount(ctx, item.InvoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice amount: %w", err)
	}

	return &InvoiceItemResponse{
		ID:       updatedItem.ID,
		Name:     updatedItem.Name,
		Price:    updatedItem.Price,
		Quantity: updatedItem.Quantity,
	}, nil
}

func (s *InvoiceService) DeleteInvoiceItem(ctx *gin.Context, itemID int32) error {
	// Get the item to find its invoice ID
	item, err := s.storeDB.GetInvoiceItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("invoice item not found: %w", err)
	}

	// Delete the item
	err = s.storeDB.DeleteInvoiceItem(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to delete invoice item: %w", err)
	}

	// Update invoice total amount
	err = s.storeDB.UpdateInvoiceAmount(ctx, item.InvoiceID)
	if err != nil {
		return fmt.Errorf("failed to update invoice amount: %w", err)
	}

	return nil
}

func (s *InvoiceService) CreateInvoiceFromTestOrder(ctx *gin.Context, testOrderID int32) (*InvoiceResponse, error) {
	// 1. Get test order details
	testOrder, err := s.storeDB.GetTestOrderByID(ctx, testOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test order: %w", err)
	}

	// 2. Get ordered tests
	orderedTests, err := s.storeDB.GetOrderedTestsByOrderID(ctx, pgtype.Int4{Int32: testOrderID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get ordered tests: %w", err)
	}

	// 3. Get appointment details for customer information
	appointment, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, int64(testOrder.AppointmentID.Int32))
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	// 4. Get pet and owner details
	pet, err := s.storeDB.GetPetByID(ctx, appointment.PetID.Int64)
	if err != nil {
		return nil, fmt.Errorf("failed to get pet details: %w", err)
	}

	// 5. Generate invoice number
	invoiceNumber := fmt.Sprintf("INV-TEST-%d-%s", testOrderID, time.Now().Format("20060102"))

	// 6. Calculate total amount from ordered tests
	var totalAmount float64
	var invoiceItems []InvoiceItem

	for _, test := range orderedTests {
		totalAmount += test.PriceAtOrder
		invoiceItems = append(invoiceItems, InvoiceItem{
			Name:     test.TestName.String,
			Price:    test.PriceAtOrder,
			Quantity: 1,
		})
	}

	// 7. Create invoice request
	dueDate := time.Now().AddDate(0, 0, 15) // Set due date to 15 days from now
	description := fmt.Sprintf("Invoice for test order #%d - Pet: %s", testOrderID, pet.Name)

	invoiceReq := CreateInvoiceRequest{
		InvoiceNumber: invoiceNumber,
		Amount:        totalAmount,
		Date:          time.Now(),
		DueDate:       dueDate,
		Status:        "unpaid",
		Description:   description,
		CustomerName:  appointment.Username.String,
		Items:         invoiceItems,
	}

	// 8. Create invoice
	invoice, err := s.CreateInvoice(ctx, invoiceReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	return invoice, nil
}
