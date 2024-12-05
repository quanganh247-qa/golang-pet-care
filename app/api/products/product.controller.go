package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ProductControllerInterface interface {
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
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
