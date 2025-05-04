package device_token

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DeviceTokenControllerInterface interface {
	insertDeviceToken(ctx *gin.Context)
	deleteDeviceToken(ctx *gin.Context)
	getDeviceTokens(ctx *gin.Context)
}

func (c *DeviceTokenController) insertDeviceToken(ctx *gin.Context) {
	var req DVTRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	// Call service to insert device token
	token, err := c.service.InsertToken(ctx, req, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Device token inserted successfully", token))

}

func (c *DeviceTokenController) deleteDeviceToken(ctx *gin.Context) {

	token := ctx.Param("token")
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	err = c.service.DeleteDevicetToken(ctx, authPayload.Username, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Device token deleted successfully", nil))
}

func (c *DeviceTokenController) getDeviceTokens(ctx *gin.Context) {

	var req DVTList
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	user, err := db.StoreDB.GetUserByID(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	tokens, err := c.service.GetDeviceTokenByUsername(ctx, user.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Get device token successfully", tokens))
}
