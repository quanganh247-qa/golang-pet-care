package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Add this to the PaymentControllerInterface
type PaymentControllerInterface interface {
	GetToken(c *gin.Context)
	GetBanks(c *gin.Context)
	GenerateQRCode(c *gin.Context)

	CreatePayPalOrder(c *gin.Context)
	CapturePayPalOrder(c *gin.Context)
	GetOrderDetails(c *gin.Context)
	UpdateOrder(c *gin.Context)
	TrackOrder(c *gin.Context)

	CreatePayOSLink(c *gin.Context)

	GetRevenueLastSevenDays(ctx *gin.Context)

	// Add this new method
	GenerateQuickLink(c *gin.Context)
	GetPatientTrends(ctx *gin.Context)
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

func (c *PaymentController) CreatePayPalOrder(ctx *gin.Context) {
	var orderRequest OrderRequest
	if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get PayPal access token
	accessToken, err := c.service.getPayPalAccessToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal access token: " + err.Error()})
		return
	}

	result, err := c.service.createPayPalOrder(accessToken, orderRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *PaymentController) CapturePayPalOrder(ctx *gin.Context) {
	orderID := ctx.Param("orderID")
	accessToken, err := c.service.getPayPalAccessToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal access token: " + err.Error()})
		return
	}
	result, err := c.service.capturePayPalOrder(accessToken, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *PaymentController) GetOrderDetails(ctx *gin.Context) {
	orderID := ctx.Param("orderID")
	accessToken, err := c.service.getPayPalAccessToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal access token: " + err.Error()})
		return
	}
	result, err := c.service.getOrderDetails(accessToken, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *PaymentController) UpdateOrder(ctx *gin.Context) {
	orderID := ctx.Param("orderID")
	accessToken, err := c.service.getPayPalAccessToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal access token: " + err.Error()})
		return
	}
	var updates []OrderUpdateRequest
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.updateOrder(accessToken, orderID, updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *PaymentController) TrackOrder(ctx *gin.Context) {
	orderID := ctx.Param("orderID")
	accessToken, err := c.service.getPayPalAccessToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PayPal access token: " + err.Error()})
		return
	}
	result, err := c.service.trackOrder(accessToken, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *PaymentController) CreatePayOSLink(ctx *gin.Context) {
	var orderRequest CreatePaymentLinkRequest
	if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := c.service.createPaymentLink(ctx, orderRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

// Add this to the existing payment.controller.go file

// GetRevenueLastSevenDays handles the request to get revenue data for the last 7 days
func (c *PaymentController) GetRevenueLastSevenDays(ctx *gin.Context) {
	revenueData, err := c.service.GetRevenueLastSevenDays(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, revenueData)
}

// GenerateQuickLink handles the request to generate a VietQR Quick Link
func (c *PaymentController) GenerateQuickLink(ctx *gin.Context) {
	var request QuickLinkRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.service.GenerateQuickLink(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetPatientTrends handles the request to get patient trend data
func (c *PaymentController) GetPatientTrends(ctx *gin.Context) {
	trends, err := c.service.GetPatientTrends(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, trends)
}
