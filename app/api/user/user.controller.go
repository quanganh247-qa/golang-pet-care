package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type UserControllerInterface interface {
	createUser(ctx *gin.Context)
	getAllUsers(ctx *gin.Context)
	loginUser(ctx *gin.Context)
	verifyEmail(ctx *gin.Context)
	getAccessToken(ctx *gin.Context)
	createDoctor(ctx *gin.Context)
	addSchedule(ctx *gin.Context)
	getDoctor(ctx *gin.Context)
<<<<<<< HEAD
	insertTimeSlots(ctx *gin.Context)
	getTimeSlots(ctx *gin.Context)
	getAllTimeSlots(ctx *gin.Context)
	updateDoctorAvailableTime(ctx *gin.Context)
	insertTokenInfo(ctx *gin.Context)
=======
>>>>>>> 1ada478 (get doctor api)
}

func (controller *UserController) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	res, err := controller.service.createUserService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Success", res))
}

func (controller *UserController) getAllUsers(ctx *gin.Context) {
	res, err := controller.service.getAllUsersService(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (controller *UserController) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	err := controller.service.loginUserService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	refreshToken, _, err := token.TokenMaker.CreateToken(req.Username, nil, util.Configs.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	accessToken, _, err := token.TokenMaker.CreateToken(req.Username, nil, util.Configs.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	host, secure := util.SetCookieSameSite(ctx)

	ctx.SetCookie("refresh_token", refreshToken, int(util.Configs.RefreshTokenDuration), "/", host, secure, true)
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", loginUSerResponse{AccessToken: accessToken}))

}

func (controller *UserController) getAccessToken(ctx *gin.Context) {
	util.SetCookieSameSite(ctx)
	cookie, err := ctx.Cookie("refresh_token")
	// nếu set default username thì sẽ luôn sử dụng username đó để tạo token
	if util.Configs.DefaultAuthenticationUsername != "" && err != nil {
		cookie, _, err = token.TokenMaker.CreateToken(util.Configs.DefaultAuthenticationUsername, nil, util.Configs.AccessTokenDuration)
	}
	if err != nil {
		ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
		return
	}
	if cookie == "" {
		ctx.JSON(http.StatusForbidden, util.ErrorResponse(errors.New("refresh_token is not provided")))
		return
	}
	payload, err := token.TokenMaker.VerifyToken(cookie)
	if err != nil {
		ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
		return
	}
	accessToken, _, err := token.TokenMaker.CreateToken(payload.Username, nil, util.Configs.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusForbidden, util.ErrorResponse(err))
		return
	}
	res := loginUSerResponse{AccessToken: accessToken}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (controller *UserController) verifyEmail(ctx *gin.Context) {
	var req VerrifyEmailTxParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	res, err := controller.service.verifyEmailService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Veried user", res))

}

func (controller *UserController) createDoctor(ctx *gin.Context) {
	var req InsertDoctorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.createDoctorService(ctx, req, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Success", res))
}

func (controller *UserController) addSchedule(ctx *gin.Context) {
	var req InsertDoctorScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.createDoctorScheduleService(ctx, req, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Inserted schedulke successfull", res))
}

func (controller *UserController) getDoctor(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	user, err := db.StoreDB.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.getDoctorByID(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}
<<<<<<< HEAD

func (controller *UserController) insertTimeSlots(ctx *gin.Context) {
	var req db.InsertTimeslotParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := controller.service.insertTimeSlots(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Inserted timeslot successfull", res))

}

func (controller *UserController) getTimeSlots(ctx *gin.Context) {
	doctorID := ctx.Param("doctor_id")
	if doctorID == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("doctor_id is required")))
		return
	}
	// day
	day := ctx.Query("day")
	// convert sitrng to int
	doctorIDInt, err := strconv.Atoi(doctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.GetTimeslotsAvailable(ctx, int64(doctorIDInt), day)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("All list of time slots are available", res))
}

func (controller *UserController) getAllTimeSlots(ctx *gin.Context) {
	doctorID := ctx.Param("doctor_id")
	if doctorID == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("doctor_id is required")))
		return
	}
	// day
	day := ctx.Query("day")
	// convert sitrng to int
	doctorIDInt, err := strconv.Atoi(doctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.GetTimeslotsAvailable(ctx, int64(doctorIDInt), day)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("All list of time slots are available", res))
}

func (controller *UserController) updateDoctorAvailableTime(ctx *gin.Context) {
	timeID := ctx.Param("id")
	if timeID == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("Time slot is required")))
		return
	}
	// convert sitrng to int64
	timeSlotId, err := strconv.ParseInt(timeID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	err = controller.service.UpdateDoctorAvailable(ctx, timeSlotId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Updated timeslot successfull", nil))
}

func (controller *UserController) insertTokenInfo(ctx *gin.Context) {
	var req InsertTokenInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.InsertTokenInfoService(ctx, req, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Inserted token info successfull", res))
}
=======
>>>>>>> 1ada478 (get doctor api)
