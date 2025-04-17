package appointment

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type AppointmentServiceInterface interface {
	CreateSOAPService(ctx *gin.Context, soap CreateSOAPRequest, appointmentID int64) (*SOAPResponse, error)
	UpdateSOAPService(ctx *gin.Context, soap UpdateSOAPRequest, appointmentID int64) (*SOAPResponse, error)
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error)
	CreateWalkInAppointment(ctx *gin.Context, req createWalkInAppointmentRequest) (*createAppointmentResponse, error)
	ConfirmPayment(ctx context.Context, appointmentID int64) error
	CheckInAppoinment(ctx *gin.Context, id, roomID int64, priority string) error
	GetAppointmentByID(ctx *gin.Context, id int64) (*Appointment, error)
	GetAppointmentsByUser(ctx *gin.Context, username string) ([]Appointment, error)
	GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error)
	GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error)
	GetAllAppointments(ctx *gin.Context, date string, option string, pagination *util.Pagination) (*util.PaginationResponse[Appointment], error)
	GetAllAppointmentsByDate(ctx *gin.Context, pagination *util.Pagination, date string) ([]Appointment, error)
	UpdateAppointmentService(ctx *gin.Context, req updateAppointmentRequest, appointmentID int64) error
	UpdateQueueItemStatusService(ctx *gin.Context, id int64, status string) error
	GetQueueService(ctx *gin.Context, username string) ([]QueueItem, error)
	GetHistoryAppointmentsByPetID(ctx *gin.Context, petID int64) ([]historyAppointmentResponse, error)
	GetSOAPByAppointmentID(ctx *gin.Context, appointmentID int64) (*SOAPResponse, error)
	GetAppointmentDistributionByService(ctx *gin.Context, startDate, endDate string) ([]AppointmentDistribution, error)
	HandleWebSocket(ctx *gin.Context)
	GetPendingNotifications(ctx *gin.Context) ([]websocket.OfflineMessage, error)
	MarkMessageDelivered(ctx *gin.Context, id int64) error
}

func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error) {
	var err error
	var timeSlot db.TimeSlot
	var doctor user.DoctorResponse
	var service db.Service
	var wg sync.WaitGroup

	// Generate a lock key for this time slot to prevent concurrent booking
	lockKey := fmt.Sprintf("timeslot:%d", req.TimeSlotID)

	// Try to acquire a distributed lock with 10-second timeout
	lockAcquired := AcquireLock(lockKey, 10*time.Second)
	if !lockAcquired {
		return nil, fmt.Errorf("another appointment is being processed for this time slot, please try again shortly")
	}

	// Make sure the lock is released when we're done
	defer ReleaseLock(lockKey)

	errChan := make(chan error, 3)
	wg.Add(3)

	go func() {
		defer wg.Done()
		d, err := s.storeDB.GetDoctor(ctx, req.DoctorID)
		if err != nil {
			errChan <- fmt.Errorf("failed to get doctor: %w", err)
			return
		}
		doctor = user.DoctorResponse{
			ID:             d.ID,
			Specialization: d.Specialization.String,
			Name:           d.Name,
			YearsOfExp:     d.YearsOfExperience.Int32,
			Education:      d.Education.String,
			Certificate:    d.CertificateNumber.String,
			Bio:            d.Bio.String,
		}
	}()

	go func() {
		defer wg.Done()
		ts, err := s.storeDB.GetTimeSlotForUpdate(ctx, req.TimeSlotID)
		if err != nil {
			errChan <- fmt.Errorf("failed to get time slot: %w", err)
			return
		}
		timeSlot = db.TimeSlot{
			StartTime:      ts.StartTime,
			EndTime:        ts.EndTime,
			MaxPatients:    ts.MaxPatients,
			BookedPatients: ts.BookedPatients,
		}
	}()

	go func() {
		defer wg.Done()
		service, err = s.storeDB.GetServiceByID(ctx, req.ServiceID)
		if err != nil {
			errChan <- fmt.Errorf("failed to get service: %w", err)
			return
		}
	}()

	wg.Wait()

	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	if timeSlot.BookedPatients.Int32 >= timeSlot.MaxPatients.Int32 {
		return nil, fmt.Errorf("time slot is fully booked")
	}
	dateTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	var startTimeFormatted string
	var endTimeFormatted string

	// Create the appointment within a transaction
	var appointment db.Appointment
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {

		startTimeFormatted = time.UnixMicro(timeSlot.StartTime.Microseconds).UTC().Format("15:04:05")
		endTimeFormatted = time.UnixMicro(timeSlot.EndTime.Microseconds).UTC().Format("15:04:05")

		appointment, err = q.CreateAppointment(ctx, db.CreateAppointmentParams{
			DoctorID:          pgtype.Int8{Int64: int64(doctor.ID), Valid: true},
			Petid:             pgtype.Int8{Int64: req.PetID, Valid: true},
			ServiceID:         pgtype.Int8{Int64: service.ID, Valid: true},
			AppointmentReason: pgtype.Text{String: req.Reason, Valid: true},
			Date:              pgtype.Timestamp{Time: dateTime, Valid: true},
			TimeSlotID:        pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
			Username:          pgtype.Text{String: username, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	detail, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, appointment.AppointmentID)

	if err != nil {
		return nil, fmt.Errorf("failed to get appointment detail: %w", err)
	}

	title := fmt.Sprintf("New Appointment for %s", detail.PetName.String)

	notification := AppointmentNotification{
		ID:            fmt.Sprintf("%s-%d-%s-%s", "app", appointment.AppointmentID, detail.PetName.String, doctor.Name),
		Title:         title,
		AppointmentID: appointment.AppointmentID,
		Doctor: Doctor{
			DoctorID:   doctor.ID,
			DoctorName: doctor.Name,
		},
		Pet: Pet{
			PetID:   detail.PetID.Int64,
			PetName: detail.PetName.String,
		},
		Reason: detail.AppointmentReason.String,
		Date:   appointment.Date.Time.Format("2006-01-02"),
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		}, ServiceName: detail.ServiceName.String,
	}

	s.sendAppointmentNotification(ctx, notification)

	// Prepare the response
	response := &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     detail.PetName.String,
		Reason:      detail.AppointmentReason.String,
		Date:        appointment.Date.Time.Format(time.RFC3339),
		ServiceName: detail.ServiceName.String,
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		State:        detail.StateName.String,
		ReminderSend: appointment.ReminderSend.Bool,
		CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		RoomType:     service.Category.String,
	}

	// Clear any cached appointment lists for this user
	if redis.Client != nil {
		redis.Client.RemoveCacheByKey(GetAppointmentListByUserCacheKey(username))
		redis.Client.RemoveCacheByKey(GetAppointmentListByDoctorCacheKey(doctor.ID))
		redis.Client.RemoveCacheByKey(GetTimeSlotsKey(doctor.ID, req.Date))
	}

	return response, nil
}

func (s *AppointmentService) sendAppointmentNotification(ctx context.Context, notification AppointmentNotification) {
	// Send notification via WebSocket
	wsMessage := websocket.WebSocketMessage{
		Type: "appointment_alert",
		Data: notification,
	}

	// Try to send to the doctor
	doctorClientID := fmt.Sprintf("doctor_%d", notification.Doctor.DoctorID)
	// Doctor is offline, store the message for later delivery
	err := s.ws.MessageStore.StoreMessage(ctx, doctorClientID, "", "appointment_alert", notification)
	if err != nil {
		log.Printf("Error storing offline appointment notification for doctor: %v", err)
	}

	s.ws.BroadcastToAll(wsMessage)

	// // Get the user's username from the appointment
	// appointmentDetail, err := s.storeDB.GetAppointmentDetailByAppointmentID(context.Background(), notification.AppointmentID)
	// if err != nil {
	// 	log.Printf("Error getting appointment details for notification: %v", err)
	// 	return
	// }

	// if appointmentDetail.Username.Valid {
	// 	// Send to the pet owner/user
	// 	userClientID := fmt.Sprintf("user_%s", appointmentDetail.Username.String)
	// 	if !s.ws.SendToClient(userClientID, wsMessage) {
	// 		// User is offline, store the message for later delivery
	// 		ctx := context.Background()
	// 		err := s.ws.MessageStore.StoreMessage(ctx, userClientID, appointmentDetail.Username.String, "appointment_alert", notification)
	// 		if err != nil {
	// 			log.Printf("Error storing offline appointment notification for user: %v", err)
	// 		}
	// 	}
	// }
}

func (s *AppointmentService) CreateWalkInAppointment(ctx *gin.Context, req createWalkInAppointmentRequest) (*createAppointmentResponse, error) {
	var err error
	var doctor user.DoctorResponse
	var service db.Service
	var wg sync.WaitGroup
	var username string
	var petID int64

	errChan := make(chan error, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		d, err := s.storeDB.GetDoctor(ctx, req.DoctorID)
		if err != nil {
			errChan <- fmt.Errorf("failed to get doctor: %w", err)
			return
		}
		doctor = user.DoctorResponse{
			ID:             d.ID,
			Specialization: d.Specialization.String,
			Name:           d.Name,
			YearsOfExp:     d.YearsOfExperience.Int32,
			Education:      d.Education.String,
			Certificate:    d.CertificateNumber.String,
			Bio:            d.Bio.String,
		}
	}()

	go func() {
		defer wg.Done()
		service, err = s.storeDB.GetServiceByID(ctx, req.ServiceID)
		if err != nil {
			errChan <- fmt.Errorf("failed to get service: %w", err)
			return
		}
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// Get current time
	now := time.Now()

	// Create the appointment within a transaction
	var appointment db.Appointment
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Case 1: New user with new pet
		if req.Owner != nil {
			// Generate a temporary username based on phone number
			tempUsername := fmt.Sprintf("temp_%s", req.Owner.OwnerPhone)

			// Create new user
			_, err = q.CreateUser(ctx, db.CreateUserParams{
				Username:    tempUsername,
				FullName:    req.Owner.OwnerName,
				PhoneNumber: pgtype.Text{String: req.Owner.OwnerPhone, Valid: true},
				Email:       req.Owner.OwnerEmail,
				Address:     pgtype.Text{String: req.Owner.OwnerAddress, Valid: true},
				Role:        pgtype.Text{String: "customer", Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to create user: %w", err)
			}
			username = tempUsername

			// If pet info is provided, create a new pet
			if req.Pet != nil {
				birthDate, err := time.Parse("2006-01-02", req.Pet.BirthDate)
				if err != nil {
					return fmt.Errorf("invalid birth date format: %w", err)
				}

				pet, err := q.CreatePet(ctx, db.CreatePetParams{
					Username:  username,
					Name:      req.Pet.Name,
					Breed:     pgtype.Text{String: req.Pet.Breed, Valid: true},
					Type:      req.Pet.Species,
					BirthDate: pgtype.Date{Time: birthDate, Valid: true},
					Gender:    pgtype.Text{String: req.Pet.Gender, Valid: true},
					Weight:    pgtype.Float8{Float64: req.Pet.Weight, Valid: true},
					Age:       pgtype.Int4{Int32: int32(req.Pet.Age), Valid: true},
				})
				if err != nil {
					return fmt.Errorf("failed to create pet: %w", err)
				}
				petID = pet.Petid
			} else {
				return fmt.Errorf("pet information is required for new users")
			}
		} else {
			// Case 2: Existing pet, just create appointment
			// Validate that the pet exists
			pet, err := q.GetPetByID(ctx, req.PetID)
			if err != nil {
				return fmt.Errorf("pet not found: %w", err)
			}
			petID = pet.Petid
			username = pet.Username
		}

		appointment, err = q.CreateAppointment(ctx, db.CreateAppointmentParams{
			DoctorID:          pgtype.Int8{Int64: int64(doctor.ID), Valid: true},
			Petid:             pgtype.Int8{Int64: petID, Valid: true},
			ServiceID:         pgtype.Int8{Int64: service.ID, Valid: true},
			AppointmentReason: pgtype.Text{String: req.Reason, Valid: true},
			Date:              pgtype.Timestamp{Time: now, Valid: true},
			TimeSlotID:        pgtype.Int8{Int64: 0, Valid: true}, // No time slot for walk-in
			Username:          pgtype.Text{String: username, Valid: true},
			Priority:          pgtype.Text{String: req.Priority, Valid: true},
			ArrivalTime:       pgtype.Timestamp{Time: now, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	detail, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, appointment.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment detail: %w", err)
	}

	// Prepare the response
	response := &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     detail.PetName.String,
		Reason:      detail.AppointmentReason.String,
		Date:        appointment.Date.Time.Format(time.RFC3339),
		ServiceName: detail.ServiceName.String,
		TimeSlot: timeslot{
			StartTime: now.Format("15:04:05"),
			EndTime:   now.Add(time.Duration(service.Duration.Int16) * time.Minute).Format("15:04:05"),
		},
		State:        detail.StateName.String,
		ReminderSend: appointment.ReminderSend.Bool,
		CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		RoomType:     service.Category.String,
	}

	// Clear any cached appointment lists
	if redis.Client != nil {
		redis.Client.RemoveCacheByKey(GetAppointmentListByUserCacheKey(username))
		redis.Client.RemoveCacheByKey(GetAppointmentListByDoctorCacheKey(doctor.ID))
	}

	return response, nil
}

func (s *AppointmentService) ConfirmPayment(ctx context.Context, appointmentID int64) error {
	// Generate a lock key for this appointment to prevent concurrent confirmation
	lockKey := fmt.Sprintf("appointment:%d", appointmentID)

	// Try to acquire a distributed lock with 5-second timeout
	lockAcquired := AcquireLock(lockKey, 5*time.Second)
	if !lockAcquired {
		return fmt.Errorf("payment confirmation is already in progress, please wait")
	}

	// Make sure the lock is released when we're done
	defer ReleaseLock(lockKey)

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Fetch appointment details
		appointment, err := q.GetAppointmentDetailByAppointmentID(ctx, appointmentID)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}

		// Check appointment state early to avoid unnecessary queries
		state, err := q.GetState(ctx, int64(appointment.StateID.Int32))
		if err != nil {
			return fmt.Errorf("failed to get state: %w", err)
		}
		if state.State == "Confirmed" {
			return errors.New("appointment is already paid")
		}

		// Fetch and validate time slot
		timeSlot, err := q.GetTimeSlotForUpdate(ctx, appointment.TimeSlotID.Int64)
		if err != nil {
			return fmt.Errorf("failed to get time slot: %w", err)
		}
		if timeSlot.BookedPatients.Int32 >= timeSlot.MaxPatients.Int32 {
			return errors.New("time slot is fully booked")
		}

		// Prepare updates in a single transaction
		if err = q.UpdateTimeSlotBookedPatients(ctx, appointment.TimeSlotID.Int64); err != nil {
			return fmt.Errorf("failed to update time slot: %w", err)
		}

		if err = q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			AppointmentID: appointmentID,
			StateID:       pgtype.Int4{Int32: 2, Valid: true}, // Consider defining a constant for StateID
		}); err != nil {
			return fmt.Errorf("failed to update appointment status: %w", err)
		}
		return nil
	})

	// If the payment was successfully confirmed, clear related caches
	if err == nil && redis.Client != nil {
		// Clear the specific appointment cache
		redis.Client.RemoveCacheByKey(GetAppointmentCacheKey(appointmentID))

		// Get the appointment details to clear user and doctor caches
		appDetail, detailErr := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, appointmentID)
		if detailErr == nil {
			// Clear user cache if username exists
			if appDetail.Username.Valid {
				redis.Client.RemoveCacheByKey(GetAppointmentListByUserCacheKey(appDetail.Username.String))
			}

			// Clear doctor cache if doctor ID exists
			if appDetail.DoctorID.Valid {
				redis.Client.RemoveCacheByKey(GetAppointmentListByDoctorCacheKey(appDetail.DoctorID.Int64))
			}
		}
	}

	return err
}

func (s *AppointmentService) CheckInAppoinment(ctx *gin.Context, id, roomID int64, priority string) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Fetch appointment details
		appointment, err := q.GetAppointmentDetailByAppointmentID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}

		req := db.UpdateAppointmentByIDParams{
			AppointmentID:     appointment.AppointmentID,
			StateID:           pgtype.Int4{Int32: 3, Valid: true},
			ArrivalTime:       pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
			RoomID:            pgtype.Int8{Int64: roomID, Valid: true},
			AppointmentReason: appointment.AppointmentReason,
			ReminderSend:      appointment.ReminderSend,
		}

		if priority != "" {
			req.Priority = pgtype.Text{String: priority, Valid: true}
		}

		// Update appointment with room assignment
		if err = q.UpdateAppointmentByID(ctx, req); err != nil {
			return fmt.Errorf("failed to update appointment status: %w", err)
		}

		// Mark room as occupied
		if err = q.AssignRoomToAppointment(ctx, db.AssignRoomToAppointmentParams{
			CurrentAppointmentID: pgtype.Int8{Int64: id, Valid: true},
			ID:                   roomID,
		}); err != nil {
			return fmt.Errorf("failed to assign room: %w", err)
		}

		return nil
	})
}

func (s *AppointmentService) GetAppointmentByID(ctx *gin.Context, id int64) (*Appointment, error) {
	// Check if Redis client is available
	if redis.Client != nil {
		// Generate cache key for this appointment
		cacheKey := GetAppointmentCacheKey(id)

		// Try to get from cache first
		var cachedAppointment Appointment
		err := redis.Client.GetWithBackground(cacheKey, &cachedAppointment)
		if err == nil {
			// Cache hit
			return &cachedAppointment, nil
		}
	}

	// Cache miss or Redis not available - get from database
	detail, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment detail: %w", err)
	}

	// Get doctor info
	doctor, err := s.storeDB.GetDoctor(ctx, detail.DoctorID.Int64)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor: %w", err)
	}

	// Format times
	startTime := time.UnixMicro(detail.StartTime.Microseconds).UTC()
	startTimeFormatted := startTime.Format("15:04:05")

	endTime := time.UnixMicro(detail.EndTime.Microseconds).UTC()
	endTimeFormatted := endTime.Format("15:04:05")

	// Prepare the response
	appointment := &Appointment{
		ID: detail.AppointmentID,
		Pet: Pet{
			PetID:    detail.Petid.Int64,
			PetName:  detail.PetName.String,
			PetBreed: detail.PetBreed.String,
		},
		Owner: Owner{
			OwnerName:    detail.OwnerName.String,
			OwnerPhone:   detail.OwnerPhone.String,
			OwnerEmail:   detail.OwnerEmail.String,
			OwnerAddress: detail.OwnerAddress.String,
		},
		Serivce: Serivce{
			ServiceName:     detail.ServiceName.String,
			ServiceDuration: detail.ServiceDuration.Int16,
			ServiceAmount:   detail.ServiceAmount.Float64,
		},
		Doctor: Doctor{
			DoctorID:   doctor.ID,
			DoctorName: doctor.Name,
		},
		Date: detail.Date.Time.Format("2006-01-02"),
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		Room:         "Examination Room", // Default value for room
		State:        detail.StateName.String,
		Reason:       detail.AppointmentReason.String,
		ReminderSend: detail.ReminderSend.Bool,
		CreatedAt:    detail.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}

	// Store in cache for future requests if Redis is available
	if redis.Client != nil {
		// Cache the appointment for 15 minutes
		err = redis.Client.SetWithBackground(GetAppointmentCacheKey(id), appointment, 15*time.Minute)
		if err != nil {
			// Just log the error, don't return it
			fmt.Printf("Error caching appointment: %v\n", err)
		}
	}

	return appointment, nil
}

func (s *AppointmentService) GetAppointmentsByUser(ctx *gin.Context, username string) ([]Appointment, error) {
	// Check if Redis client is available
	if redis.Client != nil {
		// Generate cache key for this user's appointments
		cacheKey := GetAppointmentListByUserCacheKey(username)

		// Try to get from cache first
		var cachedAppointments []Appointment
		err := redis.Client.GetWithBackground(cacheKey, &cachedAppointments)
		if err == nil {
			// Cache hit
			return cachedAppointments, nil
		}
	}

	// Cache miss or Redis not available - get from database
	results, err := s.storeDB.GetAppointmentsByUser(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}

	var appointments []Appointment
	for _, result := range results {
		// Get doctor info
		doctor, err := s.storeDB.GetDoctor(ctx, result.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("failed to get doctor: %w", err)
		}

		// Format times
		startTime := time.UnixMicro(result.StartTime.Microseconds).UTC()
		startTimeFormatted := startTime.Format("15:04:05")

		endTime := time.UnixMicro(result.EndTime.Microseconds).UTC()
		endTimeFormatted := endTime.Format("15:04:05")

		appointments = append(appointments, Appointment{
			ID: result.AppointmentID,
			Pet: Pet{
				PetID:    result.Petid.Int64,
				PetName:  result.PetName.String,
				PetBreed: result.PetBreed.String,
			},
			Owner: Owner{
				OwnerName:    result.OwnerName.String,
				OwnerPhone:   result.OwnerPhone.String,
				OwnerEmail:   result.OwnerEmail.String,
				OwnerAddress: result.OwnerAddress.String,
			},
			Serivce: Serivce{
				ServiceName:     result.ServiceName.String,
				ServiceDuration: result.ServiceDuration.Int16,
				ServiceAmount:   result.ServiceAmount.Float64,
			},
			Doctor: Doctor{
				DoctorID:   doctor.ID,
				DoctorName: doctor.Name,
			},
			Date: result.Date.Time.Format("2006-01-02"),
			TimeSlot: timeslot{
				StartTime: startTimeFormatted,
				EndTime:   endTimeFormatted,
			},
			Room:         "Examination Room", // Default value
			State:        result.StateName.String,
			Reason:       result.AppointmentReason.String,
			ReminderSend: result.ReminderSend.Bool,
			CreatedAt:    result.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	// Store in cache for future requests if Redis is available
	if redis.Client != nil && len(appointments) > 0 {
		// Cache the appointment list for 5 minutes
		err = redis.Client.SetWithBackground(GetAppointmentListByUserCacheKey(username), appointments, 5*time.Minute)
		if err != nil {
			// Just log the error, don't return it
			fmt.Printf("Error caching appointment list: %v\n", err)
		}
	}

	return appointments, nil
}

func (s *AppointmentService) GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error) {
	// Check if Redis client is available
	if redis.Client != nil {
		// Generate cache key for this doctor's appointments
		cacheKey := GetAppointmentListByDoctorCacheKey(doctorID)

		// Try to get from cache first
		var cachedAppointments []createAppointmentResponse
		err := redis.Client.GetWithBackground(cacheKey, &cachedAppointments)
		if err == nil {
			// Cache hit
			return cachedAppointments, nil
		}
	}

	var response []createAppointmentResponse

	appointments, err := s.storeDB.GetAppointmentsByDoctor(ctx, pgtype.Int8{Int64: doctorID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("Cannot get appointment by doctor")
	}

	for _, appointment := range appointments {
		doc, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("Cannot get doctor")
		}

		response = append(response, createAppointmentResponse{
			ID:          appointment.AppointmentID,
			DoctorName:  doc.Name,
			PetName:     appointment.PetName.String,
			Date:        appointment.Date.Time.Format("2006-01-02"),
			ServiceName: appointment.ServiceName.String,
			TimeSlot: timeslot{
				StartTime: time.UnixMicro(appointment.StartTime.Microseconds).UTC().Format("15:04:05"),
				EndTime:   time.UnixMicro(appointment.EndTime.Microseconds).UTC().Format("15:04:05"),
			},
			State:        appointment.StateName.String,
			ReminderSend: appointment.ReminderSend.Bool,
			CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	// Store in cache for future requests if Redis is available
	if redis.Client != nil && len(response) > 0 {
		// Cache the appointment list for 5 minutes
		err = redis.Client.SetWithBackground(GetAppointmentListByDoctorCacheKey(doctorID), response, 5*time.Minute)
		if err != nil {
			// Just log the error, don't return it
			fmt.Printf("Error caching doctor's appointment list: %v\n", err)
		}
	}

	return response, nil
}

func (s *AppointmentService) GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error) {
	// Check if Redis client is available
	if redis.Client != nil {
		// Generate cache key for these time slots
		cacheKey := GetTimeSlotsKey(doctorID, date)

		// Try to get from cache first
		var cachedTimeSlots []timeSlotResponse
		err := redis.Client.GetWithBackground(cacheKey, &cachedTimeSlots)
		if err == nil {
			// Cache hit
			return cachedTimeSlots, nil
		}
	}

	// Parse date
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format")
	}

	// Lấy danh sách khung giờ của bác sĩ trong ngày cụ thể
	timeSlots, err := s.storeDB.GetTimeSlotsByDoctorAndDate(ctx, db.GetTimeSlotsByDoctorAndDateParams{
		DoctorID: int32(doctorID),
		Date:     pgtype.Date{Time: dateTime, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot get time slots")
	}

	// Lọc ra các khung giờ còn chỗ trống
	var availableTimeSlots []timeSlotResponse
	var slotRes timeSlotResponse
	for _, slot := range timeSlots {

		if slot.BookedPatients.Int32 < slot.MaxPatients.Int32 {
			slotRes = timeSlotResponse{
				ID:        int32(slot.ID),
				StartTime: time.UnixMicro(slot.StartTime.Microseconds).UTC().Format("15:04:05"),
				EndTime:   time.UnixMicro(slot.EndTime.Microseconds).UTC().Format("15:04:05"),
				Status:    "available",
			}
		} else {
			slotRes = timeSlotResponse{
				ID:        int32(slot.ID),
				StartTime: time.UnixMicro(slot.StartTime.Microseconds).UTC().Format("15:04:05"),
				EndTime:   time.UnixMicro(slot.EndTime.Microseconds).UTC().Format("15:04:05"),
				Status:    "full",
			}
		}

		availableTimeSlots = append(availableTimeSlots, slotRes)
	}

	// Store in cache for future requests if Redis is available
	if redis.Client != nil && len(availableTimeSlots) > 0 {
		// Cache the time slots for 5 minutes
		err = redis.Client.SetWithBackground(GetTimeSlotsKey(doctorID, date), availableTimeSlots, 5*time.Minute)
		if err != nil {
			// Just log the error, don't return it
			fmt.Printf("Error caching time slots: %v\n", err)
		}
	}

	return availableTimeSlots, nil
}

func (s *AppointmentService) GetAllAppointments(ctx *gin.Context, date string, option string, pagination *util.Pagination) (*util.PaginationResponse[Appointment], error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	count, err := s.storeDB.CountAllAppointmentsByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("Cannot count appointment")
	}

	appointments, err := s.storeDB.GetAllAppointments(ctx, db.GetAllAppointmentsParams{
		Date:    date,
		Limit:   int32(pagination.PageSize),
		Offset:  int32(offset),
		Column4: option,
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot get appointment")
	}

	var a []Appointment
	for _, appointment := range appointments {

		doc, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("Cannot get doctor")
		}

		a = append(a, Appointment{
			ID: appointment.AppointmentID,
			Pet: Pet{
				PetID:    appointment.PetID.Int64,
				PetName:  appointment.PetName.String,
				PetBreed: appointment.PetBreed.String,
			},
			Room: appointment.RoomName.String,
			Serivce: Serivce{
				ServiceName:     appointment.ServiceName.String,
				ServiceDuration: appointment.ServiceDuration.Int16,
			},
			Doctor: Doctor{
				DoctorID:   appointment.DoctorID.Int64,
				DoctorName: doc.Name,
			},
			Date:         appointment.Date.Time.Format("2006-01-02"),
			State:        appointment.StateName.String,
			ReminderSend: appointment.ReminderSend.Bool,
			CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			TimeSlot: timeslot{
				StartTime: time.UnixMicro(appointment.StartTime.Microseconds).UTC().Format("15:04:05"),
				EndTime:   time.UnixMicro(appointment.EndTime.Microseconds).UTC().Format("15:04:05"),
			},
			Owner: Owner{
				OwnerName:    appointment.OwnerName.String,
				OwnerPhone:   appointment.OwnerPhone.String,
				OwnerEmail:   appointment.OwnerEmail.String,
				OwnerAddress: appointment.OwnerAddress.String,
			},
			Reason:      appointment.AppointmentReason.String,
			ArrivalTime: appointment.ArrivalTime.Time.Format("2006-01-02 15:04:05"),
		})

	}

	response := &util.PaginationResponse[Appointment]{
		Count: count,
		Rows:  &a,
	}
	return response.Build(), nil
}
func (s *AppointmentService) GetAllAppointmentsByDate(ctx *gin.Context, pagination *util.Pagination, date string) ([]Appointment, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

	appointments, err := s.storeDB.GetAllAppointments(ctx, db.GetAllAppointmentsParams{
		Date:    date,
		Limit:   int32(pagination.PageSize),
		Offset:  int32(offset),
		Column4: true,
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot get appointment")
	}

	var a []Appointment
	for _, appointment := range appointments {

		doc, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("Cannot get doctor")
		}

		a = append(a, Appointment{
			ID: appointment.AppointmentID,
			Pet: Pet{
				PetID:    appointment.PetID.Int64,
				PetName:  appointment.PetName.String,
				PetBreed: appointment.PetBreed.String,
			},
			Room: appointment.RoomName.String,
			Serivce: Serivce{
				ServiceName:     appointment.ServiceName.String,
				ServiceDuration: appointment.ServiceDuration.Int16,
			},
			Doctor: Doctor{
				DoctorID:   appointment.DoctorID.Int64,
				DoctorName: doc.Name,
			},
			Date:         appointment.Date.Time.Format("2006-01-02"),
			State:        appointment.StateName.String,
			ReminderSend: appointment.ReminderSend.Bool,
			CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			TimeSlot: timeslot{
				StartTime: time.UnixMicro(appointment.StartTime.Microseconds).UTC().Format("15:04:05"),
				EndTime:   time.UnixMicro(appointment.EndTime.Microseconds).UTC().Format("15:04:05"),
			},
			Owner: Owner{
				OwnerName:    appointment.OwnerName.String,
				OwnerPhone:   appointment.OwnerPhone.String,
				OwnerEmail:   appointment.OwnerEmail.String,
				OwnerAddress: appointment.OwnerAddress.String,
			},
			Reason: appointment.AppointmentReason.String,
		})

	}
	return a, nil
}

func (s *AppointmentService) CreateSOAPService(ctx *gin.Context, soap CreateSOAPRequest, appointmentID int64) (*SOAPResponse, error) {

	var err error
	var consultation db.Consultation
	var soapResponse SOAPResponse

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		consultation, err = q.CreateSOAP(ctx, db.CreateSOAPParams{
			AppointmentID: pgtype.Int8{Int64: appointmentID, Valid: true},
			Subjective:    pgtype.Text{String: soap.Subjective, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("Cannot create SOAP")
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot create SOAP")
	}

	return &SOAPResponse{
		ConsultationID: int64(consultation.ID),
		AppointmentID:  consultation.AppointmentID.Int64,
		Subjective:     consultation.Subjective.String,
		Objective:      soapResponse.Objective,
		Assessment:     consultation.Assessment.String,
	}, nil
}

func (s *AppointmentService) UpdateSOAPService(ctx *gin.Context, soap UpdateSOAPRequest, appointmentID int64) (*SOAPResponse, error) {
	var consultation db.Consultation
	var soapResponse SOAPResponse
	var soapRequest db.UpdateSOAPParams // Fixed typo in variable name

	existingConsultation, err := s.storeDB.GetSOAPByAppointmentID(ctx, pgtype.Int8{Int64: appointmentID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get consultation: %w", err)
	}

	// Set AppointmentID in the request
	soapRequest.AppointmentID = pgtype.Int8{Int64: appointmentID, Valid: true}

	// Handle Subjective field
	soapRequest.Subjective = pgtype.Text{
		String: soap.Subjective,
		Valid:  soap.Subjective != "",
	}
	if !soapRequest.Subjective.Valid {
		soapRequest.Subjective = existingConsultation.Subjective
	}

	// Handle Assessment field
	soapRequest.Assessment = pgtype.Text{
		String: soap.Assessment,
		Valid:  soap.Assessment != "",
	}
	if !soapRequest.Assessment.Valid {
		soapRequest.Assessment = existingConsultation.Assessment
	}

	// Handle Objective field - check if it's empty using a type-specific check
	if !isObjectiveEmpty(soap.Objective) {
		objectiveJSON, err := json.Marshal(soap.Objective)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal objective data: %w", err)
		}
		soapRequest.Objective = objectiveJSON
	} else {
		soapRequest.Objective = existingConsultation.Objective
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		consultation, err = q.UpdateSOAP(ctx, soapRequest)
		if err != nil {
			return fmt.Errorf("failed to update SOAP: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(consultation.Objective, &soapResponse.Objective); err != nil {
		return nil, fmt.Errorf("failed to unmarshal objective data: %w", err)
	}

	return &SOAPResponse{
		ConsultationID: int64(consultation.ID),
		AppointmentID:  consultation.AppointmentID.Int64,
		Subjective:     consultation.Subjective.String,
		Objective:      soapResponse.Objective,
		Assessment:     consultation.Assessment.String,
	}, nil
}

func (s *AppointmentService) UpdateAppointmentService(ctx *gin.Context, req updateAppointmentRequest, appointmentID int64) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		var err error

		existingAppointment, err := q.GetAppointmentDetailByAppointmentID(ctx, appointmentID)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}
		updateReq := db.UpdateAppointmentByIDParams{
			AppointmentID: appointmentID,
		}
		if req.StateID != nil {
			updateReq.StateID = pgtype.Int4{Int32: int32(*req.StateID), Valid: true}
		} else {
			updateReq.StateID = pgtype.Int4{Int32: int32(existingAppointment.DoctorID.Int64), Valid: true}
		}

		if req.RoomID != nil {
			room, err := s.storeDB.GetRoomByID(ctx, int64(*req.RoomID))
			if err != nil {
				return fmt.Errorf("failed to get room: %w", err)
			}
			if room.Status.String != "available" {
				return fmt.Errorf("room is not available")
			}
			updateReq.RoomID = pgtype.Int8{Int64: int64(*req.RoomID), Valid: true}
		} else {
			updateReq.RoomID = pgtype.Int8{Int64: int64(existingAppointment.RoomID.Int64), Valid: true}
		}

		if req.AppointmentReason != nil {
			updateReq.AppointmentReason = pgtype.Text{String: *req.AppointmentReason, Valid: true}
		} else {
			updateReq.AppointmentReason = pgtype.Text{String: existingAppointment.AppointmentReason.String, Valid: true}
		}

		if req.ReminderSend != nil {
			updateReq.ReminderSend = pgtype.Bool{Bool: *req.ReminderSend, Valid: true}
		} else {
			updateReq.ReminderSend = pgtype.Bool{Bool: existingAppointment.ReminderSend.Bool, Valid: true}
		}

		if req.ArrivalTime != nil {
			arrivalTime, err := time.Parse("2006-01-02 15:04:05", *req.ArrivalTime)
			if err != nil {
				return fmt.Errorf("invalid arrival time format: %w", err)
			}
			updateReq.ArrivalTime = pgtype.Timestamp{Time: arrivalTime, Valid: true}
		} else {
			updateReq.ArrivalTime = pgtype.Timestamp{Time: existingAppointment.ArrivalTime.Time, Valid: true}
		}

		if req.Priority != nil {
			updateReq.Priority = pgtype.Text{String: *req.Priority, Valid: true}
		} else {
			updateReq.Priority = pgtype.Text{String: existingAppointment.Priority.String, Valid: true}
		}

		updateReq.ReminderSend = pgtype.Bool{Bool: existingAppointment.ReminderSend.Bool, Valid: true}

		err = q.UpdateAppointmentByID(ctx, updateReq)
		if err != nil {
			return fmt.Errorf("Cannot update appointment %w", err)
		}
		return nil
	})
}

func (s *AppointmentService) GetQueueService(ctx *gin.Context, username string) ([]QueueItem, error) {
	doctor, err := s.storeDB.GetDoctorByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctor: %v", err)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	// dateTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	appointments, err := s.storeDB.GetAppointmentsQueue(ctx, db.GetAppointmentsQueueParams{
		DoctorID: pgtype.Int8{Int64: doctor.ID, Valid: true},
		Date:     now,
	}) // Assuming 1 is the state ID for waiting/arrived
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %v", err)
	}

	var queueItems []QueueItem
	for _, appointment := range appointments {

		a, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, appointment.AppointmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get appointment detail: %v", err)
		}

		doc, err := s.storeDB.GetDoctor(ctx, a.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("failed to get doctor: %v", err)
		}
		// Calculate waiting time
		waitingSince := appointment.ArrivalTime.Time
		actualWaitTime := time.Since(waitingSince).Round(time.Minute).String()

		queueItem := QueueItem{
			ID:              appointment.AppointmentID,
			PatientName:     a.PetName.String,
			Status:          a.StateName.String, // You might want to get this from a separate status field
			AppointmentType: a.ServiceName.String,
			Doctor:          doc.Name,
			Priority:        a.Priority.String,
			WaitingSince:    waitingSince.String(),
			ActualWaitTime:  actualWaitTime,
		}

		queueItems = append(queueItems, queueItem)
	}

	// Sort queue items by priority (high first) and waiting time
	sortQueueItems(queueItems)

	return queueItems, nil
}

func (s *AppointmentService) UpdateQueueItemStatusService(ctx *gin.Context, id int64, status string) error {
	// Update appointment status
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// You might need to add a new query to update appointment status
		// For now, we'll just update the state_id
		stateID := int32(1) // You should map status to appropriate state_id
		if status == "in_progress" {
			stateID = 2
		} else if status == "completed" {
			stateID = 3
		}

		err := q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			AppointmentID: id,
			StateID:       pgtype.Int4{Int32: stateID, Valid: true},
		})
		return err
	})

	if err != nil {
		return fmt.Errorf("failed to update appointment status: %v", err)
	}

	return nil
}

func (s *AppointmentService) GetSOAPByAppointmentID(ctx *gin.Context, appointmentID int64) (*SOAPResponse, error) {
	soap, err := s.storeDB.GetSOAPByAppointmentID(ctx, pgtype.Int8{Int64: appointmentID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get SOAP: %v", err)
	}

	var soapResponse SOAPResponse
	if err := json.Unmarshal(soap.Objective, &soapResponse.Objective); err != nil {
		return nil, fmt.Errorf("failed to unmarshal objective data: %w", err)
	}

	return &SOAPResponse{
		ConsultationID: int64(soap.ID),
		AppointmentID:  soap.AppointmentID.Int64,
		Subjective:     soap.Subjective.String,
		Objective:      soapResponse.Objective,
		Assessment:     soap.Assessment.String,
	}, nil
}

func (s *AppointmentService) GetHistoryAppointmentsByPetID(ctx *gin.Context, petID int64) ([]historyAppointmentResponse, error) {

	appointments, err := s.storeDB.GetHistoryAppointmentsByPetID(ctx, pgtype.Int8{Int64: petID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %v", err)
	}

	var a []historyAppointmentResponse
	for _, appointment := range appointments {
		doctor, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("failed to get doctor: %w", err)
		}
		status, err := s.storeDB.GetState(ctx, int64(appointment.StateID.Int32))
		if err != nil {
			return nil, fmt.Errorf("failed to get status: %w", err)
		}
		a = append(a, historyAppointmentResponse{
			ID:          appointment.AppointmentID,
			Reason:      appointment.AppointmentReason.String,
			Date:        appointment.Date.Time.Format("2006-01-02"),
			ServiceName: appointment.ServiceName.String,
			ArrivalTime: appointment.ArrivalTime.Time.Format("2006-01-02 15:04:05"),
			DoctorName:  doctor.Name,
			Status:      status.State,
			Room:        appointment.RoomName.String,
		})
	}
	return a, nil
}

func (s *AppointmentService) GetAppointmentDistributionByService(ctx *gin.Context, startDate, endDate string) ([]AppointmentDistribution, error) {
	// Parse the date strings
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// Get distribution data from database
	rows, err := s.storeDB.GetAppointmentDistribution(ctx, db.GetAppointmentDistributionParams{
		Column1: pgtype.Date{Time: start, Valid: true},
		Column2: pgtype.Date{Time: end, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment distribution: %w", err)
	}

	// Now you can iterate over the rows
	var distributions []AppointmentDistribution
	for _, row := range rows {
		// Convert pgtype.Numeric to float64
		var percentage float64
		if row.Percentage.Valid {
			// Get the numeric value and convert it properly
			numValue, err := row.Percentage.Float64Value()
			if err == nil {
				percentage = float64(numValue.Float64)
			}
		}

		distributions = append(distributions, AppointmentDistribution{
			ServiceID:        row.ServiceID,
			ServiceName:      row.ServiceName.String,
			AppointmentCount: row.AppointmentCount,
			Percentage:       percentage,
		})
	}

	return distributions, nil
}

func (s *AppointmentService) HandleWebSocket(ctx *gin.Context) {
	s.ws.HandleWebSocket(ctx)
}

// GetPendingNotifications retrieves any pending notifications for a client
func (s *AppointmentService) GetPendingNotifications(ctx *gin.Context) ([]websocket.OfflineMessage, error) {
	// Get pending notifications from the WebSocket message store
	pendingMessages, err := s.ws.MessageStore.GetPendingMessages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending messages: %w", err)
	}

	// Create a map of existing message IDs for quick lookup
	existingMsgIds := make(map[string]bool)
	for _, msg := range pendingMessages {
		// Create a key based on message content to avoid duplicates
		key := fmt.Sprintf("%s-%s", msg.MessageType, string(msg.Data))
		existingMsgIds[key] = true
	}

	return pendingMessages, nil
}

func (s *AppointmentService) MarkMessageDelivered(ctx *gin.Context, id int64) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.MarkMessageDelivered(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to mark message as delivered: %w", err)
		}
		return nil
	})
}
