package doctor

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DoctorControllerInterface interface {
	editDoctorProfile(ctx *gin.Context)
	loginDoctor(ctx *gin.Context)
	getDoctorProfile(ctx *gin.Context)
	getAllDoctor(ctx *gin.Context)
	getShifts(ctx *gin.Context)
	createShift(ctx *gin.Context)
	getShiftByDoctorId(ctx *gin.Context)
	getDoctorById(ctx *gin.Context)
	deleteShift(ctx *gin.Context)
}

func (c *DoctorController) loginDoctor(ctx *gin.Context) {
	var req loginDoctorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	res, err := c.service.LoginDoctorService(ctx, req)
	if err != nil {
		// Error responses are handled in the service
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("login successful", res))
}

func (c *DoctorController) editDoctorProfile(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}
	var req EditDoctorProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	if err := c.service.EditDoctorProfileService(ctx, authPayload.Username, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("doctor profile updated successfully", nil))
}

func (c *DoctorController) getDoctorProfile(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	profile, err := c.service.GetDoctorProfile(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("doctor profile retrieved successfully", profile))
}

func (c *DoctorController) getAllDoctor(ctx *gin.Context) {
	res, err := c.service.GetAllDoctorService(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (c *DoctorController) getShifts(ctx *gin.Context) {
	res, err := c.service.GetShifts(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (c *DoctorController) createShift(ctx *gin.Context) {
	var req CreateShiftRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	res, err := c.service.CreateShift(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (c *DoctorController) getShiftByDoctorId(ctx *gin.Context) {
	id := ctx.Param("doctor_id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("id is required")))
		return
	}
	doctorId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetShiftByDoctorId(ctx, doctorId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (c *DoctorController) getDoctorById(ctx *gin.Context) {
	id := ctx.Param("doctor_id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("id is required")))
		return
	}
	doctorId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetDoctorById(ctx, doctorId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (c *DoctorController) deleteShift(ctx *gin.Context) {
	id := ctx.Param("shift_id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("id is required")))
		return
	}
	shiftId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	err = c.service.DeleteShift(ctx, shiftId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", nil))
}
