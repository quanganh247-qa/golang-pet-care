package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/mail"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type UserServiceInterface interface {
	createUserService(ctx *gin.Context, req createUserRequest) (*VerrifyEmailTxParams, error)
	getUserDetailsService(ctx *gin.Context, username string) (*UserResponse, error)
	getAllUsersService(ctx *gin.Context) ([]UserResponse, error)
	loginUserService(ctx *gin.Context, req loginUserRequest) (*loginUSerResponse, error)
	logoutUsersService(ctx *gin.Context, username string, token string) error
	verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) error

	resendOTPService(ctx *gin.Context, username string) (*VerrifyEmailTxParams, error)
	updateUserService(ctx *gin.Context, username string, arg UpdateUserParams) (*UserResponse, error)
	updateUserImageService(ctx *gin.Context, username string, arg UpdateUserImageParams) error

	ForgotPasswordService(ctx *gin.Context, email string) error
	UpdatePasswordService(ctx *gin.Context, username string, arg UpdatePasswordParams) error
}

func (server *UserService) createUserService(ctx *gin.Context, req createUserRequest) (*VerrifyEmailTxParams, error) {
	var userID int64
	hashedPwd, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("cannot hash password: %v", err))
		return nil, fmt.Errorf("cannot hash password: %v", err)
	}
	var otp int64

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPwd,
		FullName:       req.FullName,
		Email:          req.Email,
		PhoneNumber:    pgtype.Text{String: req.PhoneNumber, Valid: true},
		Address:        pgtype.Text{String: req.Address, Valid: true},
		DataImage:      req.DataImage,
		OriginalImage:  pgtype.Text{String: req.OriginalImage, Valid: true},
		Role:           pgtype.Text{String: "user", Valid: true}, //
	}
	err = server.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {

		userID, err = server.storeDB.CreateUser(ctx, arg) // Check this line carefully

		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "unique_violation":
					ctx.JSON(http.StatusForbidden, "username or email already exists")
					return fmt.Errorf("username or email already exists")
				}
			} else {
				ctx.JSON(http.StatusInternalServerError, "internal server error")
				return fmt.Errorf("internal server error: %v", err)
			}
		}

		otp = util.RandomInt(100000, 999999)
		if err != nil {
			return fmt.Errorf("generate otp error: %v", err)
		}

		// Distribute the task to send a verification email
		payload := &worker.PayloadSendVerifyEmail{
			Username: req.Username,
			OTP:      otp,
		}

		go server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payload, asynq.Queue(worker.QueueDefault), asynq.MaxRetry(3))

		err = server.redis.StoreOTPInRedis(payload.Username, payload.OTP)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "failed to store otp in redis")
			return fmt.Errorf("failed to store otp in redis: %v", err)
		}

		return nil

	})

	if err != nil {
		// Delete the user if any part of the process fails
		if userID != 0 {
			deleteErr := server.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
				return q.DeleteUser(ctx, userID)
			})
			if deleteErr != nil {
				ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to delete user after error: %v", deleteErr))
				return nil, fmt.Errorf("failed to delete user after error: %w", deleteErr)
			}
		}

		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &VerrifyEmailTxParams{
		Username:   req.Username,
		SecretCode: otp,
	}, nil

}

func (server *UserService) getUserDetailsService(ctx *gin.Context, username string) (*UserResponse, error) {
	user, err := server.redis.UserInfoLoadCache(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}
	return &UserResponse{
		Username:      user.Username,
		FullName:      user.FullName,
		Role:          user.Role,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		Address:       user.Address,
		DataImage:     []byte(user.DataImage),
		OriginalImage: user.OriginalImage,
	}, nil
}

func (server *UserService) getAllUsersService(ctx *gin.Context) ([]UserResponse, error) {
	users, err := server.storeDB.GetAllUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	listUser := make([]UserResponse, 0)
	for _, u := range users {
		listUser = append(listUser, UserResponse{
			Username: u.Username,
			FullName: u.FullName,
			Email:    u.Email,
		})
	}
	return listUser, nil
}

func (service *UserService) loginUserService(ctx *gin.Context, req loginUserRequest) (*loginUSerResponse, error) {
	user, err := service.storeDB.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	_, err = service.redis.UserInfoLoadCache(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Incorrect passward")
		return nil, fmt.Errorf("Incorrect passward")
	}

	if !user.IsVerifiedEmail.Bool {
		ctx.JSON(http.StatusForbidden, "email not verified")
		return nil, fmt.Errorf("email not verified")
	}

	tokens, err := service.storeDB.InsertDeviceToken(ctx, db.InsertDeviceTokenParams{
		Username:   req.Username,
		Token:      req.Token,
		DeviceType: pgtype.Text{String: req.DeviceType, Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token device")
	}

	return &loginUSerResponse{
		User: UserResponse{
			Username:  user.Username,
			FullName:  user.FullName,
			Email:     user.Email,
			DataImage: []byte(user.DataImage),
		},
		DeviceToken: tokens.Token,
	}, nil
}

func (service *UserService) logoutUsersService(ctx *gin.Context, username string, token string) error {
	host, secure := util.SetCookieSameSite(ctx)
	ctx.SetCookie("refresh_token", "", -1, "/", host, secure, true)
	fmt.Println(username, token)
	err := service.storeDB.DeleteDeviceToken(ctx, db.DeleteDeviceTokenParams{
		Username: username,
		Token:    token,
	})
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
	service.redis.RemoveUserInfoCache(username)
	return nil
}

func (server *UserService) verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) error {

	err := server.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {

		storedOTP, err := server.redis.ReadOTPFromRedis(arg.Username)
		if err != nil {
			return fmt.Errorf("failed to verify OTP: %w", err)
		}

		if storedOTP != arg.SecretCode {
			return fmt.Errorf("invalid OTP")
		}

		// Delete OTP after successful verification
		otpKey := fmt.Sprintf("OTP-%s", arg.Username)
		if err := server.redis.DeleteOTPFromRedis(otpKey); err != nil {
			return fmt.Errorf("failed to delete OTP: %w", err)
		}
		_, err = q.VerifiedUser(ctx, arg.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "failed to verify user")
			return fmt.Errorf("failed to verify user: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}

func (service *UserService) resendOTPService(ctx *gin.Context, username string) (*VerrifyEmailTxParams, error) {
	otp := util.RandomInt(100000, 999999)
	if err := service.redis.StoreOTPInRedis(username, otp); err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to store otp in redis")
		return nil, fmt.Errorf("failed to store otp in redis: %w", err)
	}

	// Distribute the task to send a verification email
	payload := &worker.PayloadSendVerifyEmail{
		Username: username,
		OTP:      otp,
	}

	go service.taskDistributor.DistributeTaskSendVerifyEmail(ctx, payload, asynq.Queue(worker.QueueDefault), asynq.MaxRetry(3))

	return &VerrifyEmailTxParams{
		Username:   username,
		SecretCode: otp,
	}, nil
}

func (service *UserService) updateUserService(ctx *gin.Context, username string, arg UpdateUserParams) (*UserResponse, error) {

	var res db.User
	var req db.UpdateUserParams

	user, err := service.storeDB.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	fmt.Println(arg)

	if arg.FullName == "" {
		req.FullName = user.FullName
	} else {
		req.FullName = arg.FullName
	}

	if arg.Email == "" {
		req.Email = user.Email
	} else {
		req.Email = arg.Email
	}

	if arg.PhoneNumber == "" {
		req.PhoneNumber = user.PhoneNumber
	} else {
		req.PhoneNumber = pgtype.Text{String: arg.PhoneNumber, Valid: true}
	}

	if arg.Address == "" {
		req.Address = user.Address
	} else {
		req.Address = pgtype.Text{String: arg.Address, Valid: true}
	}

	req.Username = username

	err = service.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err = q.UpdateUser(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			return fmt.Errorf("internal server error: %v", err)
		}
		return nil
	})

	go service.redis.RemoveUserInfoCache(username)

	return &UserResponse{
		Username:      res.Username,
		FullName:      res.FullName,
		Email:         res.Email,
		PhoneNumber:   res.PhoneNumber.String,
		Address:       res.Address.String,
		OriginalImage: res.OriginalImage.String,
	}, nil
}

func (service *UserService) updateUserImageService(ctx *gin.Context, username string, arg UpdateUserImageParams) error {

	err := service.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		_, err := q.UpdateAvatarUser(ctx, db.UpdateAvatarUserParams{
			Username:      username,
			DataImage:     arg.DataImage,
			OriginalImage: pgtype.Text{String: arg.OriginalImage, Valid: true},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			return fmt.Errorf("internal server error: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update user image: %w", err)

	}

	// remove cache
	service.redis.RemoveUserInfoCache(username)

	return nil
}

func (s *UserService) ForgotPasswordService(ctx *gin.Context, email string) error {

	user, err := s.storeDB.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return fmt.Errorf("internal server error: %v", err)
	}

	// Generate a custom password
	customConfig := util.PasswordConfig{
		Length:        12,
		IncludeUpper:  true,
		IncludeLower:  true,
		IncludeNumber: true,
		IncludeSymbol: true, // No special characters
	}
	customPassword, err := util.GeneratePassword(customConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to generate password")
		return fmt.Errorf("failed to generate password: %w", err)
	}
	mailer := mail.NewGmailSender(s.config.EmailSenderName, s.config.EmailSenderAddress, s.config.EmailSenderPassword)
	if mailer == nil {
		log.Fatal("Failed to create mailer")
	}

	subject := "New Password"
	content := fmt.Sprintf("Your new password is: %s", customPassword)
	to := []string{email}

	// Send the new password to the user
	err = mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}
	hashedPwd, err := util.HashPassword(customPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to hash password")
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		_, err := q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
			Username:       user.Username,
			HashedPassword: hashedPwd,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			return fmt.Errorf("internal server error: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil

}

func (s *UserService) UpdatePasswordService(ctx *gin.Context, username string, arg UpdatePasswordParams) error {

	user, err := s.storeDB.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return fmt.Errorf("internal server error: %v", err)
	}
	// Check if the old password is correct
	err = util.CheckPassword(arg.OldPassword, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusForbidden, "incorrect old password")
		return fmt.Errorf("incorrect old password")
	}
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		newHashedPwd, err := util.HashPassword(arg.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "failed to hash password")
			return fmt.Errorf("failed to hash password: %w", err)
		}
		_, err = q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
			Username:       username,
			HashedPassword: newHashedPwd,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			return fmt.Errorf("internal server error: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}
	return nil
}
