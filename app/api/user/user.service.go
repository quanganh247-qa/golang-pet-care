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
	createUserService(ctx *gin.Context, req createUserRequest) error
	getUserDetailsService(ctx *gin.Context, username string) (*UserResponse, error)
	getAllUsersService(ctx *gin.Context) ([]UserResponse, error)
	loginUserService(ctx *gin.Context, req loginUserRequest) (*loginUSerResponse, error)
	logoutUsersService(ctx *gin.Context, username string, token string) error
	verifyEmailService(ctx *gin.Context, req VerrifyEmailTxParams) (VerrifyEmailTxResult, error)
	createDoctorService(ctx *gin.Context, arg InsertDoctorRequest, username string) (*DoctorResponse, error)
	createDoctorScheduleService(ctx *gin.Context, arg InsertDoctorScheduleRequest, username string) (*DoctorScheduleResponse, error)
	getDoctorByID(ctx *gin.Context, userID int64) (*DoctorResponse, error)
	insertTimeSlots(ctx *gin.Context, username string, arg db.InsertTimeslotParams) (*db.Timeslot, error)
	GetTimeslotsAvailable(ctx *gin.Context, doctorID int64, date string) ([]db.GetTimeslotsAvailableRow, error)
	GetAllTimeslots(ctx *gin.Context, doctorID int64, date string) ([]db.GetTimeslotsAvailableRow, error)
	UpdateDoctorAvailable(ctx *gin.Context, time_slot_id int64) error
	// InsertTokenInfoService(ctx *gin.Context, arg InsertTokenInfoRequest, username string) (*db.TokenInfo, error)
}

func (server *UserService) createUserService(ctx *gin.Context, req createUserRequest) error {

	hashedPwd, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprintf("cannot hash password: %v", err))
		return fmt.Errorf("cannot hash password: %v", err)
	}

	arg := db.CreateUserParams{
		Username:        req.Username,
		HashedPassword:  hashedPwd,
		FullName:        req.FullName,
		Email:           req.Email,
		PhoneNumber:     pgtype.Text{String: req.PhoneNumber, Valid: true},
		Address:         pgtype.Text{String: req.Address, Valid: true},
		DataImage:       req.DataImage,
		OriginalImage:   pgtype.Text{String: req.OriginalImage, Valid: true},
		Role:            pgtype.Text{String: "user", Valid: true}, //
		IsVerifiedEmail: pgtype.Bool{Bool: true, Valid: true},
	}

	_, err = server.storeDB.CreateUser(ctx, arg) // Check this line carefully

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

	return nil
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
	device_tokens, err := service.storeDB.GetDeviceTokenByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	var device_tokens_response []string
	for _, d := range device_tokens {
		device_tokens_response = append(device_tokens_response, d.Token)
	}

	return &loginUSerResponse{
		AccessTokenExpiresAt:  time.Now().Add(time.Hour),
		RefreshTokenExpiresAt: time.Now().Add(time.Hour * 24),
		User: UserResponse{
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
		},
		DeviceToken: device_tokens_response,
	}, nil
}

func (service *UserService) logoutUsersService(ctx *gin.Context, username string, token string) error {
	host, secure := util.SetCookieSameSite(ctx)
	ctx.SetCookie("refresh_token", "", -1, "/", host, secure, true)
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

	doctorSchedule, err := s.storeDB.InsertDoctorSchedule(ctx, db.InsertDoctorScheduleParams{
		DoctorID:        doctor.ID,
		DayOfWeek:       pgtype.Int4{Int32: arg.Day, Valid: true},
		MaxAppointments: pgtype.Int4{Int32: arg.MaxAppoin, Valid: true},
		StartTime:       pgtype.Timestamp{Time: startTime, Valid: true},
		EndTime:         pgtype.Timestamp{Time: endTime, Valid: true},
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

// Time off
func (s *UserService) insertTimeSlots(ctx *gin.Context, username string, arg db.InsertTimeslotParams) (*db.Timeslot, error) {
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
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}
	arg.DoctorID = doctor.ID
	timeslot, err := s.storeDB.InsertTimeslot(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to insert timeslot: %w", err)
	}
	return &db.Timeslot{
		ID:        timeslot.ID,
		DoctorID:  timeslot.DoctorID,
		StartTime: timeslot.StartTime,
		EndTime:   timeslot.EndTime,
		IsActive:  timeslot.IsActive,
		Day:       timeslot.Day,
	}, nil
}

func ConvertStringToPgtypeDate(dateStr string) (pgtype.Date, error) {
	// Define the date format
	const dateFormat = "2006-01-02"

	// Parse the string to time.Time
	parsedTime, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return pgtype.Date{}, fmt.Errorf("failed to parse date: %w", err)
	}
	return pgtype.Date{
		Time:  parsedTime,
		Valid: true,
	}, nil
}

func (s *UserService) GetTimeslotsAvailable(ctx *gin.Context, doctorID int64, date string) ([]db.GetTimeslotsAvailableRow, error) {

	//convert string to pgtype.Date
	pgDate, err := ConvertStringToPgtypeDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to pgtype.Date: %w", err)
	}
	var listTimeSlot []db.GetTimeslotsAvailableRow
	timeslots, err := s.storeDB.GetTimeslotsAvailable(ctx, db.GetTimeslotsAvailableParams{
		ID:  doctorID,
		Day: pgDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get timeslots available: %w", err)
	}

	for _, ts := range timeslots {
		listTimeSlot = append(listTimeSlot, db.GetTimeslotsAvailableRow{
			StartTime: ts.StartTime,
			EndTime:   ts.EndTime,
		})
	}
	return listTimeSlot, nil
}

func (s *UserService) GetAllTimeslots(ctx *gin.Context, doctorID int64, date string) ([]db.GetTimeslotsAvailableRow, error) {

	//convert string to pgtype.Date
	pgDate, err := ConvertStringToPgtypeDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to pgtype.Date: %w", err)
	}
	var listTimeSlot []db.GetTimeslotsAvailableRow
	timeslots, err := s.storeDB.GetAllTimeSlots(ctx, db.GetAllTimeSlotsParams{
		ID:  doctorID,
		Day: pgDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get timeslots : %w", err)
	}

	for _, ts := range timeslots {
		listTimeSlot = append(listTimeSlot, db.GetTimeslotsAvailableRow{
			StartTime: ts.StartTime,
			EndTime:   ts.EndTime,
		})
	}
	return listTimeSlot, nil
}

func (s *UserService) UpdateDoctorAvailable(ctx *gin.Context, timeSlotID int64) error {
	// Use a transaction to ensure that the update is atomic
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.UpdateDoctorAvailable(ctx, db.UpdateDoctorAvailableParams{
			IsActive: pgtype.Bool{Bool: false, Valid: true}, // Setting is_active to false
			ID:       timeSlotID,
		})

		if err != nil {
			return fmt.Errorf("failed to update doctor availability: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to execute transaction for updating doctor availability: %w", err)
	}

	return nil // Successfully updated
}
