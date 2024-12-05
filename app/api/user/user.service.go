package user

import (
	"database/sql"
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"log"
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> 272832d (redis cache)
=======
	"log"
>>>>>>> 1a9e82a (reset password api)
=======
	"log"
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 272832d (redis cache)
=======
	"log"
>>>>>>> 1a9e82a (reset password api)
	"math/big"
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> ae87825 (updated)
	"net/http"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 1a9e82a (reset password api)
=======
>>>>>>> 1a9e82a (reset password api)
	"github.com/quanganh247-qa/go-blog-be/app/service/mail"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/rabbitmq"
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 272832d (redis cache)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
>>>>>>> 6610455 (feat: redis queue)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/rabbitmq"
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 272832d (redis cache)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
>>>>>>> 6610455 (feat: redis queue)
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type UserServiceInterface interface {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	createUserService(ctx *gin.Context, req createUserRequest) (*VerrifyEmailTxParams, error)
=======
	// createUserService(ctx *gin.Context, req createUserRequest) (*db.User, error)
=======
>>>>>>> 272832d (redis cache)
	createUserService(ctx *gin.Context, req createUserRequest) error
>>>>>>> 0fb3f30 (user images)
=======
	createUserService(ctx *gin.Context, req createUserRequest) (*VerrifyEmailTxParams, error)
>>>>>>> edfe5ad (OTP verifycation)
=======
	// createUserService(ctx *gin.Context, req createUserRequest) (*db.User, error)
=======
>>>>>>> 272832d (redis cache)
	createUserService(ctx *gin.Context, req createUserRequest) error
>>>>>>> 0fb3f30 (user images)
=======
	createUserService(ctx *gin.Context, req createUserRequest) (*VerrifyEmailTxParams, error)
>>>>>>> edfe5ad (OTP verifycation)
	getUserDetailsService(ctx *gin.Context, username string) (*UserResponse, error)
	getAllUsersService(ctx *gin.Context) ([]UserResponse, error)
	loginUserService(ctx *gin.Context, req loginUserRequest) (*loginUSerResponse, error)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	logoutUsersService(ctx *gin.Context, username string, token string) error
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) error

<<<<<<< HEAD
	resendOTPService(ctx *gin.Context, username string) (*VerrifyEmailTxParams, error)
	updateUserService(ctx *gin.Context, username string, arg UpdateUserParams) (*UserResponse, error)
	updateUserImageService(ctx *gin.Context, username string, arg UpdateUserImageParams) error

	ForgotPasswordService(ctx *gin.Context, email string) error
	UpdatePasswordService(ctx *gin.Context, username string, arg UpdatePasswordParams) error
	GetAllRoleService(ctx *gin.Context) ([]string, error)
=======
	logoutUsersService(ctx *gin.Context, username string) error
=======
	logoutUsersService(ctx *gin.Context, username string, token string) error
>>>>>>> 9d28896 (image pet)
=======
	logoutUsersService(ctx *gin.Context, username string, token string) error
>>>>>>> 9d28896 (image pet)
	verifyEmailService(ctx *gin.Context, req VerrifyEmailTxParams) (VerrifyEmailTxResult, error)
=======
	verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) (VerrifyEmailTxResult, error)
>>>>>>> 6610455 (feat: redis queue)
=======
	verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) error
>>>>>>> edfe5ad (OTP verifycation)
=======
	verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) (VerrifyEmailTxResult, error)
>>>>>>> 6610455 (feat: redis queue)
=======
	verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) error
>>>>>>> edfe5ad (OTP verifycation)
	createDoctorService(ctx *gin.Context, arg InsertDoctorRequest, username string) (*DoctorResponse, error)
	createDoctorScheduleService(ctx *gin.Context, arg InsertDoctorScheduleRequest, username string) (*DoctorScheduleResponse, error)
	getDoctorByID(ctx *gin.Context, userID int64) (*DoctorResponse, error)
	insertTimeSlots(ctx *gin.Context, username string, arg db.InsertTimeslotParams) (*db.Timeslot, error)
	GetTimeslotsAvailable(ctx *gin.Context, doctorID int64, date string) ([]db.GetTimeslotsAvailableRow, error)
	GetAllTimeslots(ctx *gin.Context, doctorID int64, date string) ([]db.GetTimeslotsAvailableRow, error)
	UpdateDoctorAvailable(ctx *gin.Context, time_slot_id int64) error
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	// InsertTokenInfoService(ctx *gin.Context, arg InsertTokenInfoRequest, username string) (*db.TokenInfo, error)
>>>>>>> 79a3bcc (medicine api)
=======
	// InsertTokenInfoService(ctx *gin.Context, arg InsertTokenInfoRequest, username string) (*db.TokenInfo, error)
>>>>>>> 79a3bcc (medicine api)
}

<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> ae87825 (updated)
	resendOTPService(ctx *gin.Context, username string) (*VerrifyEmailTxParams, error)
	updateUserService(ctx *gin.Context, username string, arg UpdateUserParams) (*UserResponse, error)
	updateUserImageService(ctx *gin.Context, username string, arg UpdateUserImageParams) error
=======
func (server *UserService) createUserService(ctx *gin.Context, req createUserRequest) error {
<<<<<<< HEAD
>>>>>>> 0fb3f30 (user images)

	ForgotPasswordService(ctx *gin.Context, email string) error
	UpdatePasswordService(ctx *gin.Context, username string, arg UpdatePasswordParams) error
	GetAllRoleService(ctx *gin.Context) ([]string, error)
}

>>>>>>> edfe5ad (OTP verifycation)
=======
	resendOTPService(ctx *gin.Context, username string) (*VerrifyEmailTxParams, error)
	updateUserService(ctx *gin.Context, username string, arg UpdateUserParams) (*UserResponse, error)
	updateUserImageService(ctx *gin.Context, username string, arg UpdateUserImageParams) error
	GetDoctorsService(ctx *gin.Context) ([]DoctorResponse, error)
	ForgotPasswordService(ctx *gin.Context, email string) error
	UpdatePasswordService(ctx *gin.Context, username string, arg UpdatePasswordParams) error
}

>>>>>>> edfe5ad (OTP verifycation)
func (server *UserService) createUserService(ctx *gin.Context, req createUserRequest) (*VerrifyEmailTxParams, error) {
	var userID int64
=======
func (server *UserService) createUserService(ctx *gin.Context, req createUserRequest) error {
<<<<<<< HEAD

>>>>>>> 0fb3f30 (user images)
=======
	var userID int64
>>>>>>> 1f24c18 (feat: OTP with redis)
=======
	var userID int64
>>>>>>> 1f24c18 (feat: OTP with redis)
	hashedPwd, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("cannot hash password: %v", err))
		return nil, fmt.Errorf("cannot hash password: %v", err)
	}
	var otp int64

	arg := db.CreateUserParams{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> 6610455 (feat: redis queue)
		Username:       req.Username,
		HashedPassword: hashedPwd,
		FullName:       req.FullName,
		Email:          req.Email,
		PhoneNumber:    pgtype.Text{String: req.PhoneNumber, Valid: true},
		Address:        pgtype.Text{String: req.Address, Valid: true},
		DataImage:      req.DataImage,
		OriginalImage:  pgtype.Text{String: req.OriginalImage, Valid: true},
		Role:           pgtype.Text{String: "user", Valid: true}, //
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> 0fb3f30 (user images)
		Username:        req.Username,
		HashedPassword:  hashedPwd,
		FullName:        req.FullName,
		Email:           req.Email,
		PhoneNumber:     pgtype.Text{String: req.PhoneNumber, Valid: true},
		Address:         pgtype.Text{String: req.Address, Valid: true},
		DataImage:       req.DataImage,
<<<<<<< HEAD
<<<<<<< HEAD
		OriginalImage:   pgtype.Text{String: req.OriginalImage, Valid: true},
		Role:            pgtype.Text{String: "user", Valid: true}, //
		IsVerifiedEmail: pgtype.Bool{Bool: true, Valid: true},
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> 6610455 (feat: redis queue)
=======
		OriginalImage:   req.OriginalImage,
=======
		OriginalImage:   pgtype.Text{String: req.OriginalImage, Valid: true},
>>>>>>> 272832d (redis cache)
		Role:            pgtype.Text{String: "user", Valid: true}, //
		IsVerifiedEmail: pgtype.Bool{Bool: true, Valid: true},
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> 6610455 (feat: redis queue)
	}
<<<<<<< HEAD
	err = server.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> ada3717 (Docker file)
		_, err := server.storeDB.CreateUser(ctx, arg) // Check this line carefully

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

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 290baeb (fixed vaccine routes)
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
=======
	_, err = server.storeDB.CreateUser(ctx, arg) // Check this line carefully

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, "username or email already exists")
				return fmt.Errorf("username or email already exists")
>>>>>>> 0fb3f30 (user images)
			}
=======
	_, err = server.storeDB.CreateUser(ctx, arg) // Check this line carefully
=======
	err = server.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
>>>>>>> 6610455 (feat: redis queue)

		_, err = server.storeDB.CreateUser(ctx, arg) // Check this line carefully
=======
		userID, err = server.storeDB.CreateUser(ctx, arg) // Check this line carefully
>>>>>>> 1f24c18 (feat: OTP with redis)

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
<<<<<<< HEAD
		} else {
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			return fmt.Errorf("internal server error: %v", err)
<<<<<<< HEAD
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> 0fb3f30 (user images)
		}

		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
		return nil, fmt.Errorf("failed to create user: %w", err)
=======
		}

		otp := util.RandomInt(1000000, 9999999)
=======
		otp = util.RandomInt(1000000, 9999999)
>>>>>>> edfe5ad (OTP verifycation)
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
=======
=======
>>>>>>> 1f24c18 (feat: OTP with redis)
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
<<<<<<< HEAD
			ctx.JSON(http.StatusInternalServerError, "failed to enqueue task")
			return fmt.Errorf("failed to enqueue task: %v", err)
>>>>>>> 6610455 (feat: redis queue)
=======
			ctx.JSON(http.StatusInternalServerError, "failed to store otp in redis")
			return fmt.Errorf("failed to store otp in redis: %v", err)
>>>>>>> 1f24c18 (feat: OTP with redis)
		}

		return nil

	})
<<<<<<< HEAD
<<<<<<< HEAD

	if err != nil {
<<<<<<< HEAD
		return fmt.Errorf("failed to create user: %w", err)
>>>>>>> 6610455 (feat: redis queue)
=======
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
<<<<<<< HEAD
		return err
>>>>>>> 1f24c18 (feat: OTP with redis)
=======
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
>>>>>>> 6610455 (feat: redis queue)
=======

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
<<<<<<< HEAD
		return err
>>>>>>> 1f24c18 (feat: OTP with redis)
	}

<<<<<<< HEAD
<<<<<<< HEAD
=======
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

>>>>>>> edfe5ad (OTP verifycation)
=======
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

>>>>>>> edfe5ad (OTP verifycation)
	return &VerrifyEmailTxParams{
		Username:   req.Username,
		SecretCode: otp,
	}, nil
<<<<<<< HEAD
<<<<<<< HEAD

}

<<<<<<< HEAD
func (server *UserService) getUserDetailsService(ctx *gin.Context, username string) (*UserResponse, error) {
	user, err := server.redis.UserInfoLoadCache(username)
=======
	// if err != nil {
	// 	log.Println("Error publishing email:", err)
	// 	ctx.JSON(http.StatusInternalServerError, "failed to send verification email")
	// 	return nil, fmt.Errorf("failed to send verification email: %v", err)
	// }

<<<<<<< HEAD
=======
>>>>>>> 9d28896 (image pet)
	return nil
=======
>>>>>>> edfe5ad (OTP verifycation)

}

func (server *UserService) getUserDetailsService(ctx *gin.Context, username string) (*UserResponse, error) {
<<<<<<< HEAD
	user, err := server.storeDB.GetUser(ctx, username)
>>>>>>> 0fb3f30 (user images)
=======
	user, err := server.redis.UserInfoLoadCache(username)
>>>>>>> 272832d (redis cache)
=======
=======
>>>>>>> 9d28896 (image pet)
	return nil
=======
>>>>>>> edfe5ad (OTP verifycation)

}

func (server *UserService) getUserDetailsService(ctx *gin.Context, username string) (*UserResponse, error) {
<<<<<<< HEAD
	user, err := server.storeDB.GetUser(ctx, username)
>>>>>>> 0fb3f30 (user images)
=======
	user, err := server.redis.UserInfoLoadCache(username)
>>>>>>> 272832d (redis cache)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 290baeb (fixed vaccine routes)
=======
>>>>>>> 290baeb (fixed vaccine routes)
		Role:          user.Role,
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		Address:       user.Address,
		DataImage:     []byte(user.DataImage),
=======
		Email:         user.Email,
<<<<<<< HEAD
		PhoneNumber:   user.PhoneNumber.String,
		Address:       user.Address.String,
		DataImage:     user.DataImage,
>>>>>>> 0fb3f30 (user images)
=======
		PhoneNumber:   user.PhoneNumber,
		Address:       user.Address,
		DataImage:     []byte(user.DataImage),
>>>>>>> 272832d (redis cache)
=======
		Email:         user.Email,
<<<<<<< HEAD
		PhoneNumber:   user.PhoneNumber.String,
		Address:       user.Address.String,
		DataImage:     user.DataImage,
>>>>>>> 0fb3f30 (user images)
=======
		PhoneNumber:   user.PhoneNumber,
		Address:       user.Address,
		DataImage:     []byte(user.DataImage),
>>>>>>> 272832d (redis cache)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> b393bb9 (add service and add permission)
	_, err = service.redis.UserInfoLoadCache(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
		}
<<<<<<< HEAD
=======
	device_tokens, err := service.storeDB.GetDeviceTokenByUsername(ctx, req.Username)
	if err != nil {
>>>>>>> 0fb3f30 (user images)
=======
	device_tokens, err := service.storeDB.GetDeviceTokenByUsername(ctx, req.Username)
	if err != nil {
>>>>>>> 0fb3f30 (user images)
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> 21608b5 (cart and order api)
	fmt.Println(user.Username)

	// user, err := service.redis.UserInfoLoadCache(req.Username)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		ctx.JSON(http.StatusNotFound, "user not found")
	// 		return nil, fmt.Errorf("user not found")
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, "internal server error")
	// 	return nil, fmt.Errorf("internal server error: %v", err)
	// }
<<<<<<< HEAD
=======
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}
>>>>>>> b393bb9 (add service and add permission)

>>>>>>> 21608b5 (cart and order api)
=======

>>>>>>> 21608b5 (cart and order api)
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Incorrect passward")
		return nil, fmt.Errorf("Incorrect passward")
	}

	fmt.Println(user.IsVerifiedEmail.Bool)

	if !user.IsVerifiedEmail.Bool {
		ctx.JSON(http.StatusForbidden, "email not verified")
		return nil, fmt.Errorf("email not verified")
	}
=======

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Incorrect passward")
		return nil, fmt.Errorf("Incorrect passward")
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Incorrect passward")
		return nil, fmt.Errorf("Incorrect passward")
	}

>>>>>>> c3c833d (login api)
	// device_tokens, err := service.storeDB.GetDeviceTokenByUsername(ctx, req.Username)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, "internal server error")
	// 	return nil, fmt.Errorf("internal server error: %v", err)
	// }

	// var device_tokens_response []string
	// for _, d := range device_tokens {
	// 	device_tokens_response = append(device_tokens_response, d.Token)
	// }
<<<<<<< HEAD
>>>>>>> c3c833d (login api)

=======
>>>>>>> 6610455 (feat: redis queue)
=======
=======
>>>>>>> edfe5ad (OTP verifycation)
	if !user.IsVerifiedEmail.Bool {
		ctx.JSON(http.StatusForbidden, "email not verified")
		return nil, fmt.Errorf("email not verified")
	}

<<<<<<< HEAD
>>>>>>> edfe5ad (OTP verifycation)
=======

>>>>>>> c3c833d (login api)
=======
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> edfe5ad (OTP verifycation)
	tokens, err := service.storeDB.InsertDeviceToken(ctx, db.InsertDeviceTokenParams{
		Username:   req.Username,
		Token:      req.Token,
		DeviceType: pgtype.Text{String: req.DeviceType, Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token device")
<<<<<<< HEAD
<<<<<<< HEAD
=======
	var device_tokens_response []string
	for _, d := range device_tokens {
		device_tokens_response = append(device_tokens_response, d.Token)
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> c3c833d (login api)
=======
	var device_tokens_response []string
	for _, d := range device_tokens {
		device_tokens_response = append(device_tokens_response, d.Token)
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> c3c833d (login api)
	}

	return &loginUSerResponse{
		User: UserResponse{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
			Username:  user.Username,
			FullName:  user.FullName,
			Email:     user.Email,
<<<<<<< HEAD
			DataImage: []byte(user.DataImage),
		},
<<<<<<< HEAD
		DeviceToken: tokens.Token,
=======
		DeviceToken: device_tokens_response,
>>>>>>> 0fb3f30 (user images)
=======
			DataImage: user.DataImage,
=======
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
>>>>>>> c449ffc (feat: cart api)
=======
=======
>>>>>>> 30f20e5 (dât image)
			Username:  user.Username,
			FullName:  user.FullName,
			Email:     user.Email,
			DataImage: []byte(user.DataImage),
<<<<<<< HEAD
>>>>>>> 30f20e5 (dât image)
		},
<<<<<<< HEAD
		DeviceToken: tokens.Token,
>>>>>>> c3c833d (login api)
=======
		DeviceToken: device_tokens_response,
>>>>>>> 0fb3f30 (user images)
=======
			Username:  user.Username,
			FullName:  user.FullName,
			Email:     user.Email,
			DataImage: user.DataImage,
=======
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
>>>>>>> c449ffc (feat: cart api)
=======
>>>>>>> 30f20e5 (dât image)
		},
		DeviceToken: tokens.Token,
>>>>>>> c3c833d (login api)
	}, nil
}

func (service *UserService) logoutUsersService(ctx *gin.Context, username string, token string) error {
	host, secure := util.SetCookieSameSite(ctx)
	ctx.SetCookie("refresh_token", "", -1, "/", host, secure, true)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	fmt.Println(username, token)
=======
>>>>>>> 9d28896 (image pet)
=======
	fmt.Println(username, token)
>>>>>>> 8d5618d (feat: update logout)
=======
>>>>>>> 9d28896 (image pet)
=======
	fmt.Println(username, token)
>>>>>>> 8d5618d (feat: update logout)
	err := service.storeDB.DeleteDeviceToken(ctx, db.DeleteDeviceTokenParams{
		Username: username,
		Token:    token,
	})
	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	service.redis.RemoveUserInfoCache(username)
=======
=======
>>>>>>> 9d28896 (image pet)
	// service.redis.RemoveUserInfoCache(username)
>>>>>>> 9d28896 (image pet)
=======
	service.redis.RemoveUserInfoCache(username)
>>>>>>> 272832d (redis cache)
=======
	service.redis.RemoveUserInfoCache(username)
>>>>>>> 272832d (redis cache)
	return nil
}

func (server *UserService) verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) error {

	err := server.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		storedOTP, err := server.redis.ReadOTPFromRedis(arg.Username)
=======
=======
>>>>>>> 1f24c18 (feat: OTP with redis)
		err = server.readOTPFromRedis(ctx, arg.Username)
		if err != nil {

			return fmt.Errorf("incorrect otp")
		}
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, db.UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
>>>>>>> 1f24c18 (feat: OTP with redis)
		if err != nil {
<<<<<<< HEAD
<<<<<<< HEAD
			return fmt.Errorf("failed to verify OTP: %w", err)
=======
=======
>>>>>>> 6610455 (feat: redis queue)

			return err
>>>>>>> 6610455 (feat: redis queue)
		}

<<<<<<< HEAD
<<<<<<< HEAD
		if storedOTP != arg.SecretCode {
			return fmt.Errorf("invalid OTP")
=======
=======
>>>>>>> 6610455 (feat: redis queue)
		result.User, err = q.VerifiedUser(ctx, db.VerifiedUserParams{
			Username: result.VerifyEmail.Username,
			IsVerifiedEmail: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
<<<<<<< HEAD
=======
		storedOTP, err := server.redis.ReadOTPFromRedis(arg.Username)
		if err != nil {
			return fmt.Errorf("failed to verify OTP: %w", err)
		}

		if storedOTP != arg.SecretCode {
			return fmt.Errorf("invalid OTP")
		}
>>>>>>> edfe5ad (OTP verifycation)

		// Delete OTP after successful verification
		otpKey := fmt.Sprintf("OTP-%s", arg.Username)
		if err := server.redis.DeleteOTPFromRedis(otpKey); err != nil {
			return fmt.Errorf("failed to delete OTP: %w", err)
		}
		_, err = q.VerifiedUser(ctx, arg.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, "failed to verify user")
			return fmt.Errorf("failed to verify user: %w", err)
=======
=======
		storedOTP, err := server.redis.ReadOTPFromRedis(arg.Username)
		if err != nil {
			return fmt.Errorf("failed to verify OTP: %w", err)
		}

		if storedOTP != arg.SecretCode {
			return fmt.Errorf("invalid OTP")
		}
>>>>>>> edfe5ad (OTP verifycation)

		// Delete OTP after successful verification
		otpKey := fmt.Sprintf("OTP-%s", arg.Username)
		if err := server.redis.DeleteOTPFromRedis(otpKey); err != nil {
			return fmt.Errorf("failed to delete OTP: %w", err)
		}
		_, err = q.VerifiedUser(ctx, arg.Username)
		if err != nil {
<<<<<<< HEAD

			return err
>>>>>>> 6610455 (feat: redis queue)
=======
			ctx.JSON(http.StatusInternalServerError, "failed to verify user")
			return fmt.Errorf("failed to verify user: %w", err)
>>>>>>> edfe5ad (OTP verifycation)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
<<<<<<< HEAD
=======
}

func (service *UserService) resendOTPService(ctx *gin.Context, username string) (*VerrifyEmailTxParams, error) {
	otp := util.RandomInt(1000000, 9999999)
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
>>>>>>> edfe5ad (OTP verifycation)
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
	go service.redis.RemoveUserInfoCache(username)

	return nil
}

<<<<<<< HEAD
func (server *UserService) createDoctorService(ctx *gin.Context, arg InsertDoctorRequest, username string) (*DoctorResponse, error) {

	user, err := server.storeDB.GetUser(ctx, username)
	fmt.Println(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
>>>>>>> 6610455 (feat: redis queue)
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
<<<<<<< HEAD
=======

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

func (server *UserService) createDoctorService(ctx *gin.Context, arg InsertDoctorRequest, username string) (*DoctorResponse, error) {
>>>>>>> 473cd1d (uplaod image method)

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

<<<<<<< HEAD
<<<<<<< HEAD
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
	go service.redis.RemoveUserInfoCache(username)

	return nil
}

=======
>>>>>>> ae87825 (updated)
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

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
// update password
>>>>>>> a2c21c8 (update pass)
=======
>>>>>>> ae87825 (updated)
=======
// update password
>>>>>>> a2c21c8 (update pass)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 3980627 (generated test cases with keploy)
=======
>>>>>>> 3980627 (generated test cases with keploy)
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

func (s *UserService) GetAllRoleService(ctx *gin.Context) ([]string, error) {
	roles, err := s.storeDB.GetAllRole(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all role: %w", err)
	}
	roleList := make([]string, 0)
	for _, role := range roles {
		roleList = append(roleList, role.String)
	}
	return roleList, nil
}
<<<<<<< HEAD
=======
	return nil // Successfully updated
}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD

<<<<<<< HEAD
=======
	return nil // Successfully updated
}
<<<<<<< HEAD

<<<<<<< HEAD
>>>>>>> 79a3bcc (medicine api)
// func (s *UserService) InsertTokenInfoService(ctx *gin.Context, arg InsertTokenInfoRequest, username string) (*db.TokenInfo, error) {
// 	tokenInfo, err := s.storeDB.InsertTokenInfo(ctx, db.InsertTokenInfoParams{
// 		AccessToken:  arg.AccessToken,
// 		TokenType:    arg.TokenType,
// 		UserName:     username,
// 		RefreshToken: pgtype.Text{String: arg.RefreshToken.String, Valid: true},
// 		Expiry:       arg.Expiry,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to insert token info: %w", err)
// 	}
// 	return &db.TokenInfo{
// 		AccessToken:  tokenInfo.AccessToken,
// 		TokenType:    tokenInfo.TokenType,
// 		RefreshToken: tokenInfo.RefreshToken,
// 		Expiry:       tokenInfo.Expiry,
// 	}, nil
// }
<<<<<<< HEAD
>>>>>>> 79a3bcc (medicine api)
=======
=======
>>>>>>> 9d28896 (image pet)
func (s *UserService) ProccessTaskSendVerifyEmail(ctx context.Context, payload rabbitmq.PayloadVerifyEmail) error {
	// var payload PayloadVerifyEmail
	// if err := json.Unmarshal(task.Payload(), &payload); err != nil {
	// 	return fmt.Errorf("failed to unmarshal payload: %w", err)
	// }
	log.Printf("Processing task for user: %s", payload.Username)

	user, err := s.storeDB.GetUser(ctx, payload.Username)
	if err != nil {

		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exists: %w", err)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}
	log.Printf("User retrieved successfully")

	verifyEmail, err := s.storeDB.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email %w", err)
	}
	subject := "Welcome to Simple Bank"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:8088/api/v1/user/verify-email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}
	fmt.Println(subject)
	err = s.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	return nil
}
<<<<<<< HEAD
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 272832d (redis cache)
=======
=======
>>>>>>> e30b070 (Get list appoinment by user)

func (s *UserService) GetDoctorsService(ctx *gin.Context) ([]DoctorResponse, error) {
	// offset := (pagination.Page - 1) * pagination.PageSize

	rows, err := s.storeDB.GetDoctors(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get doctors")
	}
	var doctors []DoctorResponse
	for _, row := range rows {
		doctors = append(doctors, DoctorResponse{
			ID:             row.DoctorID,
			Specialization: row.Specialization.String,
			Name:           row.FullName,
			YearsOfExp:     row.YearsOfExperience.Int32,
			Education:      row.Education.String,
			Certificate:    row.CertificateNumber.String,
			Bio:            row.Bio.String,
		})
	}

	return doctors, nil

}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> e30b070 (Get list appoinment by user)
=======
=======
>>>>>>> 1a9e82a (reset password api)

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
<<<<<<< HEAD
=======
		_, err := q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
			Username:       username,
			HashedPassword: arg.Password,
>>>>>>> a2c21c8 (update pass)
=======
>>>>>>> 1a9e82a (reset password api)
=======
		_, err := q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
			Username:       username,
			HashedPassword: arg.Password,
>>>>>>> a2c21c8 (update pass)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 1a9e82a (reset password api)

	return nil

}
<<<<<<< HEAD
>>>>>>> 1a9e82a (reset password api)
=======
	return nil
}
>>>>>>> a2c21c8 (update pass)
=======
>>>>>>> 4ccd381 (Update appointment flow)
=======
>>>>>>> 79a3bcc (medicine api)
=======
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> 272832d (redis cache)
=======
>>>>>>> e30b070 (Get list appoinment by user)
=======
>>>>>>> 1a9e82a (reset password api)
=======
	return nil
}
>>>>>>> a2c21c8 (update pass)
