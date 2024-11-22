package user

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/util"
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/mail"
	"github.com/quanganh247-qa/go-blog-be/app/service/rabbitmq"
>>>>>>> 9d28896 (image pet)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
>>>>>>> 272832d (redis cache)
=======
>>>>>>> 6610455 (feat: redis queue)
)

type UserController struct {
	service UserServiceInterface
}

type UserService struct {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	storeDB         db.Store
	redis           *redis.ClientType
	taskDistributor worker.TaskDistributor
	config          util.Config
=======
	storeDB    db.Store
	emailQueue *rabbitmq.EmailQueue
	mailer     mail.EmailSender
>>>>>>> 9d28896 (image pet)
=======
	storeDB db.Store
	redis   *redis.ClientType
>>>>>>> 272832d (redis cache)
=======
	storeDB         db.Store
	redis           *redis.ClientType
	taskDistributor worker.TaskDistributor
>>>>>>> 6610455 (feat: redis queue)
}

// route
type UserApi struct {
	controller UserControllerInterface
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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 272832d (redis cache)
	Username      string    `json:"username"`
	FullName      string    `json:"full_name"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	Address       string    `json:"address"`
	Role          string    `json:"role"`
	DataImage     []byte    `json:"data_image"`
	OriginalImage string    `json:"original_image"`
	RemovedAt     time.Time `json:"removed_at"`
<<<<<<< HEAD
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
=======
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phone_number"`
	Address         string    `json:"address"`
	Role            string    `json:"role"`
	IsVerifiedEmail bool      `json:"is_verified_email"`
	DataImage       []byte    `json:"data_image"`
	OriginalImage   string    `json:"original_image"`
	RemovedAt       time.Time `json:"removed_at"`
>>>>>>> 0fb3f30 (user images)
=======
>>>>>>> 272832d (redis cache)
}

type UpdateUserParams struct {
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phone_number"`
	Address       string `json:"address"`
	DataImage     []byte `json:"data_image"`
	OriginalImage string `json:"original_image"`
}
type UpdateUserImageParams struct {
	DataImage     []byte `json:"data_image"`
	OriginalImage string `json:"original_image"`
}

type loginUserRequest struct {
<<<<<<< HEAD
<<<<<<< HEAD
	Username   string `json:"username" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required,min=6,max=25"`
=======
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=25"`
	// Token      string `json:"token" binding:"required"`
>>>>>>> c3c833d (login api)
=======
	Username   string `json:"username" binding:"required,alphanum"`
	Password   string `json:"password" binding:"required,min=6,max=25"`
>>>>>>> 290baeb (fixed vaccine routes)
	Token      string `json:"token"`
	DeviceType string `json:"device_type"`
}
type loginUSerResponse struct {
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
<<<<<<< HEAD
<<<<<<< HEAD
	DeviceToken           string       `json:"device_token"`
=======
	DeviceToken           []string     `json:"device_token"`
>>>>>>> 0fb3f30 (user images)
=======
	DeviceToken           string       `json:"device_token"`
>>>>>>> c3c833d (login api)
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
<<<<<<< HEAD
<<<<<<< HEAD
	SecretCode int64  `json:"secret_code"`
	Username   string `json:"username"`
}

type VerrifyInput struct {
	SecretCode string `json:"secret_code"`
	Username   string `json:"username"`
<<<<<<< HEAD
=======
	EmailId    int64  `json:"email_id"`
=======
>>>>>>> edfe5ad (OTP verifycation)
	SecretCode int64  `json:"secret_code"`
	Username   string `json:"username"`
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> 290baeb (fixed vaccine routes)
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
