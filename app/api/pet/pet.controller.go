package pet

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetControllerInterface interface {
	CreatePet(ctx *gin.Context)
	GetPetByID(ctx *gin.Context)
	ListPets(ctx *gin.Context)
	ListPetsByUsername(ctx *gin.Context)
	UpdatePet(ctx *gin.Context)
	DeletePet(ctx *gin.Context)
	GetPetLogsByPetID(ctx *gin.Context)
	InsertPetLog(ctx *gin.Context)
	DeletePetLog(ctx *gin.Context)
	UpdatePetLog(ctx *gin.Context)
	UpdatePetAvatar(ctx *gin.Context)
	GetAllPetLogsByUsername(ctx *gin.Context)
	GetPetOwnerByPetID(ctx *gin.Context)
	GetDetailsPetLog(ctx *gin.Context)

	// Weight tracking methods
	GetWeightHistory(ctx *gin.Context)
	DeleteWeightRecord(ctx *gin.Context)
}

func (c *PetController) CreatePet(ctx *gin.Context) {
	var req createPetRequest

	// Parse the JSON data from the "data" form field
	jsonData := ctx.PostForm("data")
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req.OriginalImage = originalImageName
	req.DataImage = dataImage

	res, err := c.service.CreatePet(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

func (c *PetController) GetPetByID(ctx *gin.Context) {
	petidStr := ctx.Param("pet_id")
	petid, err := strconv.ParseInt(petidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}

	res, err := c.service.GetPetByID(ctx, petid)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

func (c *PetController) ListPets(ctx *gin.Context) {

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	req := listPetsRequest{
		Type:  ctx.Query("type"),
		Breed: ctx.Query("breed"),
		// Age:    int(limit),
		// Weight: float64(offset),
	}

	pets, err := c.service.ListPets(ctx, req, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the pets with proper pagination response structure
	ctx.JSON(http.StatusOK, pets)
}

func (c *PetController) UpdatePet(ctx *gin.Context) {
	petid, err := strconv.ParseInt(ctx.Param("pet_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	var req updatePetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.UpdatePet(ctx, petid, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Pet updated successfully"})
}

func (c *PetController) DeletePet(ctx *gin.Context) {
	petid, err := strconv.ParseInt(ctx.Param("pet_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	err = c.service.SetPetInactive(ctx, petid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Pet set to inactive successfully"})
}

func (c *PetController) ListPetsByUsername(ctx *gin.Context) {
	// username := ctx.Param("username")
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pets, err := c.service.ListPetsByUsername(ctx, authPayload.Username, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pets)
}

func (c *PetController) GetPetLogsByPetID(ctx *gin.Context) {
	petidStr := ctx.Param("pet_id")
	petid, err := strconv.ParseInt(petidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	res, err := c.service.GetPetLogsByPetIDService(ctx, petid, pagination)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

func (c *PetController) InsertPetLog(ctx *gin.Context) {
	var req PetLogWithPetInfo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := c.service.InsertPetLogService(ctx, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Invalidate any log-related cache after insertion
	middleware.InvalidateCache("pet_logs")

	ctx.JSON(200, gin.H{"message": "Insert pet log successfully"})
}

func (c *PetController) DeletePetLog(ctx *gin.Context) {
	logidStr := ctx.Param("log_id")
	logid, err := strconv.ParseInt(logidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid log ID"})
		return
	}

	err = c.service.DeletePetLogService(ctx, logid)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Invalidate any log-related cache after deletion
	middleware.InvalidateCache("pet_logs")

	ctx.JSON(200, gin.H{"message": "Delete pet log successfully"})
}

func (c *PetController) UpdatePetLog(ctx *gin.Context) {
	var req UpdatePetLogRequeststruct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logidStr := ctx.Param("log_id")
	logid, err := strconv.ParseInt(logidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid log ID"})
		return
	}

	err = c.service.UpdatePetLogService(ctx, req, logid)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Invalidate any log-related cache after update
	middleware.InvalidateCache("pet_logs")

	ctx.JSON(200, gin.H{"message": "Update pet log successfully"})
}

func (c *PetController) UpdatePetAvatar(ctx *gin.Context) {
	petidStr := ctx.Param("pet_id")
	petid, err := strconv.ParseInt(petidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}

	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(400, util.ErrorResponse(err))
		return
	}

	err = c.service.UpdatePetAvatar(ctx, petid, updatePetAvatarRequest{
		DataImage:     dataImage,
		OriginalImage: originalImageName,
	})
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Update pet avatar successfully"})
}

// Add this method to your PetController implementation
func (c *PetController) GetAllPetLogsByUsername(ctx *gin.Context) {
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logs, err := c.service.GetAllPetLogsByUsername(ctx, authPayload.Username, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, logs)
}

func (c *PetController) GetPetOwnerByPetID(ctx *gin.Context) {
	petidStr := ctx.Param("pet_id")
	petid, err := strconv.ParseInt(petidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}
	res, err := c.service.GetPetOwnerByPetID(ctx, petid)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, res)
}

func (c *PetController) GetDetailsPetLog(ctx *gin.Context) {
	logidStr := ctx.Param("log_id")
	logid, err := strconv.ParseInt(logidStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid log ID"})
		return
	}

	res, err := c.service.GetDetailsPetLogService(ctx, logid)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

// GetWeightHistory gets the weight history for a pet
// @Summary Get weight history
// @Description Get the weight history for a pet
// @Tags pets
// @Produce json
// @Param pet_id path int true "Pet ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Param unit_type query string false "Unit type (kg or lb, default: kg)"
// @Success 200 {object} PetWeightHistoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pets/{pet_id}/weight [get]
func (c *PetController) GetWeightHistory(ctx *gin.Context) {
	// Get pet ID from path parameter
	petIDStr := ctx.Param("pet_id")
	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil || petID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	unitType := ctx.DefaultQuery("unit_type", "kg")

	// Validate unit type
	if unitType != "kg" && unitType != "lb" {
		unitType = "kg" // Default to kg
	}

	response, err := c.service.GetWeightHistory(ctx, petID, pagination, unitType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteWeightRecord deletes a weight record for a pet
// @Summary Delete weight record
// @Description Delete a weight record for a pet
// @Tags pets
// @Produce json
// @Param record_id path int true "Record ID"
// @Param pet_id path int true "Pet ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /pets/{pet_id}/weight/{record_id} [delete]
func (c *PetController) DeleteWeightRecord(ctx *gin.Context) {
	// Get record ID from path parameter
	recordIDStr := ctx.Param("record_id")
	recordID, err := strconv.ParseInt(recordIDStr, 10, 64)
	if err != nil || recordID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	// Get pet ID from path parameter
	petIDStr := ctx.Param("pet_id")
	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil || petID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	err = c.service.DeleteWeightRecord(ctx, recordID, petID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Weight record deleted successfully"})
}
