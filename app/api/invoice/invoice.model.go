package invoice

import (
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type InvoiceService struct {
	storeDB db.Store
}

type InvoiceController struct {
	service InvoiceServiceInterface
}

type InvoiceServiceInterface interface {
	CreateInvoice(ctx *gin.Context, req CreateInvoiceRequest) (*InvoiceResponse, error)
	GetInvoiceByID(ctx *gin.Context, id int32) (*InvoiceWithItemsResponse, error)
	GetInvoiceByNumber(ctx *gin.Context, invoiceNumber string) (*InvoiceWithItemsResponse, error)
	ListInvoices(ctx *gin.Context, pagination *util.Pagination) (*ListInvoicesResponse, error)
	UpdateInvoiceStatus(ctx *gin.Context, id int32, status string) error
	DeleteInvoice(ctx *gin.Context, id int32) error

	// Invoice item methods
	AddInvoiceItem(ctx *gin.Context, invoiceID int32, item InvoiceItem) (*InvoiceItemResponse, error)
	GetInvoiceItems(ctx *gin.Context, invoiceID int32) ([]InvoiceItemResponse, error)
	UpdateInvoiceItem(ctx *gin.Context, itemID int32, req UpdateInvoiceItemRequest) (*InvoiceItemResponse, error)
	DeleteInvoiceItem(ctx *gin.Context, itemID int32) error

	CreateInvoiceFromTestOrder(ctx *gin.Context, testOrderID int32) (*InvoiceResponse, error)
}

type CreateInvoiceRequest struct {
	InvoiceNumber string        `json:"invoice_number" binding:"required"`
	Amount        float64       `json:"amount" binding:"required"`
	Date          time.Time     `json:"date"`
	DueDate       time.Time     `json:"due_date" binding:"required"`
	Status        string        `json:"status"`
	Description   string        `json:"description"`
	CustomerName  string        `json:"customer_name" binding:"required"`
	Items         []InvoiceItem `json:"items" binding:"required"`
}

type InvoiceItem struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity int32   `json:"quantity"`
}

type InvoiceResponse struct {
	ID            int32     `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	Amount        float64   `json:"amount"`
	Date          time.Time `json:"date"`
	DueDate       time.Time `json:"due_date"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	CustomerName  string    `json:"customer_name"`
	CreatedAt     time.Time `json:"created_at"`
}

type InvoiceItemResponse struct {
	ID       int32   `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

type InvoiceWithItemsResponse struct {
	ID            int32                 `json:"id"`
	InvoiceNumber string                `json:"invoice_number"`
	Amount        float64               `json:"amount"`
	Date          time.Time             `json:"date"`
	DueDate       time.Time             `json:"due_date"`
	Status        string                `json:"status"`
	Description   string                `json:"description"`
	CustomerName  string                `json:"customer_name"`
	CreatedAt     time.Time             `json:"created_at"`
	Items         []InvoiceItemResponse `json:"items"`
}

type ListInvoicesResponse struct {
	Invoices []InvoiceWithItemsResponse `json:"invoices"`
	Total    int64                      `json:"total"`
}

type UpdateInvoiceStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// Add to your invoice.model.go file:
type UpdateInvoiceItemRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}
