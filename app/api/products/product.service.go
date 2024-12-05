package products

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ProductServiceInterface interface {
	GetProducts(c *gin.Context, pagination *util.Pagination) ([]ProductResponse, error)
}

// get all products
func (s *ProductService) GetProducts(c *gin.Context, pagination *util.Pagination) ([]ProductResponse, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	products, err := s.storeDB.GetAllProducts(c, db.GetAllProductsParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	var productResponse []ProductResponse
	for _, product := range products {
		productResponse = append(productResponse, ProductResponse{
			ProductID:     product.ProductID,
			Name:          product.Name,
			Price:         product.Price,
			Stock:         product.StockQuantity.Int32,
			Category:      product.Category.String,
			DataImage:     product.DataImage,
			OriginalImage: product.OriginalImage.String,
		})
	}

	return productResponse, nil
}
