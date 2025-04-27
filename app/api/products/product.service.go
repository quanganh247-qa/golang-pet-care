package products

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	redisCache "github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ProductServiceInterface interface {
	CreateProductService(c *gin.Context, req CreateProductRequest) (*ProductResponse, error)
	GetProducts(c *gin.Context, pagination *util.Pagination) ([]ProductResponse, error)
	GetProductByID(c *gin.Context, productID int64) (*ProductResponse, error)
	ImportStock(ctx context.Context, productID int64, req ImportStockRequest) (*ProductStockMovementResponse, error)
	ExportStock(ctx context.Context, productID int64, req ExportStockRequest) (*ProductStockMovementResponse, error)
	GetProductStockMovements(ctx context.Context, productID int64, pagination *util.Pagination) ([]ProductStockMovementResponse, error)
}

// get all products
func (s *ProductService) GetProducts(c *gin.Context, pagination *util.Pagination) ([]ProductResponse, error) {

	// If cache miss or error, fallback to database
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
			Description:   product.Description.String,
			Name:          product.Name,
			Price:         product.Price,
			Stock:         product.StockQuantity.Int32,
			Category:      product.Category.String,
			DataImage:     product.DataImage,
			OriginalImage: product.OriginalImage.String,
			IsAvailable:   &product.IsAvailable.Bool,
		})
	}

	return productResponse, nil
}

// get product by id
func (s *ProductService) GetProductByID(c *gin.Context, productID int64) (*ProductResponse, error) {
	// Try to get from cache first
	cachedProduct, err := redisCache.Client.ProductInfoLoadCache(productID)
	if err == nil {
		// Cache hit, convert and return
		return &ProductResponse{
			ProductID:     cachedProduct.ProductID,
			Name:          cachedProduct.Name,
			Price:         cachedProduct.Price,
			Stock:         cachedProduct.Stock,
			Category:      cachedProduct.Category,
			DataImage:     cachedProduct.DataImage,
			OriginalImage: cachedProduct.OriginalImage,
		}, nil
	}

	// If cache miss or error, fallback to database
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

	// Clear caches when a new product is created
	redisCache.Client.ClearProductInfoCache()

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

// ImportStock increases the stock quantity of a product and logs the movement.
func (s *ProductService) ImportStock(ctx context.Context, productID int64, req ImportStockRequest) (*ProductStockMovementResponse, error) {
	var movement db.ProductStockMovement
	var updatedProduct db.Product

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Get product for update to lock the row
		product, err := q.GetProductByIDForUpdate(ctx, productID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("product with ID %d not found", productID)
			}
			return fmt.Errorf("failed to get product for update: %w", err)
		}

		// Calculate total price
		totalPrice := req.UnitPrice * float64(req.Quantity)

		// Convert the calculated total price to pgtype.Numeric
		totalPriceNumeric := pgtype.Numeric{}
		err = totalPriceNumeric.Scan(fmt.Sprintf("%f", totalPrice))
		if err != nil {
			return fmt.Errorf("failed to convert total price to numeric: %w", err)
		}

		// Create stock movement record
		movementArg := db.CreateProductStockMovementParams{
			ProductID:    productID,
			MovementType: "import",
			Quantity:     int32(req.Quantity),
			Reason:       pgtype.Text{String: req.Reason, Valid: req.Reason != ""},
			Price:        totalPriceNumeric, // Use the calculated total price
		}
		movement, err = q.CreateProductStockMovement(ctx, movementArg)
		if err != nil {
			return fmt.Errorf("failed to create stock movement record: %w", err)
		}

		// Update product stock quantity
		newStock := product.StockQuantity.Int32 + int32(req.Quantity) // req.Quantity is int, needs cast to int32 for calculation
		updateArg := db.UpdateProductStockParams{
			ProductID:     productID,
			StockQuantity: pgtype.Int4{Int32: newStock, Valid: true},
		}
		updatedProduct, err = q.UpdateProductStock(ctx, updateArg)
		if err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Clear relevant cache entries
	redisCache.Client.RemoveProductInfoCache(productID)

	// Convert movement price back to float64 for response
	var movementPriceFloat float64
	if movement.Price.Valid {
		bf := new(big.Float)
		_, err := fmt.Sscan(movement.Price.Int.String(), bf)
		if err == nil {
			bf.SetInt(movement.Price.Int)
			bf.Quo(bf, big.NewFloat(0).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(movement.Price.Exp)), nil)))
			movementPriceFloat, _ = bf.Float64()
		}
	}

	return &ProductStockMovementResponse{
		ID:           movement.MovementID,
		ProductID:    movement.ProductID,
		MovementType: string(movement.MovementType),
		Quantity:     int64(movement.Quantity),
		Reason:       movement.Reason.String,
		MovementDate: movement.MovementDate,
		CurrentStock: updatedProduct.StockQuantity.Int32,
		Price:        movementPriceFloat, // Use the converted float64 price
	}, nil
}

// ExportStock decreases the stock quantity of a product and logs the movement.
func (s *ProductService) ExportStock(ctx context.Context, productID int64, req ExportStockRequest) (*ProductStockMovementResponse, error) {
	var movement db.ProductStockMovement
	var updatedProduct db.Product

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Get product for update to lock the row
		product, err := q.GetProductByIDForUpdate(ctx, productID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("product with ID %d not found", productID)
			}
			return fmt.Errorf("failed to get product for update: %w", err)
		}

		// Check for sufficient stock
		if product.StockQuantity.Int32 < int32(req.Quantity) {
			return fmt.Errorf("insufficient stock for product ID %d: available %d, requested %d", productID, product.StockQuantity.Int32, req.Quantity)
		}

		// Calculate total price if unit price is provided
		totalPriceNumeric := pgtype.Numeric{}
		if req.UnitPrice > 0 { // Only calculate and set if unit price is provided and positive
			totalPrice := req.UnitPrice * float64(req.Quantity)
			err = totalPriceNumeric.Scan(fmt.Sprintf("%f", totalPrice))
			if err != nil {
				return fmt.Errorf("failed to convert total export price to numeric: %w", err)
			}
		} else {
			totalPriceNumeric.Valid = false // Mark as NULL if unit price is not provided or not positive
		}

		// Create stock movement record
		movementArg := db.CreateProductStockMovementParams{
			ProductID:    productID,
			MovementType: "export",
			Quantity:     int32(req.Quantity),
			Reason:       pgtype.Text{String: req.Reason, Valid: req.Reason != ""},
			Price:        totalPriceNumeric, // Use the calculated total price (or NULL)
		}
		movement, err = q.CreateProductStockMovement(ctx, movementArg)
		if err != nil {
			return fmt.Errorf("failed to create stock movement record: %w", err)
		}

		// Update product stock quantity
		newStock := product.StockQuantity.Int32 - int32(req.Quantity) // req.Quantity is int, needs cast to int32 for calculation
		updateArg := db.UpdateProductStockParams{
			ProductID:     productID,
			StockQuantity: pgtype.Int4{Int32: newStock, Valid: true},
		}
		updatedProduct, err = q.UpdateProductStock(ctx, updateArg)
		if err != nil {
			return fmt.Errorf("failed to update product stock: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	redisCache.Client.RemoveProductInfoCache(productID)

	// Convert movement price back to float64 for response
	var movementPriceFloat float64
	if movement.Price.Valid {
		bf := new(big.Float)
		_, err := fmt.Sscan(movement.Price.Int.String(), bf)
		if err == nil {
			bf.SetInt(movement.Price.Int)
			bf.Quo(bf, big.NewFloat(0).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(movement.Price.Exp)), nil)))
			movementPriceFloat, _ = bf.Float64()
		}
	}

	return &ProductStockMovementResponse{
		ID:           movement.MovementID,
		ProductID:    movement.ProductID,
		MovementType: string(movement.MovementType),
		Quantity:     int64(movement.Quantity),
		Reason:       movement.Reason.String,
		MovementDate: movement.MovementDate,
		CurrentStock: updatedProduct.StockQuantity.Int32,
		Price:        movementPriceFloat,
	}, nil
}

// GetProductStockMovements retrieves the stock movement history for a specific product.
func (s *ProductService) GetProductStockMovements(ctx context.Context, productID int64, pagination *util.Pagination) ([]ProductStockMovementResponse, error) {
	arg := db.GetProductStockMovementsByProductIDParams{
		ProductID: productID,
		Limit:     int32(pagination.PageSize),
		Offset:    int32((pagination.Page - 1) * pagination.PageSize),
	}

	movements, err := s.storeDB.GetProductStockMovementsByProductID(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock movements: %w", err)
	}

	var response []ProductStockMovementResponse
	for _, m := range movements {
		// Convert movement price back to float64 for response
		var movementPriceFloat float64
		if m.Price.Valid {
			bf := new(big.Float)
			_, err := fmt.Sscan(m.Price.Int.String(), bf)
			if err == nil {
				bf.SetInt(m.Price.Int)
				bf.Quo(bf, big.NewFloat(0).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(m.Price.Exp)), nil)))
				movementPriceFloat, _ = bf.Float64()
			}
		}

		response = append(response, ProductStockMovementResponse{
			ID:           m.MovementID,
			ProductID:    m.ProductID,
			MovementType: string(m.MovementType),
			Quantity:     int64(m.Quantity),
			Reason:       m.Reason.String,
			MovementDate: m.MovementDate,
			Price:        movementPriceFloat, // Use the converted float64 price
		})
	}

	return response, nil
}
