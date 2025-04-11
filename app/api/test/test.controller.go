package test

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
)

func NewTestController(ws *websocket.WSClientManager) *TestController {
	service := NewTestService(db.StoreDB, ws)
	return &TestController{
		service: service,
		ws:      ws,
	}
}

func (c *TestController) HandleWebSocket(ctx *gin.Context) {
	// Get doctor ID from context
	doctorID, exists := ctx.Get("doctor_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convert interface{} to int64
	doctorIDInt, ok := doctorID.(int64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID format"})
		return
	}

	// Set client ID in query parameters
	clientID := fmt.Sprintf("doctor_%d", doctorIDInt)
	ctx.Request.Header.Set("X-Client-ID", clientID)

	// Handle WebSocket connection
	c.ws.HandleWebSocket(ctx)
}

// ListAllItems lists all tests and vaccines
func (c *TestController) ListAllItems(ctx *gin.Context) {
	// Get type filter from query parameters (optional)
	itemType := ctx.Query("type")

	items, err := c.service.ListItems(ctx, ItemType(itemType))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Items retrieved successfully",
		"data":    items,
	})
}

// ListTests lists only tests (for backward compatibility)
func (c *TestController) ListTests(ctx *gin.Context) {
	items, err := c.service.ListItems(ctx, TypeTest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tests retrieved successfully",
		"data":    items,
	})
}

// ListVaccines lists only vaccines
func (c *TestController) ListVaccines(ctx *gin.Context) {
	items, err := c.service.ListItems(ctx, TypeVaccine)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Vaccines retrieved successfully",
		"data":    items,
	})
}

// CreateOrder creates a new order for tests/vaccines
func (c *TestController) CreateOrder(ctx *gin.Context) {
	var req TestOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.CreateOrder(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

// CreateTestOrder is for backward compatibility
func (c *TestController) CreateTestOrder(ctx *gin.Context) {
	var req TestOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default type to test if not provided
	// if req.ItemType == "" {
	// 	req.ItemType = TypeTest
	// }

	err := c.service.CreateOrder(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Test order created successfully"})
}

// GetItemByID gets a test or vaccine by ID
func (c *TestController) GetItemByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	item, err := c.service.GetItemByID(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item details retrieved successfully",
		"data":    item,
	})
}

// GetTestByID is for backward compatibility
func (c *TestController) GetTestByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
		return
	}

	item, err := c.service.GetItemByID(ctx, int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test details retrieved successfully",
		"data":    item,
	})
}

// SoftDeleteItem removes an item (test or vaccine)
func (c *TestController) SoftDeleteItem(ctx *gin.Context) {
	itemID := ctx.Param("item_id")
	if itemID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	err := c.service.SoftDeleteItem(ctx, itemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Item deleted successfully",
	})
}

// SoftDeleteTest is for backward compatibility
func (c *TestController) SoftDeleteTest(ctx *gin.Context) {
	testID := ctx.Param("test_id")
	if testID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Test ID is required"})
		return
	}

	err := c.service.SoftDeleteItem(ctx, testID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test deleted successfully",
	})
}

// GetOrderedItemsByAppointment gets ordered items by appointment ID
func (c *TestController) GetOrderedItemsByAppointment(ctx *gin.Context) {
	appointmentID := ctx.Query("appointment_id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "appointment_id is required"})
		return
	}

	// Get optional type filter
	itemType := ctx.Query("type")

	// Convert appointmentID to int64
	appointmentIDInt, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	orderedItems, err := c.service.GetOrderedItemsByAppointment(ctx, appointmentIDInt, ItemType(itemType))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ordered items retrieved successfully",
		"data":    orderedItems,
	})
}

// GetOrderedTestsByAppointment is for backward compatibility
func (c *TestController) GetOrderedTestsByAppointment(ctx *gin.Context) {
	appointmentID := ctx.Query("appointment_id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "appointment_id is required"})
		return
	}

	// Convert appointmentID to int64
	appointmentIDInt, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	orderedItems, err := c.service.GetOrderedItemsByAppointment(ctx, appointmentIDInt, TypeTest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ordered tests retrieved successfully",
		"data":    orderedItems,
	})
}

func (c *TestController) GetAllAppointmentsWithOrders(ctx *gin.Context) {

	appointments, err := c.service.GetAllAppointmentsWithOrders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get appointments successfully",
		"data":    appointments,
	})
}

func (c *TestController) GetTestByAppointment(ctx *gin.Context) {
	appointmentID := ctx.Query("appointment_id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "appointment_id is required"})
		return
	}

	// Convert appointmentID to int64
	appointmentIDInt, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	tests, err := c.service.GetTestByAppointment(ctx, appointmentIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tests retrieved successfully",
		"data":    tests,
	})
}
