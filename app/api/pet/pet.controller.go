package pet

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

type PetControllerInterface interface {
	CreatePet(ctx *gin.Context)
	GetPet(ctx *gin.Context)
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

func (c *PetController) GetPet(ctx *gin.Context) {

	petId := ctx.Param("petId")
	petIdInt, err := strconv.ParseInt(petId, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid pet ID"})
		return
	}
	res, err := c.service.GetPet(ctx, petIdInt)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, res)
}
