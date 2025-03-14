<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
	loginDoctor(ctx *gin.Context)
	getDoctorProfile(ctx *gin.Context)
	getAllDoctor(ctx *gin.Context)
	getShifts(ctx *gin.Context)
	createShift(ctx *gin.Context)
	getShiftByDoctorId(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
	getDoctorById(ctx *gin.Context)
=======
>>>>>>> 4ccd381 (Update appointment flow)
=======
	getDoctorById(ctx *gin.Context)
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
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
	id := ctx.Query("id")
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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)

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
<<<<<<< HEAD
=======
package doctor
>>>>>>> ffc9071 (AI suggestion)
=======
package doctor
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 4ccd381 (Update appointment flow)
=======
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
package doctor
>>>>>>> ffc9071 (AI suggestion)
=======
package doctor
>>>>>>> ada3717 (Docker file)
