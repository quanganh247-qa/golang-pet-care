package test

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
)

func NewTestController(es *elasticsearch.ESService, ws *websocket.WSClientManager) *TestController {
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

// func (c *TestController) CreateTest(ctx *gin.Context) {
// 	var req TestOrderRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	test, err := c.service.CreateTestOrder(ctx, req.PetID, req.DoctorID, req.TestType)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, test)
// }

// func (c *TestController) UpdateTestStatus(ctx *gin.Context) {
// 	testID, err := strconv.ParseInt(ctx.Param("test_id"), 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
// 		return
// 	}

// 	var req struct {
// 		Status string `json:"status" binding:"required"`
// 	}

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := c.service.UpdateTestStatus(ctx, testID, req.Status); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Test status updated successfully"})
// }

// func (c *TestController) AddTestResult(ctx *gin.Context) {
// 	testID, err := strconv.ParseInt(ctx.Param("test_id"), 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid test ID"})
// 		return
// 	}

// 	var result TestResult
// 	if err := ctx.ShouldBindJSON(&result); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := c.service.AddTestResult(ctx, testID, result); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Test result added successfully"})
// }

// func (c *TestController) GetTestsByPetsID(ctx *gin.Context) {
// 	petID, err := strconv.ParseInt(ctx.Query("pet_id"), 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
// 		return
// 	}
// 	tests, err := c.service.GetTestsByPetID(ctx, petID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Get tests successfully",
// 		"data":    tests,
// 	})
// }

func (c *TestController) ListTests(ctx *gin.Context) {
	tests, err := c.service.ListTests(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Get tests successfully",
		"data":    tests,
	})
}

func (c *TestController) CreateTestOrder(ctx *gin.Context) {
	var req TestOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := c.service.CreateTestOrder(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Test order created successfully"})

}
