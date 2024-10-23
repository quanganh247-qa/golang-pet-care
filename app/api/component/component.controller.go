package component

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ComponentControllerInterface interface {
	createComponent(ctx *gin.Context)
	getComponentByID(ctx *gin.Context)
	getComponentsByUser(ctx *gin.Context)
	removeComponent(ctx *gin.Context)
}

func (c *ComponentController) createComponent(ctx *gin.Context) {
	var req createComponentsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	res, err := c.service.createComponent(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Component is created sucessfully", res))
}

func (c *ComponentController) getComponentByID(ctx *gin.Context) {
	id := ctx.Param("id")
	idNumber, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	projectId := ctx.Request.URL.Query().Get("project_id")
	project_id, err := strconv.ParseInt(projectId, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.getComponentByID(ctx, idNumber, int32(project_id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (c *ComponentController) getComponentsByUser(ctx *gin.Context) {

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	authenticatePayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	projectId := ctx.Request.URL.Query().Get("project_id")
	project_id, err := strconv.ParseInt(projectId, 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.getComponentsByUser(ctx, authenticatePayload.Username, int32(project_id), *pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Get Components success", res))
}

func (c *ComponentController) removeComponent(ctx *gin.Context) {
	id := ctx.Param("id")
	idNumber, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	projectId := ctx.Request.URL.Query().Get("project_id")
	project_id, err := strconv.ParseInt(projectId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.removeComponent(ctx, idNumber, int32(project_id))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Removed sucessfully", res))
}
