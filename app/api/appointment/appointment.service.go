package appointment

import (
	"context"
	"errors"
	"fmt"
<<<<<<< HEAD
	"sync"
=======
	"strconv"
>>>>>>> cfbe865 (updated service response)
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type AppointmentServiceInterface interface {
<<<<<<< HEAD
	CreateSOAPService(ctx *gin.Context, soap CreateSOAPRequest, appointmentID int64) (*SOAPResponse, error)
	UpdateSOAPService(ctx *gin.Context, soap UpdateSOAPRequest, appointmentID int64) (*SOAPResponse, error)
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error)
	ConfirmPayment(ctx context.Context, appointmentID int64) error
	CheckInAppoinment(ctx *gin.Context, id, roomID int64, priority string) error
	GetAppointmentByID(ctx *gin.Context, id int64) (*Appointment, error)
	GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error)
	GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error)
	GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error)
	GetAllAppointments(ctx *gin.Context, date string, option string, pagination *util.Pagination) ([]Appointment, error)
	GetAllAppointmentsByDate(ctx *gin.Context, pagination *util.Pagination, date string) ([]Appointment, error)
	UpdateAppointmentService(ctx *gin.Context, req updateAppointmentRequest, appointmentID int64) error
	UpdateQueueItemStatusService(ctx *gin.Context, id int64, status string) error
	GetQueueService(ctx *gin.Context, username string) ([]QueueItem, error)
	GetHistoryAppointmentsByPetID(ctx *gin.Context, petID int64) ([]historyAppointmentResponse, error)
	GetSOAPByAppointmentID(ctx *gin.Context, appointmentID int64) (*SOAPResponse, error)
=======
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error)
	UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error
	GetAppointmentsOfDoctorService(ctx *gin.Context, doctorID int64) ([]AppointmentWithDetails, error)
	GetAppointmentByID(ctx *gin.Context, id int64) (*db.Appointment, error)
<<<<<<< HEAD
>>>>>>> 7e35c2e (get appointment detail)
=======
	GetAppointmentsByPetOfUser(ctx *gin.Context, username string) ([]AppointmentWithDetails, error)
>>>>>>> e30b070 (Get list appoinment by user)
}

func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error) {

<<<<<<< HEAD
	var err error
	var timeSlot db.TimeSlot
	var doctor user.DoctorResponse
	var service db.Service
	var wg sync.WaitGroup

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
=======
	var arg db.CreateAppointmentParams

	// convert string to int64
	doctorID, err := strconv.ParseInt(req.DoctorID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error while converting doctor id: %w", err)
	}

	doctor, err := s.storeDB.GetDoctor(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("error while getting doctor: %w", err)
>>>>>>> cfbe865 (updated service response)
	}
	dateTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	var startTimeFormatted string
	var endTimeFormatted string

<<<<<<< HEAD
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
=======
	if req.Date != "" {
		//convert string to time.TIme
		dateTime, err := time.Parse("2006-01-02T15:04:05Z", req.Date)
		if err != nil {
			return nil, fmt.Errorf("error while converting date: %w", err)
		}
		arg.Date = pgtype.Timestamp{Time: dateTime, Valid: true}
	}
	arg.DoctorID = pgtype.Int8{Int64: doctor.ID, Valid: true}
	arg.Petid = pgtype.Int8{Int64: req.PetID, Valid: true}
	arg.ServiceID = pgtype.Int8{Int64: req.ServiceID, Valid: true}

	appointment, err := s.storeDB.CreateAppointment(ctx, arg)
>>>>>>> cfbe865 (updated service response)
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	detail, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, appointment.AppointmentID)

	if err != nil {
<<<<<<< HEAD
		return nil, fmt.Errorf("failed to get appointment detail: %w", err)
	}

	// Prepare the response
	return &createAppointmentResponse{
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
=======
		return nil, fmt.Errorf("error while getting service: %w", err)
	}

	pet, err := s.storeDB.GetPetByID(ctx, appointment.Petid.Int64)
	if err != nil {
		return nil, fmt.Errorf("error while getting pet: %w", err)
	}

	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     pet.Name,
		ServiceName: service.Name,
		Note:        req.Note,
>>>>>>> cfbe865 (updated service response)
	}, nil
}
func (s *AppointmentService) ConfirmPayment(ctx context.Context, appointmentID int64) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Fetch appointment details
		appointment, err := q.GetAppointmentDetailByAppointmentID(ctx, appointmentID)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}

		// Check appointment state early to avoid unnecessary queries
		state, err := q.GetState(ctx, appointment.StateID.Int64)
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
	var err error

	appointment, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Cannot get appointment detail")
	}

	doctor, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
	if err != nil {
		return nil, fmt.Errorf("Cannot get doctor detail")
	}

	// Format start and end times
	startTime := time.UnixMicro(appointment.StartTime.Microseconds).UTC()
	startTimeFormatted := startTime.Format("15:04:05")

	endTime := time.UnixMicro(appointment.EndTime.Microseconds).UTC()
	endTimeFormatted := endTime.Format("15:04:05")

	// Prepare the response
	return &Appointment{
		ID: appointment.AppointmentID,
		Doctor: Doctor{
			DoctorID:   doctor.ID,
			DoctorName: doctor.Name,
		},
		Owner: Owner{
			OwnerName:    appointment.OwnerName.String,
			OwnerPhone:   appointment.OwnerPhone.String,
			OwnerEmail:   appointment.OwnerEmail.String,
			OwnerAddress: appointment.OwnerAddress.String,
		},
		Pet: Pet{
			PetID:    appointment.PetID.Int64,
			PetName:  appointment.PetName.String,
			PetBreed: appointment.PetBreed.String,
		},
		Date: appointment.Date.Time.Format("2006-01-02"),
		Serivce: Serivce{
			ServiceName:     appointment.ServiceName.String,
			ServiceDuration: appointment.ServiceDuration.Int16,
		},
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		Reason:       appointment.AppointmentReason.String,
		ReminderSend: appointment.ReminderSend.Bool,
		State:        appointment.StateName.String,
		CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *AppointmentService) GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error) {

	var a []createAppointmentResponse
	appointments, err := s.storeDB.GetAppointmentsByUser(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("Cannot get appointment by user")
	}
	for _, appointment := range appointments {
		doc, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("failed to get doctor: %w", err)
		}
		a = append(a, createAppointmentResponse{
			ID:          appointment.AppointmentID,
			DoctorName:  doc.Name,
			PetName:     appointment.PetName.String,
			Date:        appointment.Date.Time.Format("2006-01-02"),
			ServiceName: appointment.ServiceName.String,
			TimeSlot: timeslot{
				StartTime: time.UnixMicro(appointment.StartTime.Microseconds).UTC().Format("15:04:05"),
				EndTime:   time.UnixMicro(appointment.EndTime.Microseconds).UTC().Format("15:04:05"),
			},
			State:     appointment.State.String,
			CreatedAt: appointment.CreatedAt.Time.Format(time.RFC3339),
		})
	}
	return a, nil
}

func (s *AppointmentService) GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error) {

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

	return response, nil
}

func (s *AppointmentService) GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error) {
	// Parse ngày
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

	return availableTimeSlots, nil
}

func (s *AppointmentService) GetAllAppointments(ctx *gin.Context, date string, option string, pagination *util.Pagination) ([]Appointment, error) {
	offset := (pagination.Page - 1) * pagination.PageSize

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
			Reason: appointment.AppointmentReason.String,
		})

	}
	return a, nil
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
		Objective:      consultation.Objective.String,
		Assessment:     consultation.Assessment.String,
		// Plan:           consultation.Plan.Int64,
	}, nil
}

func (s *AppointmentService) UpdateSOAPService(ctx *gin.Context, soap UpdateSOAPRequest, appointmentID int64) (*SOAPResponse, error) {

	var err error
	var consultation db.Consultation
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		consultation, err = q.UpdateSOAP(ctx, db.UpdateSOAPParams{
			AppointmentID: pgtype.Int8{Int64: appointmentID, Valid: true},
			Subjective:    pgtype.Text{String: soap.Subjective, Valid: true},
			Objective:     pgtype.Text{String: soap.Objective, Valid: true},
			Assessment:    pgtype.Text{String: soap.Assessment, Valid: true},
			// Plan:          pgtype.Text{String: soap.Plan, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("Cannot update SOAP")
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Cannot update SOAP")
	}
	return &SOAPResponse{
		ConsultationID: int64(consultation.ID),
		AppointmentID:  consultation.AppointmentID.Int64,
		Subjective:     consultation.Subjective.String,
		Objective:      consultation.Objective.String,
		Assessment:     consultation.Assessment.String,
		// Plan:           consultation.Plan.String,
	}, nil
}

func (s *AppointmentService) UpdateAppointmentService(ctx *gin.Context, req updateAppointmentRequest, appointmentID int64) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		var err error
		updateReq := db.UpdateAppointmentByIDParams{
			AppointmentID: appointmentID,
		}
		if req.StateID != nil {
			updateReq.StateID = pgtype.Int4{Int32: int32(*req.StateID), Valid: true}
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
		}

		if req.AppointmentReason != nil {
			updateReq.AppointmentReason = pgtype.Text{String: *req.AppointmentReason, Valid: true}
		}

		if req.ReminderSend != nil {
			updateReq.ReminderSend = pgtype.Bool{Bool: *req.ReminderSend, Valid: true}
		}

		if req.ArrivalTime != nil {
			arrivalTime, err := time.Parse("2006-01-02 15:04:05", *req.ArrivalTime)
			if err != nil {
				return fmt.Errorf("invalid arrival time format: %w", err)
			}
			updateReq.ArrivalTime = pgtype.Timestamp{Time: arrivalTime, Valid: true}
		}

		if req.Priority != nil {
			updateReq.Priority = pgtype.Text{String: *req.Priority, Valid: true}
		}

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

	appointments, err := s.storeDB.GetAppointmentsQueue(ctx, pgtype.Int8{Int64: doctor.ID, Valid: true}) // Assuming 1 is the state ID for waiting/arrived
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
			WaitingSince:    waitingSince.Format("3:04 PM"),
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

	return &SOAPResponse{
		ConsultationID: int64(soap.ID),
		AppointmentID:  soap.AppointmentID.Int64,
		Subjective:     soap.Subjective.String,
		Objective:      soap.Objective.String,
		Assessment:     soap.Assessment.String,
		// Plan:           soap.Plan.String,
	}, nil
}

<<<<<<< HEAD
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
		a = append(a, historyAppointmentResponse{
			ID:          appointment.AppointmentID,
			Reason:      appointment.AppointmentReason.String,
			Date:        appointment.Date.Time.Format("2006-01-02"),
			ServiceName: appointment.ServiceName.String,
			ArrivalTime: appointment.ArrivalTime.Time.Format("2006-01-02 15:04:05"),
			DoctorName:  doctor.Name,
			Room:        appointment.RoomName.String,
		})
	}
	return a, nil
=======
// get by id
func (s *AppointmentService) GetAppointmentByID(ctx *gin.Context, id int64) (*db.Appointment, error) {
	appointment, err := s.storeDB.GetAppointmentDetailById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting appointment by id: %w", err)
	}
	return &appointment, nil
>>>>>>> 7e35c2e (get appointment detail)
}

// get by id
func (s *AppointmentService) GetAppointmentsByPetOfUser(ctx *gin.Context, username string) ([]AppointmentWithDetails, error) {
	rows, err := s.storeDB.GetAppointmentsByPetOfUser(ctx, username)
	if err != nil {
		return nil, err
	}
	var a []AppointmentWithDetails
	for _, row := range rows {
		service, err := s.storeDB.GetServiceByID(ctx, row.ServiceID.Int64)
		if err != nil {
			return nil, err
		}
		pet, err := s.storeDB.GetPetByID(ctx, row.Petid.Int64)
		if err != nil {
			return nil, err
		}
		a = append(a, AppointmentWithDetails{
			AppointmentID: row.AppointmentID,
			PetName:       pet.Name,
			ServiceName:   service.Name,
			Date:          row.Date.Time.Format(time.RFC3339),
			Status:        row.Status.String,
			CreatedAt:     row.CreatedAt.Time.Format(time.RFC3339),
		})
	}
	return a, nil
}
