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
