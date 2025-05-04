package user

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type UserController struct {
	service UserServiceInterface
}

type UserService struct {
	storeDB         db.Store
	redis           *redis.ClientType
	taskDistributor worker.TaskDistributor
	config          util.Config
}

// route
type UserApi struct {
	controller UserControllerInterface
}

type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

type createUserRequest struct {
	Username        string `json:"username" binding:"required,alphanum"`
	Password        string `json:"password" binding:"required,min=6,max=25"`
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	PhoneNumber     string `json:"phone_number"`
	Address         string `json:"address"`
	Role            string `json:"role" binding:"required"`
	IsVerifiedEmail bool   `json:"is_verified_email"`
	DataImage       []byte `json:"-"`
	OriginalImage   string `json:"original_image"`
}

type UserResponse struct {
	UserID        int64     `json:"user_id"`
	Username      string    `json:"username"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	Role          string    `json:"role"`
	DataImage     []byte    `json:"data_image"`
	OriginalImage string    `json:"original_image"`
	RemovedAt     time.Time `json:"removed_at"`
}

type UpdateUserParams struct {
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}
type UpdateUserImageParams struct {
	DataImage     []byte `json:"data_image"`
	OriginalImage string `json:"original_image"`
}

type loginUserRequest struct {
	Username   string `json:"username" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required,min=6,max=25"`
	Token      string `json:"token"`
	DeviceType string `json:"device_type"`
}
type loginUSerResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
	DeviceToken           string       `json:"device_token"`
	DataImage             string       `json:"data_image"`
}

type InsertDoctorRequest struct {
	UserID            int64  `json:"user_id"`
	Specialization    string `json:"specialization"`
	YearsOfExperience int32  `json:"years_of_experience"`
	Education         string `json:"education"`
	CertificateNumber string `json:"certificate_number"`
	Bio               string `json:"bio"`
	ConsultationFee   int32  `json:"consultation_fee"`
}

type DoctorResponse struct {
	ID             int64  `json:"id"`
	Specialization string `json:"specialization"`
	Name           string `json:"name"`
	YearsOfExp     int32  `json:"years_of_exp"`
	Education      string `json:"education"`
	Certificate    string `json:"certificate"`
	Bio            string `json:"bio"`
}

type InsertDoctorScheduleRequest struct {
	DoctorID  int64  `json:"doctor_id"`
	Day       int32  `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	MaxAppoin int32  `json:"max_appointments"`
}

type DoctorScheduleResponse struct {
	ID              int64  `json:"id"`
	DoctorID        int64  `json:"doctor_id"`
	Day             int32  `json:"day"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	MaxAppointments int32  `json:"max_appointments"`
}

type VerrifyEmailTxParams struct {
	SecretCode int64  `json:"secret_code"`
	Username   string `json:"username"`
}

type VerrifyInput struct {
	SecretCode string `json:"secret_code"`
	Username   string `json:"username"`
}

type VerrifyEmailTxResult struct {
	User        db.User
	VerifyEmail db.VerifyEmail
}

type InsertTokenInfoRequest struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken pgtype.Text `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	Expiry       time.Time   `json:"expiry"`
}

// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPasswordResponse represents the response for forgot password
type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type UpdatePasswordParams struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}

type CreateStaffRequest struct {
	Username        string `json:"username" binding:"required,alphanum"`
	Password        string `json:"password" binding:"required,min=6,max=25"`
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	PhoneNumber     string `json:"phone_number"`
	Specialization  string `json:"specialization"`
	Address         string `json:"address"`
	Role            string `json:"role"`
	IsVerifiedEmail bool   `json:"is_verified_email"`
}

type CreateStaffResponse struct {
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}
