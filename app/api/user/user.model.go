package user

import (
	"time"

	"github.com/gofrs/uuid"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/rabbitmq"
)

type UserController struct {
	service UserServiceInterface
}

type UserService struct {
	storeDB    db.Store
	emailQueue *rabbitmq.EmailQueue
}
type UserApi struct {
	controller UserControllerInterface
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=25"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=25"`
}
type loginUSerResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type VerrifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerrifyEmailTxResult struct {
	User        db.User
	VerifyEmail db.VerifyEmail
}
