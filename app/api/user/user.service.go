package user

import (
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type UserServiceInterface interface {
	createUserService(ctx *gin.Context, req createUserRequest) (*db.User, error)
	getAllUsersService(ctx *gin.Context) ([]UserResponse, error)
	loginUserService(ctx *gin.Context, req loginUserRequest) error
	verifyEmailService(ctx *gin.Context, req VerrifyEmailTxParams) (VerrifyEmailTxResult, error)
	createDoctorService(ctx *gin.Context, arg InsertDoctorRequest, username string) (*DoctorResponse, error)
	createDoctorScheduleService(ctx *gin.Context, arg InsertDoctorScheduleRequest, username string) (*DoctorScheduleResponse, error)
	getDoctorByID(ctx *gin.Context, userID int64) (*DoctorResponse, error)
}

func (server *UserService) createUserService(ctx *gin.Context, req createUserRequest) (*db.User, error) {

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

	// // **Send Email using RabbitMQ**
	// emailPayload := rabbitmq.PayloadVerifyEmail{
	// 	Username: arg.Username,
	// }

	// err = server.emailQueue.PublishEmail(emailPayload)

	// if err != nil {
	// 	log.Println("Error publishing email:", err)
	// 	ctx.JSON(http.StatusInternalServerError, "failed to send verification email")
	// 	return nil, fmt.Errorf("failed to send verification email: %v", err)
	// }

	return &db.User{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
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
			Username:  u.Username,
			FullName:  u.FullName,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
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

func (server *UserService) verifyEmailService(ctx *gin.Context, arg VerrifyEmailTxParams) (VerrifyEmailTxResult, error) {

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
		return VerrifyEmailTxResult{}, fmt.Errorf("failed to verify email: %w", err)
	}

	return result, nil
}

func (server *UserService) createDoctorService(ctx *gin.Context, arg InsertDoctorRequest, username string) (*DoctorResponse, error) {

	user, err := server.storeDB.GetUser(ctx, username)
	fmt.Println(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}
	doctor, err := server.storeDB.InsertDoctor(ctx, db.InsertDoctorParams{
		UserID:            user.ID,
		Specialization:    pgtype.Text{String: arg.Specialization, Valid: true},
		YearsOfExperience: pgtype.Int4{Int32: arg.YearsOfExperience, Valid: true},
		Education:         pgtype.Text{String: arg.Education, Valid: true},
		CertificateNumber: pgtype.Text{String: arg.CertificateNumber, Valid: true},
		Bio:               pgtype.Text{String: arg.Bio, Valid: true},
		ConsultationFee:   pgtype.Numeric{Int: big.NewInt(int64(arg.ConsultationFee)), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create doctor: %w", err)
	}

	return &DoctorResponse{
		ID:             doctor.ID,
		Specialization: doctor.Specialization.String,
		Name:           user.FullName,
		YearsOfExp:     doctor.YearsOfExperience.Int32,
		Education:      doctor.Education.String,
		Certificate:    doctor.CertificateNumber.String,
		Bio:            doctor.Bio.String,
	}, nil
}

func (s *UserService) createDoctorScheduleService(ctx *gin.Context, arg InsertDoctorScheduleRequest, username string) (*DoctorScheduleResponse, error) {
	user, err := s.storeDB.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}
	doctor, err := s.storeDB.GetDoctor(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "doctor not found")
			return nil, fmt.Errorf("doctor not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	// string to pgtype.Time
	startTime, err := time.Parse("15:04:05", arg.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid start time")
		return nil, fmt.Errorf("invalid start time: %w", err)
	}
	endTime, err := time.Parse("15:04:05", arg.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid end time")
		return nil, fmt.Errorf("invalid end time: %w", err)
	} // time.Time to pgtype.Time
	pgStartTime := pgtype.Time{
		Microseconds: int64(startTime.Hour()*3600+startTime.Minute()*60+startTime.Second()) * 1e6,
		Valid:        true,
	}
	pgEndTime := pgtype.Time{
		Microseconds: int64(endTime.Hour()*3600+endTime.Minute()*60+endTime.Second()) * 1e6,
		Valid:        true,
	}

	doctorSchedule, err := s.storeDB.InsertDoctorSchedule(ctx, db.InsertDoctorScheduleParams{
		DoctorID:        doctor.ID,
		DayOfWeek:       pgtype.Int4{Int32: arg.Day, Valid: true},
		MaxAppointments: pgtype.Int4{Int32: arg.MaxAppoin, Valid: true},
		StartTime:       pgStartTime,
		EndTime:         pgEndTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create doctor schedule: %w", err)
	}

	return &DoctorScheduleResponse{
		ID:              doctorSchedule.ID,
		DoctorID:        doctorSchedule.DoctorID,
		Day:             doctorSchedule.DayOfWeek.Int32,
		StartTime:       arg.StartTime,
		EndTime:         arg.EndTime,
		MaxAppointments: doctorSchedule.MaxAppointments.Int32,
	}, nil
}

func (s *UserService) getDoctorByID(ctx *gin.Context, userID int64) (*DoctorResponse, error) {

	doctor, err := s.storeDB.GetDoctor(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "doctor not found")
			return nil, fmt.Errorf("doctor not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	return &DoctorResponse{
		ID:             doctor.ID,
		Specialization: doctor.Specialization.String,
		Name:           doctor.Name,
		YearsOfExp:     doctor.YearsOfExperience.Int32,
		Education:      doctor.Education.String,
		Certificate:    doctor.CertificateNumber.String,
		Bio:            doctor.Bio.String,
	}, nil
}
