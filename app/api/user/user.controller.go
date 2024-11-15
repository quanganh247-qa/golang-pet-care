package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type UserControllerInterface interface {
	createUser(ctx *gin.Context)
	getUserDetails(ctx *gin.Context)
	getAllUsers(ctx *gin.Context)
	loginUser(ctx *gin.Context)
	logoutUser(ctx *gin.Context)
	verifyEmail(ctx *gin.Context)
	getAccessToken(ctx *gin.Context)
<<<<<<< HEAD
	resendOTP(ctx *gin.Context)
	updatetUser(ctx *gin.Context)
	updatetUserAvatar(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	sessioninfo(ctx *gin.Context)
	userinfo(ctx *gin.Context)
	GetAllRole(ctx *gin.Context)
=======
	createDoctor(ctx *gin.Context)
	addSchedule(ctx *gin.Context)
	getDoctor(ctx *gin.Context)
	insertTimeSlots(ctx *gin.Context)
	getTimeSlots(ctx *gin.Context)
	getAllTimeSlots(ctx *gin.Context)
	updateDoctorAvailableTime(ctx *gin.Context)
	// insertTokenInfo(ctx *gin.Context)
>>>>>>> 79a3bcc (medicine api)
}

func (controller *UserController) createUser(ctx *gin.Context) {
	var req createUserRequest
<<<<<<< HEAD

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
=======
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	fullName := ctx.PostForm("full_name")
	email := ctx.PostForm("email")
	phoneNumber := ctx.PostForm("phone_number")
	address := ctx.PostForm("address")
	role := ctx.PostForm("role")

	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max
>>>>>>> 0fb3f30 (user images)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	// Handle image file
	file, header, err := ctx.Request.FormFile("image")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("image is required")))
		return
	}
	defer file.Close()

	// Get the original image name
	originalImageName := header.Filename

	// Read the file content into a byte array
	dataImage, err := ioutil.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read data image"})
		return
	}
	// get original image

	req.DataImage = dataImage
	req.Username = username
	req.Password = password
	req.FullName = fullName
	req.Email = email
	req.PhoneNumber = phoneNumber
	req.Address = address
	req.Role = role
	req.OriginalImage = originalImageName

	err = controller.service.createUserService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Success", nil))
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

// LoginUser           	godoc
// @Summary 			Login user
// @Description 		Login user
// @Tags 				users
// @Accept  			json
// @Produce  			json
// @Param 				user body loginUserRequest true "User info"
// @Success 			200 {object} loginUSerResponse
// @Router 				/user/login [post]
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
	accessToken, _, err := token.TokenMaker.CreateToken(req.Username,
		map[string]bool{"user": true},
		util.Configs.AccessTokenDuration,
	)
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
<<<<<<< HEAD
<<<<<<< HEAD
	token := ctx.Query("token")
=======
	token := ctx.Param("token")
>>>>>>> 9d28896 (image pet)
=======
	token := ctx.Query("token")
>>>>>>> 8d5618d (feat: update logout)

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
<<<<<<< HEAD
=======
	if util.Configs.DefaultAuthenticationUsername != "" && err != nil {
		cookie, _, err = token.TokenMaker.CreateToken(util.Configs.DefaultAuthenticationUsername, nil, util.Configs.AccessTokenDuration)
	}
>>>>>>> 8d5618d (feat: update logout)
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
<<<<<<< HEAD
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
=======
	var req VerrifyEmailTxParams

	emailID := ctx.Query("email_id")
	secretCode := ctx.Query("secret_code")

	if emailID == "" || secretCode == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("missing email_id or secret_code in query parameters")))
		return
	}

	// convert strign to int 64
	emailIDInt, err := strconv.ParseInt(emailID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("invalid email_id parameter")))
		return
	}

	req = VerrifyEmailTxParams{
		EmailId:    emailIDInt,
		SecretCode: secretCode,
	}

	res, err := controller.service.verifyEmailService(ctx, req)
>>>>>>> 9d28896 (image pet)
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

func (controller *UserController) sessioninfo(ctx *gin.Context) {
	sessionContainer := session.GetSessionFromRequestContext(ctx.Request.Context())
	if sessionContainer == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no session found"})
		return
	}

	sessionData, err := sessionContainer.GetSessionDataInDatabase()
	if err != nil {
		if err = supertokens.ErrorHandler(err, ctx.Request, ctx.Writer); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sessionHandle":      sessionContainer.GetHandle(),
		"userId":             sessionContainer.GetUserID(),
		"accessTokenPayload": sessionContainer.GetAccessTokenPayload(),
		"sessionData":        sessionData,
	})
}

func (controller *UserController) userinfo(ctx *gin.Context) {
	sessionContainer := session.GetSessionFromRequestContext(ctx.Request.Context())
	if sessionContainer == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no session found"})
		return
	}

	userInfo, err := thirdparty.GetUserByID(sessionContainer.GetUserID())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", userInfo))
}

func (controller *UserController) GetAllRole(ctx *gin.Context) {
	res, err := controller.service.GetAllRoleService(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}
<<<<<<< HEAD
=======

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

// func (controller *UserController) insertTokenInfo(ctx *gin.Context) {
// 	var req InsertTokenInfoRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
// 		return
// 	}
// 	authPayload, err := middleware.GetAuthorizationPayload(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
// 		return
// 	}
// 	res, err := controller.service.InsertTokenInfoService(ctx, req, authPayload.Username)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusCreated, util.SuccessResponse("Inserted token info successfull", res))
// }
>>>>>>> 79a3bcc (medicine api)
