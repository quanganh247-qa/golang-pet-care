package appointment

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type AppointmentServiceInterface interface {
	CreateSOAPService(ctx *gin.Context, soap CreateSOAPRequest, appointmentID int64) (*SOAPResponse, error)
	UpdateSOAPService(ctx *gin.Context, soap UpdateSOAPRequest, appointmentID int64) (*SOAPResponse, error)
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error)
	ConfirmPayment(ctx context.Context, appointmentID int64) error
	GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error)
	GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error)
	GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error)
	GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error)
	GetAllAppointments(ctx *gin.Context) ([]createAppointmentResponse, error)
}

func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error) {

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
			DoctorID:   pgtype.Int8{Int64: int64(doctor.ID), Valid: true},
			Petid:      pgtype.Int8{Int64: req.PetID, Valid: true},
			ServiceID:  pgtype.Int8{Int64: service.ID, Valid: true},
			Date:       pgtype.Timestamp{Time: dateTime, Valid: true},
			TimeSlotID: pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
			Username:   pgtype.Text{String: username, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

		// Cập nhật khung giờ
		err = q.UpdateTimeSlotBookedPatients(ctx, req.TimeSlotID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	detail, err := s.storeDB.GetAppointmentDetail(ctx, db.GetAppointmentDetailParams{
		ID:    appointment.ServiceID.Int64,
		Petid: appointment.Petid.Int64,
		ID_2:  int64(appointment.StateID.Int32),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment detail: %w", err)
	}

	// Prepare the response
	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     detail.PetName,
		Date:        appointment.Date.Time.Format(time.RFC3339),
		ServiceName: detail.ServiceName.String,
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		State:        detail.StateName,
		Notes:        appointment.Notes.String,
		ReminderSend: appointment.ReminderSend.Bool,
		CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *AppointmentService) ConfirmPayment(ctx context.Context, appointmentID int64) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		appointment, err := q.GetAppointmentDetailByAppointmentID(ctx, appointmentID)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}

		state, err := q.GetState(ctx, appointment.StateID.Int64)
		if err != nil {
			return fmt.Errorf("failed to get state: %w", err)
		}
		if state.State == "Confirmed" {
			return fmt.Errorf("appointment is already paid")
		}
		timeSlot, err := q.GetTimeSlotForUpdate(ctx, appointment.TimeSlotID.Int64)
		if err != nil {
			return fmt.Errorf("failed to get time slot: %w", err)
		}

		if timeSlot.BookedPatients.Int32 >= timeSlot.MaxPatients.Int32 {
			return fmt.Errorf("time slot is fully booked")
		}

		err = q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			AppointmentID: appointmentID,
			StateID:       pgtype.Int4{Int32: 2, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		if err = q.UpdateTimeSlotBookedPatients(ctx, appointment.TimeSlotID.Int64); err != nil {
			return fmt.Errorf("failed to update time slot: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (s *AppointmentService) GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error) {
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
	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     appointment.PetName.String,
		Date:        appointment.Date.Time.Format(time.RFC3339),
		ServiceName: appointment.ServiceName.String,
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		Notes:        appointment.Notes.String,
		ReminderSend: appointment.ReminderSend.Bool,
		State:        appointment.StateName.String,
		CreatedAt:    appointment.CreatedAt.Time.Format(time.RFC3339),
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
			Notes:        appointment.Notes.String,
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

func (s *AppointmentService) GetAllAppointments(ctx *gin.Context) ([]createAppointmentResponse, error) {
	appointments, err := s.storeDB.GetAllAppointments(ctx)
	if err != nil {
		return nil, fmt.Errorf("Cannot get appointment")
	}

	var a []createAppointmentResponse
	for _, appointment := range appointments {

		doc, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
		if err != nil {
			return nil, fmt.Errorf("Cannot get doctor")
		}

		a = append(a, createAppointmentResponse{
			ID:           appointment.AppointmentID,
			DoctorName:   doc.Name,
			PetName:      appointment.PetName.String,
			ServiceName:  appointment.ServiceName.String,
			Date:         appointment.Date.Time.Format("2006-01-02"),
			State:        appointment.StateName.String,
			Notes:        appointment.Notes.String,
			ReminderSend: appointment.ReminderSend.Bool,
			CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			TimeSlot: timeslot{
				StartTime: time.UnixMicro(appointment.StartTime.Microseconds).UTC().Format("15:04:05"),
				EndTime:   time.UnixMicro(appointment.EndTime.Microseconds).UTC().Format("15:04:05"),
			},
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
			Objective:     pgtype.Text{String: soap.Objective, Valid: true},
			Assessment:    pgtype.Text{String: soap.Assessment, Valid: true},
			Plan:          pgtype.Text{String: soap.Plan, Valid: true},
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
		Plan:           consultation.Plan.String,
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
			Plan:          pgtype.Text{String: soap.Plan, Valid: true},
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
		Plan:           consultation.Plan.String,
	}, nil
}
