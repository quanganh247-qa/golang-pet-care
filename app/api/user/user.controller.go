package user

import (
	"encoding/json"
	"errors"
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
	resendOTP(ctx *gin.Context)
	updatetUser(ctx *gin.Context)
	updatetUserAvatar(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	sessioninfo(ctx *gin.Context)
	userinfo(ctx *gin.Context)
	GetAllRole(ctx *gin.Context)
	createStaff(ctx *gin.Context)
}

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
	// token := ctx.Query("token")

	var req LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	err = controller.service.logoutUsersService(ctx, authPayload.Username, req.Token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", nil))
}

func (controller *UserController) getAccessToken(ctx *gin.Context) {
	util.SetCookieSameSite(ctx)
	cookie, err := ctx.Cookie("refresh_token")
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

func (controller *UserController) createStaff(ctx *gin.Context) {
	// First check if the requesting user has admin role
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, util.ErrorResponse(err))
		return
	}

	// Get user details to check role
	user, err := controller.service.getUserDetailsService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	// Only admin can create staff accounts
	if user.Role != "admin" {
		ctx.JSON(http.StatusForbidden, util.ErrorResponse(errors.New("only admin can create staff accounts")))
		return
	}

	var req CreateStaffRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	// Set verified email to true for staff by default
	req.IsVerifiedEmail = true

	res, err := controller.service.createStaffService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Staff created successfully", res))
}
