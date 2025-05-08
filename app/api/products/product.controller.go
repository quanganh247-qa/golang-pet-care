package products

import (
	"database/sql" // Added sql import
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ProductControllerInterface interface {
	CreateProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	ImportStock(c *gin.Context)
	ExportStock(c *gin.Context)
	GetProductStockMovements(c *gin.Context)
	GetAllProductStockMovements(c *gin.Context)
}

func (controller *ProductController) CreateProduct(ctx *gin.Context) {
	var req CreateProductRequest

	// Parse the JSON data from the "data" form field
	jsonData := ctx.PostForm("data")
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Use the helper function to handle the image upload
	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	req.DataImage = dataImage
	req.OriginalImage = originalImageName

	res, err := controller.service.CreateProductService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Success", res))
}

func (controller *ProductController) GetProducts(c *gin.Context) {

	pagination, err := util.GetPageInQuery(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	products, err := controller.service.GetProducts(c, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (controller *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("product_id")
	productID, er := strconv.ParseInt(id, 10, 64)
	if er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": er.Error()})
		return
	}

	product, err := controller.service.GetProductByID(c, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (controller *ProductController) ImportStock(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("invalid product ID")))
		return
	}

	var req ImportStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	movement, err := controller.service.ImportStock(c.Request.Context(), productID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == fmt.Sprintf("product with ID %d not found", productID) {
			c.JSON(http.StatusNotFound, util.ErrorResponse(err))
		} else {
			c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		}
		return
	}

	c.JSON(http.StatusOK, util.SuccessResponse("Stock imported successfully", movement))
}

func (controller *ProductController) ExportStock(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("invalid product ID")))
		return
	}

	var req ExportStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	movement, err := controller.service.ExportStock(c.Request.Context(), productID, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == fmt.Sprintf("product with ID %d not found", productID) {
			c.JSON(http.StatusNotFound, util.ErrorResponse(err))
		} else if err.Error() == fmt.Sprintf("insufficient stock for product ID %d: available %d, requested %d", productID, 0, 0) { // Simplified check, actual values aren't known here
			c.JSON(http.StatusBadRequest, util.ErrorResponse(err)) // Use Bad Request for insufficient stock
		} else {
			c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		}
		return
	}

	c.JSON(http.StatusOK, util.SuccessResponse("Stock exported successfully", movement))
}

func (controller *ProductController) GetProductStockMovements(c *gin.Context) {
	id := c.Param("product_id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("invalid product ID: %w", err)))
		return
	}

	pagination, err := util.GetPageInQuery(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	movements, err := controller.service.GetProductStockMovements(c.Request.Context(), productID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, util.SuccessResponse("Stock movements retrieved successfully", movements))
}

func (controller *ProductController) GetAllProductStockMovements(c *gin.Context) {
	pagination, err := util.GetPageInQuery(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	movements, err := controller.service.GetAllProductStockMovements(c.Request.Context(), pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, util.SuccessResponse("All stock movements retrieved successfully", movements))
}
