package user

import (
	"encoding/json"
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
	getUserDetails(ctx *gin.Context)
	getAllUsers(ctx *gin.Context)
	loginUser(ctx *gin.Context)
	logoutUser(ctx *gin.Context)
	verifyEmail(ctx *gin.Context)
	getAccessToken(ctx *gin.Context)
	createDoctor(ctx *gin.Context)
	addSchedule(ctx *gin.Context)
	getDoctor(ctx *gin.Context)
	insertTimeSlots(ctx *gin.Context)
	getTimeSlots(ctx *gin.Context)
	getAllTimeSlots(ctx *gin.Context)
	updateDoctorAvailableTime(ctx *gin.Context)
	resendOTP(ctx *gin.Context)
	updatetUser(ctx *gin.Context)
	updatetUserAvatar(ctx *gin.Context)
	GetDoctors(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
}

// createUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags user
// @Accept  json
// @Produce  json
// @Param data formData string true "User data"
// @Param image formData file true "User image"
// @Success 201 {object} VerrifyEmailTxParams
// @Router /user [post]
// @Security ApiKeyAuth
// @in header
// @name Authorization
func (controller *UserController) createUser(ctx *gin.Context) {
	var req createUserRequest

	// Parse the JSON data from the "data" form field
	jsonData := ctx.PostForm("data")
	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Use the helper function to handle the image upload
	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	req.DataImage = dataImage
	req.OriginalImage = originalImageName

	res, err := controller.service.createUserService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Success", res))
}

func (controller *UserController) getUserDetails(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	user, err := controller.service.getUserDetailsService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", user))
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
	user, err := controller.service.loginUserService(ctx, req)
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
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	host, secure := util.SetCookieSameSite(ctx)

	ctx.SetCookie("refresh_token", refreshToken, int(util.Configs.RefreshTokenDuration), "/", host, secure, true)
	ctx.JSON(http.StatusOK, util.SuccessResponse("Login succesfully", user))
}

func (controller *UserController) logoutUser(ctx *gin.Context) {
	token := ctx.Query("token")

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	err = controller.service.logoutUsersService(ctx, authPayload.Username, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", nil))
}

func (controller *UserController) getAccessToken(ctx *gin.Context) {
	util.SetCookieSameSite(ctx)
	cookie, err := ctx.Cookie("refresh_token")
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
	var req VerrifyInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	otpInt, err := strconv.ParseInt(req.SecretCode, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	arg := VerrifyEmailTxParams{
		SecretCode: otpInt,
		Username:   req.Username,
	}

	err = controller.service.verifyEmailService(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Email verified", nil))

}

func (controller *UserController) resendOTP(ctx *gin.Context) {
	username := ctx.Param("username")
	if username == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(errors.New("username is required")))
		return
	}
	res, err := controller.service.resendOTPService(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Resend OTP successfull", res))
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

func (controller *UserController) updatetUser(ctx *gin.Context) {

	var arg UpdateUserParams
	if err := ctx.ShouldBindJSON(&arg); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.updateUserService(ctx, authPayload.Username, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Update user succesfully", res))
}

func (controller *UserController) updatetUserAvatar(ctx *gin.Context) {

	// Use the helper function to handle the image upload
	dataImage, originalImageName, err := util.HandleImageUpload(ctx, "image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	req := UpdateUserImageParams{
		DataImage:     dataImage,
		OriginalImage: originalImageName,
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	err = controller.service.updateUserImageService(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", nil))
}

func (controller *UserController) GetDoctors(ctx *gin.Context) {
	// pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
	// 	return
	// }

	res, err := controller.service.GetDoctorsService(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (controller *UserController) ForgotPassword(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	err := controller.service.ForgotPasswordService(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", nil))
}

func (controller *UserController) UpdatePassword(ctx *gin.Context) {
	var req UpdatePasswordParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	err = controller.service.UpdatePasswordService(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", nil))
}
