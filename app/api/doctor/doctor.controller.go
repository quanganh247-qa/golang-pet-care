package doctor

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
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
	createLeaveRequest(ctx *gin.Context)
	getLeaveRequests(ctx *gin.Context)
	updateLeaveRequest(ctx *gin.Context)
	getDoctorAttendance(ctx *gin.Context)
	getDoctorWorkload(ctx *gin.Context)
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

func (c *DoctorController) createLeaveRequest(ctx *gin.Context) {
	// Get doctor ID from params
	doctorIDStr := ctx.Param("doctor_id")
	doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var req CreateLeaveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	leave, err := c.service.CreateLeaveRequest(ctx, doctorID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, util.SuccessResponse("leave request created successfully", leave))
}

func (c *DoctorController) getLeaveRequests(ctx *gin.Context) {
	doctorIDStr := ctx.Param("doctor_id")
	doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	leaves, err := c.service.GetLeaveRequests(ctx, doctorID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("leave requests retrieved successfully", leaves))
}

func (c *DoctorController) updateLeaveRequest(ctx *gin.Context) {
	leaveIDStr := ctx.Param("leave_id")
	leaveID, err := strconv.ParseInt(leaveIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	// Get reviewer info from auth token
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	var req UpdateLeaveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	// Get reviewer's user ID
	user, err := db.StoreDB.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	err = c.service.UpdateLeaveRequest(ctx, leaveID, user.ID, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("leave request updated successfully", nil))
}

func (c *DoctorController) getDoctorAttendance(ctx *gin.Context) {
	doctorIDStr := ctx.Param("doctor_id")
	doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	startDate, err := time.Parse("2006-01-02", ctx.Query("start_date"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("invalid start_date format")))
		return
	}

	endDate, err := time.Parse("2006-01-02", ctx.Query("end_date"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("invalid end_date format")))
		return
	}

	attendance, err := c.service.GetDoctorAttendance(ctx, doctorID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("attendance records retrieved successfully", attendance))
}

func (c *DoctorController) getDoctorWorkload(ctx *gin.Context) {
	doctorIDStr := ctx.Param("doctor_id")
	doctorID, err := strconv.ParseInt(doctorIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	startDate, err := time.Parse("2006-01-02", ctx.Query("start_date"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("invalid start_date format")))
		return
	}

	endDate, err := time.Parse("2006-01-02", ctx.Query("end_date"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("invalid end_date format")))
		return
	}

	workload, err := c.service.GetDoctorWorkload(ctx, doctorID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("workload metrics retrieved successfully", workload))
}
