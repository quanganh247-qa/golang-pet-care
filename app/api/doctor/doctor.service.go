package doctor

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DoctorServiceInterface interface {
	LoginDoctorService(ctx *gin.Context, req loginDoctorRequest) (*loginDoctorResponse, error)
	GetDoctorProfile(ctx *gin.Context, username string) (*DoctorDetail, error)
	GetAllDoctorService(ctx *gin.Context) ([]DoctorDetail, error)
	GetShifts(ctx *gin.Context) ([]Shift, error)
	CreateShift(ctx *gin.Context, req CreateShiftRequest) (*Shift, error)
	GetShiftByDoctorId(ctx *gin.Context, doctorId int64) ([]ShiftResponse, error)
	GetDoctorById(ctx *gin.Context, doctorId int64) (DoctorDetail, error)
}

func (service *DoctorService) LoginDoctorService(ctx *gin.Context, req loginDoctorRequest) (*loginDoctorResponse, error) {
	// First get the user by username
	user, err := service.storeDB.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, "user not found")
			return nil, fmt.Errorf("user not found")
		}
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return nil, fmt.Errorf("internal server error: %v", err)
	}

	// Check if the user is a doctor (role check)
	if user.Role.String == "user" {
		ctx.JSON(http.StatusForbidden, "access denied: not a doctor account")
		return nil, fmt.Errorf("access denied: not a doctor account")
	}

	// Verify password
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "incorrect password")
		return nil, fmt.Errorf("incorrect password")
	}

	// Check if email is verified
	if !user.IsVerifiedEmail.Bool {
		ctx.JSON(http.StatusForbidden, "email not verified")
		return nil, fmt.Errorf("email not verified")
	}

	// Get doctor details
	doctor, err := service.storeDB.GetDoctorByUserId(ctx, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "doctor profile not found")
		return nil, fmt.Errorf("doctor profile not found: %v", err)
	}

	accessToken, accessPayload, err := token.TokenMaker.CreateToken(req.Username,
		map[string]bool{"doctor": true},
		util.Configs.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to create access token")
		return nil, fmt.Errorf("failed to create access token: %v", err)
	}

	// Create refresh token
	refreshToken, refreshPayload, err := token.TokenMaker.CreateToken(
		user.Username,
		map[string]bool{"doctor": true},
		util.Configs.RefreshTokenDuration, // 7 days
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to create refresh token")
		return nil, fmt.Errorf("failed to create refresh token: %v", err)
	}

	return &loginDoctorResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		Doctor: DoctorDetail{
			DoctorID:       doctor.ID,
			Username:       user.Username,
			DoctorName:     user.FullName,
			Email:          user.Email,
			Specialization: doctor.Specialization.String,
			YearsOfExp:     doctor.YearsOfExperience.Int32,
			Education:      doctor.Education.String,
			Certificate:    doctor.CertificateNumber.String,
			Bio:            doctor.Bio.String,
			DataImage:      []byte(user.DataImage),
		},
	}, nil
}

func (service *DoctorService) GetDoctorProfile(ctx *gin.Context, username string) (*DoctorDetail, error) {
	// Get user details
	user, err := service.storeDB.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Check if the user is a doctor
	if user.Role.String != "doctor" {
		return nil, fmt.Errorf("user is not a doctor")
	}

	// Get doctor details
	doctor, err := service.storeDB.GetDoctorByUserId(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor profile: %v", err)
	}

	return &DoctorDetail{
		DoctorID:       doctor.ID,
		Username:       user.Username,
		DoctorName:     user.FullName,
		Email:          user.Email,
		Specialization: doctor.Specialization.String,
		YearsOfExp:     doctor.YearsOfExperience.Int32,
		Education:      doctor.Education.String,
		Certificate:    doctor.CertificateNumber.String,
		Bio:            doctor.Bio.String,
		DataImage:      []byte(user.DataImage),
	}, nil
}

func (service *DoctorService) GetAllDoctorService(ctx *gin.Context) ([]DoctorDetail, error) {
	doctor, err := service.storeDB.ListDoctors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all doctor: %w", err)
	}
	doctorList := make([]DoctorDetail, 0)
	for _, d := range doctor {
		doctorList = append(doctorList, DoctorDetail{
			DoctorID:       d.DoctorID,
			Username:       d.Username,
			DoctorName:     d.FullName,
			Email:          d.Email,
			Role:           d.Role.String,
			Specialization: d.Specialization.String,
			YearsOfExp:     d.YearsOfExperience.Int32,
			Education:      d.Education.String,
			Certificate:    d.CertificateNumber.String,
			Bio:            d.Bio.String,
			DataImage:      []byte(d.DataImage),
		})
	}
	return doctorList, nil
}

// Get shifts
func (service *DoctorService) GetShifts(ctx *gin.Context) ([]Shift, error) {
	shifts, err := service.storeDB.GetShifts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shifts: %w", err)
	}
	shiftList := make([]Shift, 0)
	for _, s := range shifts {

		parsedTime, err := time.Parse("2006-01-02 15:04:05", s.StartTime.Time.Format("2006-01-02 15:04:05"))
		if err != nil {
			fmt.Println("Error parsing time:", err)
		}

		shiftList = append(shiftList, Shift{
			ID:               s.ID,
			StartTime:        parsedTime,
			EndTime:          s.EndTime.Time,
			AssignedPatients: s.AssignedPatients.Int32,
			CreatedAt:        s.CreatedAt.Time,
			DoctorID:         s.DoctorID,
		})
	}
	return shiftList, nil
}

// Create shift
func (service *DoctorService) CreateShift(ctx *gin.Context, req CreateShiftRequest) (*Shift, error) {
	shift, err := service.storeDB.CreateShift(ctx, db.CreateShiftParams{
		StartTime: pgtype.Timestamp{Time: req.StartTime, Valid: true},
		EndTime:   pgtype.Timestamp{Time: req.EndTime, Valid: true},
		DoctorID:  req.DoctorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create shift: %w", err)
	}
	return &Shift{
		ID:               shift.ID,
		StartTime:        shift.StartTime.Time,
		EndTime:          shift.EndTime.Time,
		AssignedPatients: shift.AssignedPatients.Int32,
		CreatedAt:        shift.CreatedAt.Time,
		DoctorID:         shift.DoctorID,
	}, nil
}

// Get shift by doctor id
func (service *DoctorService) GetShiftByDoctorId(ctx *gin.Context, doctorId int64) ([]ShiftResponse, error) {
	shifts, err := service.storeDB.GetShiftByDoctorId(ctx, doctorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get shifts: %w", err)
	}
	shiftList := make([]ShiftResponse, 0)
	for _, s := range shifts {
		shiftList = append(shiftList, ShiftResponse{
			ID:               s.ID,
			StartTime:        s.StartTime.Time.Format("2006-01-02 15:04:05"),
			EndTime:          s.EndTime.Time.Format("2006-01-02 15:04:05"),
			AssignedPatients: s.AssignedPatients.Int32,
			DoctorID:         s.DoctorID,
		})
	}
	return shiftList, nil
}

func (service *DoctorService) GetDoctorById(ctx *gin.Context, doctorId int64) (DoctorDetail, error) {
	doctor, err := service.storeDB.GetDoctor(ctx, doctorId)
	if err != nil {
		return DoctorDetail{}, fmt.Errorf("failed to get doctor: %w", err)
	}
	return DoctorDetail{
		DoctorID:       doctor.ID,
		DoctorName:     doctor.Name,
		Specialization: doctor.Specialization.String,
		YearsOfExp:     doctor.YearsOfExperience.Int32,
		Education:      doctor.Education.String,
		Role:           doctor.Role.String,
		Certificate:    doctor.CertificateNumber.String,
		Bio:            doctor.Bio.String,
	}, nil
}
