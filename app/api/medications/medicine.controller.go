package medications

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicineControllerInterface interface {
	CreateMedicine(ctx *gin.Context)
	GetMedicineByID(ctx *gin.Context)
	ListMedicines(ctx *gin.Context)
	UpdateMedicine(ctx *gin.Context)
	// DeleteMedicine(ctx *gin.Context)
	GetAllMedicines(ctx *gin.Context)

	// Inventory management methods
	CreateMedicineTransaction(ctx *gin.Context)
	GetMedicineTransactions(ctx *gin.Context)
	GetMedicineTransactionsByMedicineID(ctx *gin.Context)

	// Supplier management
	CreateSupplier(ctx *gin.Context)
	GetSupplierByID(ctx *gin.Context)
	GetAllSuppliers(ctx *gin.Context)
	UpdateSupplier(ctx *gin.Context)

	// Expiration alerts
	GetExpiringMedicines(ctx *gin.Context)
	GetLowStockMedicines(ctx *gin.Context)

	// WebSocket connection
	HandleWebSocket(ctx *gin.Context)
}

func (c *MedicineController) CreateMedicine(ctx *gin.Context) {
	var req createMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.CreateMedicine(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *MedicineController) GetMedicineByID(ctx *gin.Context) {
	medicineidStr := ctx.Param("medicine_id")
	medicineid, err := strconv.ParseInt(medicineidStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Medicine ID"})
		return
	}

	res, err := c.service.GetMedicineByID(ctx, medicineid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *MedicineController) ListMedicines(ctx *gin.Context) {

	petID := ctx.Param("pet_id")
	petIDInt, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Pet ID"})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	Medicines, err := c.service.ListMedicines(ctx, pagination, petIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Medicines)
}

func (c *MedicineController) UpdateMedicine(ctx *gin.Context) {
	medicineid, err := strconv.ParseInt(ctx.Param("medicine_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Medicine ID"})
		return
	}

	var req createMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(req)

	err = c.service.UpdateMedicine(ctx, medicineid, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medicine updated successfully"})
}

// func (c *MedicineController) ListMedicinesByUsername(ctx *gin.Context) {
// 	// username := ctx.Param("username")
// 	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
// 		return
// 	}
// 	authPayload, err := middleware.GetAuthorizationPayload(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	fmt.Println(authPayload.Username)
// 	Medicines, err := c.service.ListMedicinesByUsername(ctx, authPayload.Username, pagination)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, Medicines)
// }

func (c *MedicineController) GetAllMedicines(ctx *gin.Context) {

	Medicines, err := c.service.GetAllMedicines(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Medicines)

}

// Medicine Inventory Management

func (c *MedicineController) CreateMedicineTransaction(ctx *gin.Context) {
	var req MedicineTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.CreateMedicineTransaction(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *MedicineController) GetMedicineTransactions(ctx *gin.Context) {
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	transactions, err := c.service.GetMedicineTransactions(ctx, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (c *MedicineController) GetMedicineTransactionsByMedicineID(ctx *gin.Context) {
	medicineID, err := strconv.ParseInt(ctx.Param("medicine_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Medicine ID"})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	transactions, err := c.service.GetMedicineTransactionsByMedicineID(ctx, medicineID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

// Supplier Management

func (c *MedicineController) CreateSupplier(ctx *gin.Context) {
	var req MedicineSupplierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.service.CreateSupplier(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *MedicineController) GetSupplierByID(ctx *gin.Context) {
	supplierID, err := strconv.ParseInt(ctx.Param("supplier_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Supplier ID"})
		return
	}

	supplier, err := c.service.GetSupplierByID(ctx, supplierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, supplier)
}

func (c *MedicineController) GetAllSuppliers(ctx *gin.Context) {
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	suppliers, err := c.service.GetAllSuppliers(ctx, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, suppliers)
}

func (c *MedicineController) UpdateSupplier(ctx *gin.Context) {
	supplierID, err := strconv.ParseInt(ctx.Param("supplier_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Supplier ID"})
		return
	}

	var req MedicineSupplierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.UpdateSupplier(ctx, supplierID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Supplier updated successfully"})
}

// Expiration and Stock Alerts

func (c *MedicineController) GetExpiringMedicines(ctx *gin.Context) {
	daysParam := ctx.Query("days")
	days, err := strconv.Atoi(daysParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
		return
	}

	expiringMedicines, err := c.service.GetExpiringMedicines(ctx, days)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, expiringMedicines)
}

func (c *MedicineController) GetLowStockMedicines(ctx *gin.Context) {
	lowStockMedicines, err := c.service.GetLowStockMedicines(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, lowStockMedicines)
}

// WebSocket handler for real-time notifications

func (c *MedicineController) HandleWebSocket(ctx *gin.Context) {
	// Get user information for client ID
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Set client ID in query parameters
	clientID := fmt.Sprintf("user_%s", authPayload.Username)
	ctx.Request.Header.Set("X-Client-ID", clientID)

	// Handle WebSocket connection
	c.service.HandleWebSocket(ctx)
}
