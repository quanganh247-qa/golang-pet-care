package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentControllerInterface interface {
	GetToken(c *gin.Context)
	GetBanks(c *gin.Context)
	GenerateQRCode(c *gin.Context)
}

func (c *PaymentController) GetToken(ctx *gin.Context) {
	result, err := c.service.GetToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *PaymentController) GetBanks(ctx *gin.Context) {
	result, err := c.service.GetBanksService(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *PaymentController) GenerateQRCode(ctx *gin.Context) {
	var qrRequest QRRequest
	if err := ctx.ShouldBindJSON(&qrRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.GenerateQRService(ctx, qrRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
