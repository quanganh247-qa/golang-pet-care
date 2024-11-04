package pet

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

type PetControllerInterface interface {
	CreatePet(ctx *gin.Context)
	GetPetByID(ctx *gin.Context)
	ListPets(ctx *gin.Context)
	ListPetsByUsername(ctx *gin.Context)
	UpdatePet(ctx *gin.Context)
	DeletePet(ctx *gin.Context)
}

func (c *PetController) CreatePet(ctx *gin.Context) {
	var req createPetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)

	res, err := c.service.CreatePet(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}

func (c *PetController) GetPetByID(ctx *gin.Context) {
	petidStr := ctx.Param("petid")
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
	// Lấy các tham số từ query, ví dụ như limit và offset
	limitStr := ctx.DefaultQuery("limit", "10")  // Giá trị mặc định là 10
	offsetStr := ctx.DefaultQuery("offset", "0") // Giá trị mặc định là 0

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	req := listPetsRequest{
		Type:   ctx.Query("type"),
		Breed:  ctx.Query("breed"),
		Age:    int(limit),
		Weight: float64(offset),
	}

	pets, err := c.service.ListPets(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pets)
}

func (c *PetController) UpdatePet(ctx *gin.Context) {
	petid, err := strconv.ParseInt(ctx.Param("petid"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	var req createPetRequest
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
	petid, err := strconv.ParseInt(ctx.Param("petid"), 10, 64)
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
	username := ctx.Param("username")
	limitStr := ctx.DefaultQuery("limit", "10")
	offsetStr := ctx.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	offset, err := strconv.ParseInt(offsetStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
		return
	}

	pets, err := c.service.ListPetsByUsername(ctx, username, int32(limit), int32(offset))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pets)
}
