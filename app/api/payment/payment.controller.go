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

	GetRevenueLastSevenDays(ctx *gin.Context)

	// Add this new method
	GenerateQuickLink(c *gin.Context)
	GetPatientTrends(ctx *gin.Context)

	// Add cash payment method
	CreateCashPayment(c *gin.Context)
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

// CreateCashPayment handles the request to create a payment using cash method
func (c *PaymentController) CreateCashPayment(ctx *gin.Context) {
	var request CashPaymentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kiểm tra nếu tiền nhận phải >= số tiền thanh toán
	if request.CashReceived < request.Amount {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tiền nhận phải lớn hơn hoặc bằng số tiền thanh toán"})
		return
	}

	// Tính toán tiền thối lại nếu không được cung cấp
	if request.CashChange == 0 {
		request.CashChange = request.CashReceived - request.Amount
	}

	// Gọi service để xử lý thanh toán
	result, err := c.service.CreateCashPayment(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
