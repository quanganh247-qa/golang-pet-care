package products

import (
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
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Stock         int32   `json:"stock"`
	Category      string  `json:"category"`
	DataImage     []byte  `json:"data_image"`
	OriginalImage string  `json:"original_name"`
}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 1ec1fee (create product api)
=======
>>>>>>> 1ec1fee (create product api)

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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> bd5945b (get list products)
=======
>>>>>>> 1ec1fee (create product api)
=======
>>>>>>> bd5945b (get list products)
=======
>>>>>>> 1ec1fee (create product api)
