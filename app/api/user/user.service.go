package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/rabbitmq"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type UserServiceInterface interface {
	createUserService(ctx *gin.Context, req createUserRequest) (*db.User, error)
	getAllUsersService(ctx *gin.Context) ([]UserResponse, error)
	loginUserService(ctx *gin.Context, req loginUserRequest) error
}

func (server *UserService) createUserService(ctx *gin.Context, req createUserRequest) (*db.User, error) {
	if server.storeDB == nil {
		log.Println("storeDB is nil!")
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("storeDB is nil")
	}

	hashedPwd, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("cannot hash password: %v", err))
		return nil, fmt.Errorf("cannot hash password: %v", err)
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPwd,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.storeDB.CreateUser(ctx, arg) // Check this line carefully

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, "username or email already exists")
				return nil, fmt.Errorf("username or email already exists")
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			return nil, fmt.Errorf("internal server error: %v", err)
		}
	}

	// **Send Email using RabbitMQ**
	emailPayload := rabbitmq.PayloadVerifyEmail{
		Username: arg.Username,
	}

	err = server.emailQueue.PublishEmail(emailPayload)

	if err != nil {
		log.Println("Error publishing email:", err)
		ctx.JSON(http.StatusInternalServerError, "failed to send verification email")
		return nil, fmt.Errorf("failed to send verification email: %v", err)
	}

	return &db.User{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
	// return nil, nil
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
			Username:          u.Username,
			FullName:          u.FullName,
			Email:             u.Email,
			PasswordChangedAt: u.PasswordChangedAt,
			CreatedAt:         u.CreatedAt,
		})
	}
	return listUser, nil
}

func (service *UserService) loginUserService(ctx *gin.Context, req loginUserRequest) error {
	_, err := service.storeDB.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return fmt.Errorf("internal server error: %v", err)
	}
	return nil
}

func (server *UserService) VerifyEmail(ctx context.Context, arg VerrifyEmailTxParams) error {

	var result VerrifyEmailTxResult

	err := server.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, db.UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err = q.UpdateUser(ctx, db.UpdateUserParams{
			Username: result.User.Username,
			IsVerifiedEmail: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		return err
	})
	if err != nil {
		return fmt.Errorf("Transaction verify email")
	}

	return nil
}
