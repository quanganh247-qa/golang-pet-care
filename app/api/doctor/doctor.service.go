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
	EditDoctorProfileService(ctx *gin.Context, username string, req EditDoctorProfileRequest) error
	LoginDoctorService(ctx *gin.Context, req loginDoctorRequest) (*loginDoctorResponse, error)
	GetDoctorProfile(ctx *gin.Context, username string) (*DoctorDetail, error)
	GetAllDoctorService(ctx *gin.Context) ([]DoctorDetail, error)
	GetShifts(ctx *gin.Context) ([]Shift, error)
	CreateShift(ctx *gin.Context, req CreateShiftRequest) (*Shift, error)
	GetShiftByDoctorId(ctx *gin.Context, doctorId int64) ([]ShiftResponse, error)
	GetDoctorById(ctx *gin.Context, doctorId int64) (DoctorDetail, error)
	CreateLeaveRequest(ctx *gin.Context, doctorID int64, req CreateLeaveRequest) (*LeaveRequest, error)
	GetLeaveRequests(ctx *gin.Context, doctorID int64, page, pageSize int) ([]LeaveRequest, error)
	UpdateLeaveRequest(ctx *gin.Context, leaveID int64, reviewerID int64, req UpdateLeaveRequest) error
	GetDoctorAttendance(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) ([]AttendanceRecord, error)
	GetDoctorWorkload(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) (*WorkloadMetrics, error)
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

func (service *DoctorService) EditDoctorProfileService(ctx *gin.Context, username string, req EditDoctorProfileRequest) error {
	user, err := service.storeDB.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	doctor, err := service.storeDB.GetDoctorByUserId(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to get doctor profile: %v", err)
	}

	var reqUpdate db.UpdateDoctorParams
	reqUpdate.ID = doctor.ID
	if req.Specialization != "" {
		reqUpdate.Specialization = pgtype.Text{String: req.Specialization, Valid: true}
	}
	if req.YearsOfExp != 0 {
		reqUpdate.YearsOfExperience = pgtype.Int4{Int32: req.YearsOfExp, Valid: true}
	}
	if req.Education != "" {
		reqUpdate.Education = pgtype.Text{String: req.Education, Valid: true}
	}
	if req.Certificate != "" {
		reqUpdate.CertificateNumber = pgtype.Text{String: req.Certificate, Valid: true}
	}
	if req.Bio != "" {
		reqUpdate.Bio = pgtype.Text{String: req.Bio, Valid: true}
	}

	_, err = service.storeDB.UpdateDoctor(ctx, reqUpdate)
	if err != nil {
		return fmt.Errorf("failed to update doctor profile: %v", err)
	}
	return nil
}

func (service *DoctorService) GetDoctorProfile(ctx *gin.Context, username string) (*DoctorDetail, error) {
	// Get user details
	user, err := service.storeDB.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Get doctor details
	doctor, err := service.storeDB.GetDoctorByUserId(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor profile: %v", err)
	}

	return &DoctorDetail{
		DoctorID:       doctor.ID,
		Username:       user.Username,
		FullName:       user.FullName,
		Address:        user.Address.String,
		PhoneNumber:    user.PhoneNumber.String,
		Email:          user.Email,
		Role:           user.Role.String,
		DoctorName:     user.FullName,
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

func (service *DoctorService) CreateLeaveRequest(ctx *gin.Context, doctorID int64, req CreateLeaveRequest) (*LeaveRequest, error) {
	// First check if there are any existing shifts during the leave period
	shifts, err := service.storeDB.GetShiftByDoctorId(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing shifts: %v", err)
	}

	// Check for conflicting shifts
	for _, shift := range shifts {
		if shift.StartTime.Time.Before(req.EndDate) && shift.EndTime.Time.After(req.StartDate) {
			return nil, fmt.Errorf("leave request conflicts with existing shifts")
		}
	}

	// Create the leave request
	leave, err := service.storeDB.CreateLeaveRequest(ctx, db.CreateLeaveRequestParams{
		DoctorID:  doctorID,
		StartDate: pgtype.Timestamp{Time: req.StartDate, Valid: true},
		EndDate:   pgtype.Timestamp{Time: req.EndDate, Valid: true},
		LeaveType: req.LeaveType,
		Reason:    pgtype.Text{String: req.Reason, Valid: true},
		Status:    pgtype.Text{String: "pending", Valid: true},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create leave request: %v", err)
	}

	return &LeaveRequest{
		ID:        leave.ID,
		DoctorID:  leave.DoctorID,
		StartDate: leave.StartDate.Time,
		EndDate:   leave.EndDate.Time,
		LeaveType: leave.LeaveType,
		Reason:    leave.Reason.String,
		Status:    leave.Status.String,
		CreatedAt: leave.CreatedAt.Time,
		UpdatedAt: leave.UpdatedAt.Time,
	}, nil
}

func (service *DoctorService) GetLeaveRequests(ctx *gin.Context, doctorID int64, page, pageSize int) ([]LeaveRequest, error) {
	offset := (page - 1) * pageSize
	leaves, err := service.storeDB.GetLeaveRequests(ctx,
		db.GetLeaveRequestsParams{
			DoctorID: doctorID,
			Limit:    int32(pageSize),
			Offset:   int32(offset),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get leave requests: %v", err)
	}

	var leaveRequests []LeaveRequest
	for _, leave := range leaves {
		leaveRequests = append(leaveRequests, LeaveRequest{
			ID:          leave.ID,
			DoctorID:    leave.DoctorID,
			StartDate:   leave.StartDate.Time,
			EndDate:     leave.EndDate.Time,
			LeaveType:   leave.LeaveType,
			Reason:      leave.Reason.String,
			Status:      leave.Status.String,
			ReviewedBy:  leave.ReviewedBy.Int64,
			ReviewNotes: leave.ReviewNotes.String,
			CreatedAt:   leave.CreatedAt.Time,
			UpdatedAt:   leave.UpdatedAt.Time,
		})
	}

	return leaveRequests, nil
}

func (service *DoctorService) UpdateLeaveRequest(ctx *gin.Context, leaveID int64, reviewerID int64, req UpdateLeaveRequest) error {
	_, err := service.storeDB.UpdateLeaveRequest(ctx, db.UpdateLeaveRequestParams{
		ID:          leaveID,
		Status:      pgtype.Text{String: req.Status, Valid: true},
		ReviewedBy:  pgtype.Int8{Int64: reviewerID, Valid: true},
		ReviewNotes: pgtype.Text{String: req.ReviewNotes, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update leave request: %v", err)
	}
	return nil
}

func (service *DoctorService) GetDoctorAttendance(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) ([]AttendanceRecord, error) {
	records, err := service.storeDB.GetDoctorAttendance(ctx, db.GetDoctorAttendanceParams{
		DoctorID:    doctorID,
		StartTime:   pgtype.Timestamp{Time: startDate, Valid: true},
		StartTime_2: pgtype.Timestamp{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor attendance: %v", err)
	}

	var attendance []AttendanceRecord
	for _, record := range records {
		attendance = append(attendance, AttendanceRecord{
			WorkDate:              record.WorkDate.Time,
			TotalAppointments:     int(record.TotalAppointments),
			CompletedAppointments: int(record.CompletedAppointments),
			WorkHours:             float64(record.WorkHours),
		})
	}

	return attendance, nil
}

func (service *DoctorService) GetDoctorWorkload(ctx *gin.Context, doctorID int64, startDate, endDate time.Time) (*WorkloadMetrics, error) {
	metrics, err := service.storeDB.GetDoctorWorkload(ctx, db.GetDoctorWorkloadParams{
		DoctorID:    doctorID,
		StartTime:   pgtype.Timestamp{Time: startDate, Valid: true},
		StartTime_2: pgtype.Timestamp{Time: endDate, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor workload: %v", err)
	}

	// Convert Numeric to float64
	avgWorkHours, err := metrics.AvgWorkHoursPerDay.Float64Value()
	if err != nil {
		return nil, fmt.Errorf("failed to convert average work hours: %v", err)
	}

	return &WorkloadMetrics{
		TotalAppointments:     int(metrics.TotalAppointments),
		CompletedAppointments: int(metrics.CompletedAppointments),
		AvgWorkHoursPerDay:    avgWorkHours.Float64,
		TotalWorkDays:         int(metrics.TotalWorkDays),
	}, nil
}
