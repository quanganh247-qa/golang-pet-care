package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type UserControllerInterface interface {
	createUser(ctx *gin.Context)

}

func (controller *UserController) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	// authPayload, err := middleware.GetAuthorizationPayload(ctx)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
	// 	return
	// }
	res, err := controller.service.createUserService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Success", res))
}