package appointment

import (
	"context"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"errors"
=======
>>>>>>> b393bb9 (add service and add permission)
=======
	"errors"
>>>>>>> 4ccd381 (Update appointment flow)
=======
>>>>>>> b393bb9 (add service and add permission)
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"sync"
=======
	"strconv"
>>>>>>> cfbe865 (updated service response)
=======
	"log"
=======
>>>>>>> dc47646 (Optimize SQL query)
	"sync"
>>>>>>> 685da65 (latest update)
=======
	"strconv"
>>>>>>> cfbe865 (updated service response)
=======
	"log"
=======
>>>>>>> dc47646 (Optimize SQL query)
	"sync"
>>>>>>> 685da65 (latest update)
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type AppointmentServiceInterface interface {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> e859654 (Elastic search)
	CreateSOAPService(ctx *gin.Context, soap CreateSOAPRequest, appointmentID int64) (*SOAPResponse, error)
	UpdateSOAPService(ctx *gin.Context, soap UpdateSOAPRequest, appointmentID int64) (*SOAPResponse, error)
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error)
	ConfirmPayment(ctx context.Context, appointmentID int64) error
	CheckInAppoinment(ctx *gin.Context, id, roomID int64, priority string) error
<<<<<<< HEAD
<<<<<<< HEAD
	GetAppointmentByID(ctx *gin.Context, id int64) (*Appointment, error)
=======
	GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error)
>>>>>>> 4ccd381 (Update appointment flow)
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
<<<<<<< HEAD
<<<<<<< HEAD
	GetAppointmentsOfDoctorService(ctx *gin.Context, doctorID int64) ([]AppointmentWithDetails, error)
	GetAppointmentByID(ctx *gin.Context, id int64) (*db.Appointment, error)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 7e35c2e (get appointment detail)
=======
	GetAppointmentsByPetOfUser(ctx *gin.Context, username string) ([]AppointmentWithDetails, error)
>>>>>>> e30b070 (Get list appoinment by user)
=======
	GetAppointmentByID(ctx *gin.Context, id int64) (*Appointment, error)
	GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error)
	GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error)
	GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error)
	GetAllAppointments(ctx *gin.Context, date string, option string, pagination *util.Pagination) ([]Appointment, error)
	GetAllAppointmentsByDate(ctx *gin.Context, pagination *util.Pagination, date string) ([]Appointment, error)
	UpdateAppointmentService(ctx *gin.Context, req updateAppointmentRequest, appointmentID int64) error
<<<<<<< HEAD
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
	UpdateQueueItemStatusService(ctx *gin.Context, id int64, status string) error
	GetQueueService(ctx *gin.Context, username string) ([]QueueItem, error)
	GetHistoryAppointmentsByPetID(ctx *gin.Context, petID int64) ([]historyAppointmentResponse, error)
	GetSOAPByAppointmentID(ctx *gin.Context, appointmentID int64) (*SOAPResponse, error)
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
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
=======
=======
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error)
	ConfirmPayment(ctx context.Context, appointmentID int64) error
>>>>>>> b393bb9 (add service and add permission)
	GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error)
	GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error)
	GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error)
	GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error)
	GetAllAppointments(ctx *gin.Context) ([]createAppointmentResponse, error)
}

<<<<<<< HEAD
func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error) {
<<<<<<< HEAD
	// Validate input
	if req.DoctorID == 0 || req.PetID == 0 || req.ServiceID == 0 || req.TimeSlotID == 0 || req.Date == "" {
		return nil, fmt.Errorf("missing required fields: doctor_id, pet_id, service_id, time_slot_id, or date")
>>>>>>> 685da65 (latest update)
=======

<<<<<<< HEAD
	var arg db.CreateAppointmentParams

	// convert string to int64
	doctorID, err := strconv.ParseInt(req.DoctorID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error while converting doctor id: %w", err)
>>>>>>> cfbe865 (updated service response)
=======
=======
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error)
	ConfirmPayment(ctx context.Context, appointmentID int64) error
>>>>>>> b393bb9 (add service and add permission)
	GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error)
	GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error)
	GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error)
	GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error)
	GetAllAppointments(ctx *gin.Context) ([]createAppointmentResponse, error)
}

<<<<<<< HEAD
func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error) {
	// Validate input
	if req.DoctorID == 0 || req.PetID == 0 || req.ServiceID == 0 || req.TimeSlotID == 0 || req.Date == "" {
		return nil, fmt.Errorf("missing required fields: doctor_id, pet_id, service_id, time_slot_id, or date")
>>>>>>> 685da65 (latest update)
	}
=======
func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error) {
>>>>>>> b393bb9 (add service and add permission)
=======
func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error) {
>>>>>>> b393bb9 (add service and add permission)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	// Fetch doctor details
	doctor, err := s.storeDB.GetDoctor(ctx, req.DoctorID)
=======
	doctor, err := s.storeDB.GetDoctor(ctx, doctorID)
>>>>>>> cfbe865 (updated service response)
	if err != nil {
<<<<<<< HEAD
		return nil, fmt.Errorf("error while getting doctor: %w", err)
>>>>>>> cfbe865 (updated service response)
=======
		return nil, fmt.Errorf("failed to fetch doctor: %w", err)
>>>>>>> 685da65 (latest update)
=======
=======
>>>>>>> dc47646 (Optimize SQL query)
	var err error
	var timeSlot db.TimeSlot
	var doctor user.DoctorResponse
	var service db.Service
	var wg sync.WaitGroup

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> dc47646 (Optimize SQL query)
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
<<<<<<< HEAD
>>>>>>> dc47646 (Optimize SQL query)
	}
	dateTime, err := time.Parse("2006-01-02", req.Date)
=======
	if req.Date != "" {
		//convert string to time.TIme
		dateTime, err := time.Parse("2006-01-02T15:04:05Z", req.Date)
=======
	// Fetch doctor details
	doctor, err := s.storeDB.GetDoctor(ctx, req.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch doctor: %w", err)
=======
>>>>>>> dc47646 (Optimize SQL query)
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
<<<<<<< HEAD
		// Create the appointment
		appointment, err = q.CreateAppointment(ctx, arg)
>>>>>>> 685da65 (latest update)
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

<<<<<<< HEAD
	appointment, err := s.storeDB.CreateAppointment(ctx, arg)
>>>>>>> cfbe865 (updated service response)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	var startTimeFormatted string
	var endTimeFormatted string

<<<<<<< HEAD
<<<<<<< HEAD
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
=======
	// Parse and validate the date
=======
	if timeSlot.BookedPatients.Int32 >= timeSlot.MaxPatients.Int32 {
		return nil, fmt.Errorf("time slot is fully booked")
	}
>>>>>>> dc47646 (Optimize SQL query)
	dateTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	var startTimeFormatted string
	var endTimeFormatted string

	// Create the appointment within a transaction
	var appointment db.Appointment
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
<<<<<<< HEAD
		// Create the appointment
		appointment, err = q.CreateAppointment(ctx, arg)
>>>>>>> 685da65 (latest update)
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

<<<<<<< HEAD
	appointment, err := s.storeDB.CreateAppointment(ctx, arg)
>>>>>>> cfbe865 (updated service response)
=======
		// Update the time slot status to "booked"
		err = q.UpdateTimeSlotStatus(ctx, db.UpdateTimeSlotStatusParams{
			ID:     req.TimeSlotID,
			Status: pgtype.Text{String: "pending", Valid: true},
=======

<<<<<<< HEAD
		// Lấy thông tin khung giờ và khóa bản ghi
<<<<<<< HEAD
=======

<<<<<<< HEAD
		// Lấy thông tin khung giờ và khóa bản ghi
<<<<<<< HEAD
>>>>>>> b393bb9 (add service and add permission)
		timeSlot, err := q.GetTimeSlot(ctx, db.GetTimeSlotParams{
			ID:       req.TimeSlotID,
			Date:     pgtype.Date{Time: dateTime, Valid: true},
			DoctorID: int32(req.DoctorID),
<<<<<<< HEAD
>>>>>>> b393bb9 (add service and add permission)
		})
=======
		timeSlot, err := q.GetTimeSlotForUpdate(ctx, req.TimeSlotID)
>>>>>>> ada3717 (Docker file)
		if err != nil {
			return fmt.Errorf("failed to get time slot: %w", err)
		}
		// Kiểm tra xem khung giờ còn chỗ trống không
		if timeSlot.BookedPatients.Int32 >= timeSlot.MaxPatients.Int32 {
			return fmt.Errorf("time slot is fully booked")
		}

=======
>>>>>>> dc47646 (Optimize SQL query)
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
>>>>>>> 685da65 (latest update)
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
<<<<<<< HEAD
=======
		return nil, fmt.Errorf("error while getting service: %w", err)
=======
=======
		// Update the time slot status to "booked"
		err = q.UpdateTimeSlotStatus(ctx, db.UpdateTimeSlotStatusParams{
			ID:     req.TimeSlotID,
			Status: pgtype.Text{String: "pending", Valid: true},
=======
>>>>>>> b393bb9 (add service and add permission)
		})
=======
		timeSlot, err := q.GetTimeSlotForUpdate(ctx, req.TimeSlotID)
>>>>>>> ada3717 (Docker file)
		if err != nil {
			return fmt.Errorf("failed to get time slot: %w", err)
		}
		// Kiểm tra xem khung giờ còn chỗ trống không
		if timeSlot.BookedPatients.Int32 >= timeSlot.MaxPatients.Int32 {
			return fmt.Errorf("time slot is fully booked")
		}

=======
>>>>>>> dc47646 (Optimize SQL query)
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

<<<<<<< HEAD
>>>>>>> 685da65 (latest update)
	// Fetch related data
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service: %w", err)
<<<<<<< HEAD
>>>>>>> 685da65 (latest update)
=======
>>>>>>> 685da65 (latest update)
	}

	pet, err := s.storeDB.GetPetByID(ctx, appointment.Petid.Int64)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pet: %w", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	state, err := s.storeDB.GetState(ctx, int64(appointment.StateID.Int32))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch state: %w", err)
=======
	detail, err := s.storeDB.GetAppointmentDetail(ctx, db.GetAppointmentDetailParams{
		ID:    appointment.ServiceID.Int64,
		Petid: appointment.Petid.Int64,
		ID_2:  int64(appointment.StateID.Int32),
	})
=======
	detail, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, appointment.AppointmentID)

>>>>>>> 4ccd381 (Update appointment flow)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment detail: %w", err)
>>>>>>> dc47646 (Optimize SQL query)
=======
	state, err := s.storeDB.GetState(ctx, int64(appointment.StateID.Int32))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch state: %w", err)
>>>>>>> e859654 (Elastic search)
=======
	detail, err := s.storeDB.GetAppointmentDetail(ctx, db.GetAppointmentDetailParams{
		ID:    appointment.ServiceID.Int64,
		Petid: appointment.Petid.Int64,
		ID_2:  int64(appointment.StateID.Int32),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment detail: %w", err)
>>>>>>> dc47646 (Optimize SQL query)
	}

	// Prepare the response
=======
>>>>>>> cfbe865 (updated service response)
=======
	startTime := time.UnixMicro(timeSlot.StartTime.Microseconds).UTC()
	startTimeFormatted := startTime.Format("15:04:05")

	// Format the end time
	endTime := time.UnixMicro(timeSlot.EndTime.Microseconds).UTC()
	endTimeFormatted := endTime.Format("15:04:05")

=======
>>>>>>> b393bb9 (add service and add permission)
	// Prepare the response
>>>>>>> 685da65 (latest update)
	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
<<<<<<< HEAD
<<<<<<< HEAD
		PetName:     detail.PetName.String,
		Reason:      detail.AppointmentReason.String,
		Date:        appointment.Date.Time.Format(time.RFC3339),
<<<<<<< HEAD
<<<<<<< HEAD
=======
		PetName:     pet.Name,
		Date:        appointment.Date.Time.Format(time.RFC3339),
<<<<<<< HEAD
>>>>>>> 403647b (time of appointment)
		ServiceName: service.Name,
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> cfbe865 (updated service response)
		Note:        req.Note,
>>>>>>> cfbe865 (updated service response)
=======
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
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
=======
=======
		ServiceName: service.Name.String,
>>>>>>> b393bb9 (add service and add permission)
=======
		ServiceName: detail.ServiceName.String,
>>>>>>> dc47646 (Optimize SQL query)
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		State:        detail.StateName.String,
		ReminderSend: appointment.ReminderSend.Bool,
		CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
>>>>>>> 685da65 (latest update)
}
<<<<<<< HEAD

<<<<<<< HEAD
func (s *AppointmentService) CheckInAppoinment(ctx *gin.Context, id, roomID int64, priority string) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Fetch appointment details
		appointment, err := q.GetAppointmentDetailByAppointmentID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}

<<<<<<< HEAD
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
=======
=======
		ServiceName: service.Name.String,
>>>>>>> b393bb9 (add service and add permission)
=======
		PetName:     detail.PetName,
		Date:        appointment.Date.Time.Format(time.RFC3339),
		ServiceName: detail.ServiceName.String,
>>>>>>> dc47646 (Optimize SQL query)
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		State:        detail.StateName,
		Notes:        appointment.Notes.String,
		ReminderSend: appointment.ReminderSend.Bool,
		CreatedAt:    appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
>>>>>>> 685da65 (latest update)
}

<<<<<<< HEAD
func (s *AppointmentService) GetAppointmentByID(ctx *gin.Context, id int64) (*Appointment, error) {
	var err error

<<<<<<< HEAD
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
=======
	// Fetch the appointment details
	appointment, err := s.storeDB.GetAppointmentDetailById(ctx, id)
	if err != nil {
		return fmt.Errorf("error while fetching appointment: %w", err)
=======
	// Fetch the appointment details
	appointment, err := s.storeDB.GetAppointmentDetailById(ctx, id)
	if err != nil {
		return fmt.Errorf("error while fetching appointment: %w", err)
	}

	// Check if the appointment is already approved or rejected
	if appointment.Status.String == "approved" {
		return fmt.Errorf("appointment is already approved")
	}
	if appointment.Status.String == "rejected" {
		return fmt.Errorf("appointment is already rejected")
	}

	// Perform the approval within a transaction
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Update the appointment status to "approved"
		err := q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			AppointmentID: id,
			Status:        pgtype.Text{String: req.Status, Valid: true},
		})
=======
func (s *AppointmentService) ConfirmPayment(ctx context.Context, appointmentID int64) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
<<<<<<< HEAD
		// Lấy thông tin cuộc hẹn
		appointment, err := q.GetAppointmentDetailById(ctx, appointmentID)
>>>>>>> b393bb9 (add service and add permission)
=======
		appointment, err := q.GetAppointmentDetailByAppointmentID(ctx, appointmentID)
>>>>>>> dc47646 (Optimize SQL query)
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
>>>>>>> 685da65 (latest update)
	}

<<<<<<< HEAD
	// Check if the appointment is already approved or rejected
	if appointment.Status.String == "approved" {
		return fmt.Errorf("appointment is already approved")
	}
	if appointment.Status.String == "rejected" {
		return fmt.Errorf("appointment is already rejected")
	}

	// Perform the approval within a transaction
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Update the appointment status to "approved"
		err := q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			AppointmentID: id,
			Status:        pgtype.Text{String: req.Status, Valid: true},
		})
=======
func (s *AppointmentService) ConfirmPayment(ctx context.Context, appointmentID int64) error {
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
<<<<<<< HEAD
		// Lấy thông tin cuộc hẹn
		appointment, err := q.GetAppointmentDetailById(ctx, appointmentID)
>>>>>>> b393bb9 (add service and add permission)
=======
=======
func (s *AppointmentService) ConfirmPayment(ctx context.Context, appointmentID int64) error {
	return s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Fetch appointment details
>>>>>>> 4ccd381 (Update appointment flow)
		appointment, err := q.GetAppointmentDetailByAppointmentID(ctx, appointmentID)
>>>>>>> dc47646 (Optimize SQL query)
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
<<<<<<< HEAD
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
>>>>>>> 685da65 (latest update)
	}

<<<<<<< HEAD
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

<<<<<<< HEAD
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

=======
>>>>>>> e859654 (Elastic search)
	return nil
=======
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
>>>>>>> 4ccd381 (Update appointment flow)
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
=======
func (s *AppointmentService) GetAppointmentsOfDoctorService(ctx *gin.Context, doctorID int64) ([]AppointmentWithDetails, error) {

	// // Get all appointments with related data in a single query
	// appointments, err := s.storeDB.GetAppointmentsOfDoctorWithDetails(ctx, doctorID)
	// if err != nil {
	// 	return nil, fmt.Errorf("fetching appointments: %w", err)
	// }

	// // Pre-allocate slice with known capacity
	// result := make([]AppointmentWithDetails, 0, len(appointments))

	// for _, appt := range appointments {
	// 	// ts := timeslot{
	// 	// 	StartTime: appt.StartTime.Time.Format("2006-01-02 15:04:05"),
	// 	// 	EndTime:   appt.EndTime.Time.Format("2006-01-02 15:04:05"),
	// 	// }
	// 	service, err := s.storeDB.GetServiceByID(ctx, appt.ServiceID.Int64)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	pet, err := s.storeDB.GetPetByID(ctx, appt.Pet)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	doc, err := s.storeDB.GetDoctor(ctx, appt.DoctorID.Int64)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	result = append(result, AppointmentWithDetails{
	// 		AppointmentID: appt.AppointmentID,
	// 		PetName:       appt.PetName.String,
	// 		ServiceName:   appt.ServiceName.String,
	// 		DoctorName:    appt.DoctorName,
	// 		Date:          appt.Date.Time.Format("2006-01-02 15:04:05"),
	// 		Status:        appt.Status.String,
	// 		Notes:         appt.Notes.String,
	// 		ReminderSend:  appt.ReminderSend.Bool,
	// 	})
	// }

	return nil, nil

>>>>>>> 98e9e45 (ratelimit and recovery function)
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
=======
func (s *AppointmentService) GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error) {
<<<<<<< HEAD
	// Fetch appointment details
>>>>>>> 685da65 (latest update)
	appointment, err := s.storeDB.GetAppointmentDetailById(ctx, id)
=======
=======
func (s *AppointmentService) GetAppointmentByID(ctx *gin.Context, id int64) (*Appointment, error) {
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
	var err error

	appointment, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, id)
>>>>>>> dc47646 (Optimize SQL query)
	if err != nil {
		return nil, fmt.Errorf("Cannot get appointment detail")
	}
<<<<<<< HEAD
	return &appointment, nil
>>>>>>> 7e35c2e (get appointment detail)
}
=======
>>>>>>> 685da65 (latest update)

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

=======
>>>>>>> e859654 (Elastic search)
	return nil
}

<<<<<<< HEAD
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
=======
func (s *AppointmentService) GetAppointmentsOfDoctorService(ctx *gin.Context, doctorID int64) ([]AppointmentWithDetails, error) {

	// // Get all appointments with related data in a single query
	// appointments, err := s.storeDB.GetAppointmentsOfDoctorWithDetails(ctx, doctorID)
	// if err != nil {
	// 	return nil, fmt.Errorf("fetching appointments: %w", err)
	// }

	// // Pre-allocate slice with known capacity
	// result := make([]AppointmentWithDetails, 0, len(appointments))

	// for _, appt := range appointments {
	// 	// ts := timeslot{
	// 	// 	StartTime: appt.StartTime.Time.Format("2006-01-02 15:04:05"),
	// 	// 	EndTime:   appt.EndTime.Time.Format("2006-01-02 15:04:05"),
	// 	// }
	// 	service, err := s.storeDB.GetServiceByID(ctx, appt.ServiceID.Int64)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	pet, err := s.storeDB.GetPetByID(ctx, appt.Pet)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	doc, err := s.storeDB.GetDoctor(ctx, appt.DoctorID.Int64)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	result = append(result, AppointmentWithDetails{
	// 		AppointmentID: appt.AppointmentID,
	// 		PetName:       appt.PetName.String,
	// 		ServiceName:   appt.ServiceName.String,
	// 		DoctorName:    appt.DoctorName,
	// 		Date:          appt.Date.Time.Format("2006-01-02 15:04:05"),
	// 		Status:        appt.Status.String,
	// 		Notes:         appt.Notes.String,
	// 		ReminderSend:  appt.ReminderSend.Bool,
	// 	})
	// }

	return nil, nil

>>>>>>> 98e9e45 (ratelimit and recovery function)
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
=======
func (s *AppointmentService) GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error) {
<<<<<<< HEAD
	// Fetch appointment details
>>>>>>> 685da65 (latest update)
	appointment, err := s.storeDB.GetAppointmentDetailById(ctx, id)
=======
	var err error

	appointment, err := s.storeDB.GetAppointmentDetailByAppointmentID(ctx, id)
>>>>>>> dc47646 (Optimize SQL query)
	if err != nil {
		return nil, fmt.Errorf("Cannot get appointment detail")
	}
<<<<<<< HEAD
	return &appointment, nil
>>>>>>> 7e35c2e (get appointment detail)
}
=======
>>>>>>> 685da65 (latest update)

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
