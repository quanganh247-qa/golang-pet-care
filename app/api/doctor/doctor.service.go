package doctor

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
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
	DeleteShift(ctx *gin.Context, shiftId int64) error
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

	// Invalidate cache
	if service.redis != nil {
		// Clear doctor by ID cache
		service.redis.RemoveDoctorCache(doctor.ID)
		// Clear doctor by username cache
		service.redis.RemoveDoctorByUsernameCache(username)
		// Clear all doctors list cache
		service.redis.ClearAllDoctorCache()
	}

	return nil
}

func (service *DoctorService) GetDoctorProfile(ctx *gin.Context, username string) (*DoctorDetail, error) {
	// Try to get from cache first
	if service.redis != nil {
		doctorInfo, err := service.redis.DoctorByUsernameLoadCache(username)
		if err == nil {
			// Found in cache
			return &DoctorDetail{
				DoctorID:       doctorInfo.DoctorID,
				Username:       doctorInfo.Username,
				FullName:       doctorInfo.FullName,
				Address:        doctorInfo.Address,
				PhoneNumber:    doctorInfo.PhoneNumber,
				Email:          doctorInfo.Email,
				Role:           doctorInfo.Role,
				DoctorName:     doctorInfo.DoctorName,
				Specialization: doctorInfo.Specialization,
				YearsOfExp:     doctorInfo.YearsOfExp,
				Education:      doctorInfo.Education,
				Certificate:    doctorInfo.Certificate,
				Bio:            doctorInfo.Bio,
				DataImage:      doctorInfo.DataImage,
			}, nil
		}
	}

	// Cache miss, get from database
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

	doctorDetail := &DoctorDetail{
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
	}

	return doctorDetail, nil
}

func (service *DoctorService) GetAllDoctorService(ctx *gin.Context) ([]DoctorDetail, error) {
	// Try to get from cache first
	if service.redis != nil {
		doctorList, err := service.redis.DoctorsListLoadCache()
		if err == nil {
			// Found in cache
			result := make([]DoctorDetail, len(doctorList))
			for i, doc := range doctorList {
				result[i] = DoctorDetail{
					DoctorID:       doc.DoctorID,
					Username:       doc.Username,
					DoctorName:     doc.DoctorName,
					Email:          doc.Email,
					Role:           doc.Role,
					Specialization: doc.Specialization,
					YearsOfExp:     doc.YearsOfExp,
					Education:      doc.Education,
					Certificate:    doc.Certificate,
					Bio:            doc.Bio,
					DataImage:      doc.DataImage,
				}
			}
			return result, nil
		}
	}

	// Cache miss, get from database
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
			DataImage:      d.DataImage,
		})
	}
	return doctorList, nil
}

// Get shifts
func (service *DoctorService) GetShifts(ctx *gin.Context) ([]Shift, error) {
	// Try to get from cache first
	if service.redis != nil {
		shiftInfos, err := service.redis.AllShiftsLoadCache()
		if err == nil {
			// Found in cache
			result := make([]Shift, len(shiftInfos))
			for i, s := range shiftInfos {
				result[i] = Shift{
					ID:         s.ID,
					DoctorID:   s.DoctorID,
					StartTime:  s.StartTime.Format("2006-01-02 15:04:05"),
					EndTime:    s.EndTime.Format("2006-01-02 15:04:05"),
					Date:       s.Date,
					CreatedAt:  s.CreatedAt.Format("2006-01-02 15:04:05"),
					DoctorName: s.DoctorName,
				}
			}
			return result, nil
		}
	}

	// Cache miss, get from database
	shifts, err := service.storeDB.GetShifts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shifts: %w", err)
	}
	shiftList := make([]Shift, 0)
	for _, s := range shifts {

		doctor, err := service.storeDB.GetDoctor(ctx, s.DoctorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get doctor: %w", err)
		}

		timeSlots, err := service.storeDB.GetTimeSlotsByShiftID(ctx, s.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get time slots: %w", err)
		}

		// Skip if no time slots found
		if len(timeSlots) == 0 {
			// Use created_at time if no time slots
			shiftList = append(shiftList, Shift{
				ID:         s.ID,
				StartTime:  s.CreatedAt.Time.Format("2006-01-02 15:04:05"),
				EndTime:    s.CreatedAt.Time.Add(time.Hour).Format("2006-01-02 15:04:05"),
				Date:       s.Date.Time.Format("2006-01-02"),
				CreatedAt:  s.CreatedAt.Time.Format("2006-01-02 15:04:05"),
				DoctorID:   s.DoctorID,
				DoctorName: doctor.Name,
			})
			continue
		}

		// Convert microseconds to hours/minutes/seconds for formatting
		startHours := int((timeSlots[0].StartTime.Microseconds / 1000000) / 3600)
		startMins := int(((timeSlots[0].StartTime.Microseconds / 1000000) % 3600) / 60)
		startSecs := int((timeSlots[0].StartTime.Microseconds / 1000000) % 60)

		endHours := int((timeSlots[len(timeSlots)-1].EndTime.Microseconds / 1000000) / 3600)
		endMins := int(((timeSlots[len(timeSlots)-1].EndTime.Microseconds / 1000000) % 3600) / 60)
		endSecs := int((timeSlots[len(timeSlots)-1].EndTime.Microseconds / 1000000) % 60)

		// Format times as strings in the standard time format
		startTimeStr := fmt.Sprintf("%02d:%02d:%02d", startHours, startMins, startSecs)
		endTimeStr := fmt.Sprintf("%02d:%02d:%02d", endHours, endMins, endSecs)

		// Format the created_at timestamp
		createdAtStr := s.CreatedAt.Time.Format("2006-01-02 15:04:05")

		shiftList = append(shiftList, Shift{
			ID:         s.ID,
			StartTime:  startTimeStr,
			EndTime:    endTimeStr,
			Date:       s.Date.Time.Format("2006-01-02"),
			CreatedAt:  createdAtStr,
			DoctorID:   s.DoctorID,
			DoctorName: doctor.Name,
		})
	}
	return shiftList, nil
}

// Create shift
func (service *DoctorService) CreateShift(ctx *gin.Context, req CreateShiftRequest) (*Shift, error) {
	// Parse string time inputs to time.Time objects
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start time format: %w", err)
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end time format: %w", err)
	}

	// Extract just the date portion from the start time string
	dateStr := strings.Split(req.StartTime, " ")[0]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	shift, err := service.storeDB.CreateShift(ctx, db.CreateShiftParams{
		Date:     pgtype.Date{Time: date, Valid: true},
		DoctorID: req.DoctorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create shift: %w", err)
	}

	// Invalidate cache
	if service.redis != nil {
		// Clear doctor shifts cache
		service.redis.ClearDoctorShiftsCache(req.DoctorID)
	}

	// Create time slots for the shift, 1 hour each
	// First convert doctor ID to int32 as required by the CreateTimeSlotParams
	doctorID32 := int32(req.DoctorID)

	slot := (endTime.Hour() - startTime.Hour())
	if slot <= 0 {
		// Handle cases where the shift spans over midnight or is less than an hour
		slot = int(endTime.Sub(startTime).Hours())
	}

	// Get the location (timezone) from the input time
	location := startTime.Location()

	// Create a time slot for each hour in the shift
	currentSlotStart := startTime
	for currentSlotStart.Before(endTime) {
		// Calculate the end time for this slot (1 hour later)
		currentSlotEnd := currentSlotStart.Add(time.Hour)

		// If this would go past the shift end time, use the shift end time instead
		if currentSlotEnd.After(endTime) {
			currentSlotEnd = endTime
		}
		// Extract the time components for start and end using the correct location
		startLocalTime := currentSlotStart.In(location)
		endLocalTime := currentSlotEnd.In(location)
		// Parse these strings into pgtype.Time values
		startTimeOnly := pgtype.Time{}
		endTimeOnly := pgtype.Time{}

		// Set the microseconds directly using the total seconds in the day
		startSeconds := startLocalTime.Hour()*3600 + startLocalTime.Minute()*60 + startLocalTime.Second()
		endSeconds := endLocalTime.Hour()*3600 + endLocalTime.Minute()*60 + endLocalTime.Second()

		startTimeOnly.Microseconds = int64(startSeconds) * 1000000
		startTimeOnly.Valid = true

		endTimeOnly.Microseconds = int64(endSeconds) * 1000000
		endTimeOnly.Valid = true

		// Create the time slot
		_, err = service.storeDB.CreateTimeSlot(ctx, db.CreateTimeSlotParams{
			DoctorID:       doctorID32,
			ShiftID:        shift.ID,
			Date:           pgtype.Date{Time: date, Valid: true},
			StartTime:      startTimeOnly,
			EndTime:        endTimeOnly,
			MaxPatients:    pgtype.Int4{Int32: int32(slot), Valid: true},
			BookedPatients: pgtype.Int4{Int32: 0, Valid: true},
		})

		if err != nil {
			// Log the error but continue creating other slots
			fmt.Printf("Error creating time slot: %v\n", err)
		}

		// Move to the next slot
		currentSlotStart = currentSlotEnd
	}

	return &Shift{
		ID:        shift.ID,
		StartTime: startTime.Format("2006-01-02 15:04:05"),
		EndTime:   endTime.Format("2006-01-02 15:04:05"),
		Date:      shift.Date.Time.Format("2006-01-02"),
		CreatedAt: shift.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		DoctorID:  shift.DoctorID,
	}, nil
}

// Get shift by doctor id
func (service *DoctorService) GetShiftByDoctorId(ctx *gin.Context, doctorId int64) ([]ShiftResponse, error) {
	// Try to get from cache first
	if service.redis != nil {
		shiftInfos, err := service.redis.DoctorShiftsLoadCache(doctorId)
		if err == nil {
			// Found in cache
			result := make([]ShiftResponse, len(shiftInfos))
			for i, s := range shiftInfos {
				result[i] = ShiftResponse{
					ID:         s.ID,
					DoctorID:   s.DoctorID,
					StartTime:  s.StartTime.Format("2006-01-02 15:04:05"),
					EndTime:    s.EndTime.Format("2006-01-02 15:04:05"),
					Date:       s.Date,
					DoctorName: s.DoctorName,
				}
			}
			return result, nil
		}
	}

	// Cache miss, get from database
	shifts, err := service.storeDB.GetShiftByDoctorId(ctx, doctorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get shifts: %w", err)
	}
	shiftList := make([]ShiftResponse, 0)
	for _, s := range shifts {
		doctor, err := service.storeDB.GetDoctor(ctx, s.DoctorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get doctor: %w", err)
		}

		// Convert microseconds to hours/minutes/seconds for formatting
		startHours := int((s.StartTime.Microseconds / 1000000) / 3600)
		startMins := int(((s.StartTime.Microseconds / 1000000) % 3600) / 60)
		startSecs := int((s.StartTime.Microseconds / 1000000) % 60)

		endHours := int((s.EndTime.Microseconds / 1000000) / 3600)
		endMins := int(((s.EndTime.Microseconds / 1000000) % 3600) / 60)
		endSecs := int((s.EndTime.Microseconds / 1000000) % 60)

		// Format times as strings in the standard time format
		startTimeStr := fmt.Sprintf("%02d:%02d:%02d", startHours, startMins, startSecs)
		endTimeStr := fmt.Sprintf("%02d:%02d:%02d", endHours, endMins, endSecs)

		shiftList = append(shiftList, ShiftResponse{
			ID:         s.ID,
			Date:       s.Date.Time.Format("2006-01-02"),
			StartTime:  startTimeStr,
			EndTime:    endTimeStr,
			DoctorID:   s.DoctorID,
			DoctorName: doctor.Name,
		})
	}
	return shiftList, nil
}

func (service *DoctorService) GetDoctorById(ctx *gin.Context, doctorId int64) (DoctorDetail, error) {
	// Try to get from cache first
	if service.redis != nil {
		doctorInfo, err := service.redis.DoctorLoadCache(doctorId)
		if err == nil {
			// Found in cache
			return DoctorDetail{
				DoctorID:       doctorInfo.DoctorID,
				Username:       doctorInfo.Username,
				DoctorName:     doctorInfo.DoctorName,
				Email:          doctorInfo.Email,
				Role:           doctorInfo.Role,
				Specialization: doctorInfo.Specialization,
				YearsOfExp:     doctorInfo.YearsOfExp,
				Education:      doctorInfo.Education,
				Certificate:    doctorInfo.Certificate,
				Bio:            doctorInfo.Bio,
				DataImage:      doctorInfo.DataImage,
			}, nil
		}
	}

	// Cache miss, get from database
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

// Delete shift
func (service *DoctorService) DeleteShift(ctx *gin.Context, shiftId int64) error {
	// Get the shift to determine the doctor ID before deleting
	var doctorID int64
	if service.redis != nil {
		// Try to get all shifts from cache
		shiftInfos, err := service.redis.AllShiftsLoadCache()
		if err == nil {
			// Found in cache, find the shift with the given ID
			for _, s := range shiftInfos {
				if s.ID == shiftId {
					doctorID = s.DoctorID
					break
				}
			}
		}
	}

	// If we couldn't get the doctor ID from cache, execute the transaction as normal
	err := service.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// If we didn't get the doctor ID from cache, just delete the shift
		// We'll clear all shifts in cache afterward
		return q.DeleteShift(ctx, shiftId)
	})

	// Invalidate cache
	if service.redis != nil && err == nil {
		// If we have the doctor ID, clear that doctor's shifts
		if doctorID != 0 {
			service.redis.ClearDoctorShiftsCache(doctorID)
		} else {
			// Otherwise clear all shifts caches
			service.redis.ClearAllDoctorCache()
		}
	}

	return err
}
