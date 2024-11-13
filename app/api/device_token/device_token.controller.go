package device_token

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DeviceTokenControllerInterface interface {
	insertDeviceToken(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	deleteDeviceToken(ctx *gin.Context)
=======
>>>>>>> 0fb3f30 (user images)
=======
	deleteDeviceToken(ctx *gin.Context)
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 0fb3f30 (user images)
=======
	deleteDeviceToken(ctx *gin.Context)
>>>>>>> 9d28896 (image pet)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 9d28896 (image pet)

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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> 9d28896 (image pet)
