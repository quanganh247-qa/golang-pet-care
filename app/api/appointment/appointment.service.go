package appointment

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type AppointmentServiceInterface interface {
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error)
	UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error
	GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error)
	GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error)
}

func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error) {
	// Validate input
	if req.DoctorID == 0 || req.PetID == 0 || req.ServiceID == 0 || req.TimeSlotID == 0 || req.Date == "" {
		return nil, fmt.Errorf("missing required fields: doctor_id, pet_id, service_id, time_slot_id, or date")
	}

	// Fetch doctor details
	doctor, err := s.storeDB.GetDoctor(ctx, req.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch doctor: %w", err)
	}

	// Parse and validate the date
	dateTime, err := time.Parse("2006-01-02T15:04:05Z", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// Fetch the time slot
	timeSlot, err := s.storeDB.GetTimeSlotById(ctx, req.TimeSlotID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch time slot: %w", err)
	}

	// Prepare appointment parameters
	arg := db.CreateAppointmentParams{
		DoctorID:     pgtype.Int8{Int64: doctor.ID, Valid: true},
		Petid:        pgtype.Int8{Int64: req.PetID, Valid: true},
		ServiceID:    pgtype.Int8{Int64: req.ServiceID, Valid: true},
		Date:         pgtype.Timestamp{Time: dateTime, Valid: true},
		TimeSlotID:   pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
		Notes:        pgtype.Text{String: req.Note, Valid: true},
		Status:       pgtype.Text{String: "pending", Valid: true},
		ReminderSend: pgtype.Bool{Bool: false, Valid: true},
	}

	// Create the appointment within a transaction
	var appointment db.Appointment
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Create the appointment
		appointment, err = q.CreateAppointment(ctx, arg)
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

		// Update the time slot status to "booked"
		err = q.UpdateTimeSlotStatus(ctx, db.UpdateTimeSlotStatusParams{
			ID:     req.TimeSlotID,
			Status: pgtype.Text{String: "pending", Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to update time slot status: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	// Fetch related data
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service: %w", err)
	}

	pet, err := s.storeDB.GetPetByID(ctx, appointment.Petid.Int64)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pet: %w", err)
	}

	startTime := time.UnixMicro(timeSlot.StartTime.Microseconds).UTC()
	startTimeFormatted := startTime.Format("15:04:05")

	// Format the end time
	endTime := time.UnixMicro(timeSlot.EndTime.Microseconds).UTC()
	endTimeFormatted := endTime.Format("15:04:05")

	// Prepare the response
	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     pet.Name,
		Date:        appointment.Date.Time.Format(time.RFC3339),
		ServiceName: service.Name,
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		Status:       appointment.Status.String,
		Notes:        appointment.Notes.String,
		ReminderSend: appointment.ReminderSend.Bool,
		CreatedAt:    appointment.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}

func (s *AppointmentService) UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error {

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
		if err != nil {
			return fmt.Errorf("error while updating appointment status: %w", err)
		}

		// Notify the user (optional)
		// You can implement a notification system here (e.g., email, SMS, in-app notification)
		log.Printf("Appointment %d has been approved. Notifying user...", id)

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

func (s *AppointmentService) GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error) {
	// Fetch appointment details
	appointment, err := s.storeDB.GetAppointmentDetailById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting appointment by id: %w", err)
	}

	// Fetch doctor details
	doctor, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
	if err != nil {
		return nil, fmt.Errorf("error while fetching doctor: %w", err)
	}

	// Fetch pet details
	pet, err := s.storeDB.GetPetByID(ctx, appointment.Petid.Int64)
	if err != nil {
		return nil, fmt.Errorf("error while fetching pet: %w", err)
	}

	// Fetch service details
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
	if err != nil {
		return nil, fmt.Errorf("error while fetching service: %w", err)
	}

	// Fetch time slot details
	timeSlot, err := s.storeDB.GetTimeSlotById(ctx, appointment.TimeSlotID.Int64)
	if err != nil {
		return nil, fmt.Errorf("error while fetching time slot: %w", err)
	}

	// Format start and end times
	startTime := time.UnixMicro(timeSlot.StartTime.Microseconds).UTC()
	startTimeFormatted := startTime.Format("15:04:05")

	endTime := time.UnixMicro(timeSlot.EndTime.Microseconds).UTC()
	endTimeFormatted := endTime.Format("15:04:05")

	// Prepare the response
	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     pet.Name,
		Date:        appointment.Date.Time.Format(time.RFC3339),
		ServiceName: service.Name,
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		Status:       appointment.Status.String,
		Notes:        appointment.Notes.String,
		ReminderSend: appointment.ReminderSend.Bool,
		CreatedAt:    appointment.CreatedAt.Time.Format(time.RFC3339),
	}, nil
}
func (s *AppointmentService) GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error) {
	rows, err := s.storeDB.GetAppointmentsByUser(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		return nil, err
	}

	var a []createAppointmentResponse
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, row := range rows {
		wg.Add(1)
		go func(row db.GetAppointmentsByUserRow) {
			defer wg.Done()

			service, err := s.storeDB.GetServiceByID(ctx, row.ServiceID.Int64)
			if err != nil {
				log.Printf("Failed to get service for appointment %d: %v", row.AppointmentID, err)
				return
			}

			pet, err := s.storeDB.GetPetByID(ctx, row.Petid)
			if err != nil {
				log.Printf("Failed to get pet for appointment %d: %v", row.AppointmentID, err)
				return
			}

			doc, err := s.storeDB.GetDoctor(ctx, row.DoctorID.Int64)
			if err != nil {
				log.Printf("Failed to get doctor for appointment %d: %v", row.AppointmentID, err)
				return
			}

			// Fetch time slot details
			timeSlot, err := s.storeDB.GetTimeSlotById(ctx, row.TimeSlotID.Int64)
			if err != nil {
				log.Printf("Failed to get time slot for appointment %d: %v", row.AppointmentID, err)
				return
			}

			// Format start and end times
			startTime := time.UnixMicro(timeSlot.StartTime.Microseconds).UTC()
			startTimeFormatted := startTime.Format("15:04:05")

			endTime := time.UnixMicro(timeSlot.EndTime.Microseconds).UTC()
			endTimeFormatted := endTime.Format("15:04:05")

			mu.Lock()
			a = append(a, createAppointmentResponse{
				ID:          row.AppointmentID,
				PetName:     pet.Name,
				DoctorName:  doc.Name,
				ServiceName: service.Name,
				TimeSlot: timeslot{
					StartTime: startTimeFormatted,
					EndTime:   endTimeFormatted,
				},
				Date:      row.Date.Time.Format(time.RFC3339),
				Status:    row.Status.String,
				CreatedAt: row.CreatedAt.Time.Format(time.RFC3339),
			})
			mu.Unlock()
		}(row)
	}

	wg.Wait()
	return a, nil
}

// cancle appointment
func (s *AppointmentService) CancelAppointment(ctx *gin.Context, id int64) error {
	// Fetch the appointment details
	appointment, err := s.storeDB.GetAppointmentDetailById(ctx, id)
	if err != nil {
		return fmt.Errorf("error while fetching appointment: %w", err)
	}

	// Check if the appointment is already cancelled
	if appointment.Status.String == "cancelled" {
		return fmt.Errorf("appointment is already cancelled")
	}

	// Perform the cancellation within a transaction
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Update the appointment status to "cancelled"
		err := q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			AppointmentID: id,
			Status:        pgtype.Text{String: "cancelled", Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error while updating appointment status: %w", err)
		}

		// Free the associated time slot by updating its status to "available"
		err = q.UpdateTimeSlotStatus(ctx, db.UpdateTimeSlotStatusParams{
			ID:     appointment.TimeSlotID.Int64,
			Status: pgtype.Text{String: "available", Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error while updating time slot status: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}
