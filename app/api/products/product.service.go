package products

import (
	"fmt"

	"github.com/gin-gonic/gin"
<<<<<<< HEAD
	"github.com/jackc/pgx/v5/pgtype"
=======
>>>>>>> bd5945b (get list products)
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ProductServiceInterface interface {
<<<<<<< HEAD
	CreateProductService(c *gin.Context, req CreateProductRequest) (*ProductResponse, error)
	GetProducts(c *gin.Context, pagination *util.Pagination) ([]ProductResponse, error)
	GetProductByID(c *gin.Context, productID int64) (*ProductResponse, error)
=======
	GetProducts(c *gin.Context, pagination *util.Pagination) ([]ProductResponse, error)
>>>>>>> bd5945b (get list products)
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
<<<<<<< HEAD

// get product by id
func (s *ProductService) GetProductByID(c *gin.Context, productID int64) (*ProductResponse, error) {
	product, err := s.storeDB.GetProductByID(c, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	productResponse := ProductResponse{
		ProductID:     product.ProductID,
		Name:          product.Name,
		Price:         product.Price,
		Stock:         product.StockQuantity.Int32,
		Category:      product.Category.String,
		DataImage:     product.DataImage,
		OriginalImage: product.OriginalImage.String,
	}

	return &productResponse, nil
}

// create product
func (s *ProductService) CreateProductService(c *gin.Context, req CreateProductRequest) (*ProductResponse, error) {
	// insert product
	product, err := s.storeDB.InsertProduct(c, db.InsertProductParams{
		Name:          req.Name,
		Description:   pgtype.Text{String: req.Description, Valid: true},
		Price:         req.Price,
		StockQuantity: pgtype.Int4{Int32: int32(req.StockQuantity), Valid: true},
		Category:      pgtype.Text{String: req.Category, Valid: true},
		DataImage:     req.DataImage,
		OriginalImage: pgtype.Text{String: req.OriginalImage, Valid: true},
		IsAvailable:   pgtype.Bool{Bool: *req.IsAvailable, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	productResponse := ProductResponse{
		ProductID:     product.ProductID,
		Name:          product.Name,
		Price:         product.Price,
		Stock:         product.StockQuantity.Int32,
		Category:      product.Category.String,
		DataImage:     product.DataImage,
		OriginalImage: product.OriginalImage.String,
	}

	return &productResponse, nil

}
=======
>>>>>>> bd5945b (get list products)
