package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
type PaymentControllerInterface interface {
	GetToken(c *gin.Context)
	GetBanks(c *gin.Context)
	GenerateQRCode(c *gin.Context)

	CreatePayPalOrder(c *gin.Context)
	CapturePayPalOrder(c *gin.Context)
	GetOrderDetails(c *gin.Context)
	UpdateOrder(c *gin.Context)
	TrackOrder(c *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)

	CreatePayOSLink(c *gin.Context)
}

func (c *PaymentController) GetToken(ctx *gin.Context) {
=======
type VietQRControllerInterface interface {
=======
type PaymentControllerInterface interface {
>>>>>>> e859654 (Elastic search)
	GetToken(c *gin.Context)
	GetBanks(c *gin.Context)
	GenerateQRCode(c *gin.Context)
=======
>>>>>>> ada3717 (Docker file)
}

<<<<<<< HEAD
func (c *VietQRController) GetToken(ctx *gin.Context) {
>>>>>>> c449ffc (feat: cart api)
=======
func (c *PaymentController) GetToken(ctx *gin.Context) {
>>>>>>> e859654 (Elastic search)
=======
type VietQRControllerInterface interface {
=======
type PaymentControllerInterface interface {
>>>>>>> e859654 (Elastic search)
	GetToken(c *gin.Context)
	GetBanks(c *gin.Context)
	GenerateQRCode(c *gin.Context)
=======
>>>>>>> ada3717 (Docker file)
}

<<<<<<< HEAD
func (c *VietQRController) GetToken(ctx *gin.Context) {
>>>>>>> c449ffc (feat: cart api)
=======
func (c *PaymentController) GetToken(ctx *gin.Context) {
>>>>>>> e859654 (Elastic search)
	result, err := c.service.GetToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (c *PaymentController) GetBanks(ctx *gin.Context) {
=======
func (c *VietQRController) GetBanks(ctx *gin.Context) {
>>>>>>> c449ffc (feat: cart api)
=======
func (c *PaymentController) GetBanks(ctx *gin.Context) {
>>>>>>> e859654 (Elastic search)
=======
func (c *VietQRController) GetBanks(ctx *gin.Context) {
>>>>>>> c449ffc (feat: cart api)
=======
func (c *PaymentController) GetBanks(ctx *gin.Context) {
>>>>>>> e859654 (Elastic search)
	result, err := c.service.GetBanksService(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (c *PaymentController) GenerateQRCode(ctx *gin.Context) {
=======
func (c *VietQRController) GenerateQRCode(ctx *gin.Context) {
>>>>>>> c449ffc (feat: cart api)
=======
func (c *PaymentController) GenerateQRCode(ctx *gin.Context) {
>>>>>>> e859654 (Elastic search)
=======
func (c *VietQRController) GenerateQRCode(ctx *gin.Context) {
>>>>>>> c449ffc (feat: cart api)
=======
func (c *PaymentController) GenerateQRCode(ctx *gin.Context) {
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> ada3717 (Docker file)

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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)

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
<<<<<<< HEAD
=======
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)
=======
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> ada3717 (Docker file)
