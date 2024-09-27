package project

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ProjectControllerInterface interface {
	createProject(ctx *gin.Context)
}

func (c *ProjectController) createProject(ctx *gin.Context) {
	var req createProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	fmt.Println("HAHAHA ", authPayload.Username)

	res, err := c.service.createProject(ctx, req, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Create sucessfully", res))
}
