package products

import (
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 1ec1fee (create product api)
	"encoding/json"
	"net/http"
	"strconv"
<<<<<<< HEAD
=======
	"net/http"
>>>>>>> bd5945b (get list products)
=======
>>>>>>> 63e2c90 (get product by id)

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ProductControllerInterface interface {
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 1ec1fee (create product api)
	CreateProduct(c *gin.Context)
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
<<<<<<< HEAD
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
=======
	GetProducts(c *gin.Context)
>>>>>>> bd5945b (get list products)
=======
>>>>>>> 63e2c90 (get product by id)
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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 63e2c90 (get product by id)

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
<<<<<<< HEAD
=======
>>>>>>> bd5945b (get list products)
=======
>>>>>>> 63e2c90 (get product by id)
