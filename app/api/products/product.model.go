package products

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

type ProductApi struct {
	controller ProductControllerInterface
}

type ProductController struct {
	service ProductServiceInterface
}

type ProductService struct {
	storeDB db.Store
	redis   *redis.ClientType
}

type ProductResponse struct {
	ProductID     int64   `json:"product_id"`
	Description   string  `json:"description"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Stock         int32   `json:"stock"`
	Category      string  `json:"category"`
	DataImage     []byte  `json:"data_image"`
	OriginalImage string  `json:"original_name"`
	IsAvailable   *bool   `json:"is_available"` // Availability status (optional, default true)

}

// CreateProductRequest represents the structure for creating a new product
type CreateProductRequest struct {
	Name          string  `json:"name" validate:"required"`              // Product name (required)
	Description   string  `json:"description,omitempty"`                 // Product description (optional)
	Price         float64 `json:"price" validate:"required"`             // Product price (required)
	StockQuantity int     `json:"stock_quantity,omitempty" default:"0"`  // Stock quantity (optional, default 0)
	Category      string  `json:"category,omitempty"`                    // Product category (optional)
	DataImage     []byte  `json:"data_image,omitempty"`                  // Binary image data (optional)
	OriginalImage string  `json:"original_image,omitempty"`              // Image file name or URL (optional)
	IsAvailable   *bool   `json:"is_available,omitempty" default:"true"` // Availability status (optional, default true)
}

// ImportStockRequest represents the structure for importing stock
type ImportStockRequest struct {
	Quantity  int     `json:"quantity" validate:"required,gt=0"`   // Quantity to import (must be positive)
	Reason    string  `json:"reason,omitempty"`                    // Reason for the import (optional)
	UnitPrice float64 `json:"unit_price" validate:"required,gt=0"` // Unit price at the time of import (required, positive)
}

// ExportStockRequest represents the structure for exporting stock
type ExportStockRequest struct {
	Quantity  int     `json:"quantity" validate:"required,gt=0"`    // Quantity to export (must be positive)
	Reason    string  `json:"reason,omitempty"`                     // Reason for the export (optional)
	UnitPrice float64 `json:"unit_price" validate:"omitempty,gt=0"` // Unit price at the time of export (optional, positive)
}

// ProductStockMovementResponse represents a single stock movement record
type ProductStockMovementResponse struct {
	ID           int64     `json:"id"`
	ProductID    int64     `json:"product_id"`
	MovementType string    `json:"movement_type"` // "import" or "export"
	Quantity     int64     `json:"quantity"`
	Reason       string    `json:"reason,omitempty"`
	MovementDate time.Time `json:"movement_date"`
	CurrentStock int32     `json:"current_stock,omitempty"` // Only relevant for the response of Import/Export, not GetHistory
	Price        float64   `json:"price"`                   // Add price field
}
