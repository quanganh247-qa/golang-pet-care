package appointment

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type AppointmentServiceInterface interface {
}

// creating an appointment by time slot available of doctor
func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error) {

	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.timeSlotID)
	if err != nil {
		return nil, fmt.Errorf("error while getting time slot: %w", err)
	}

	doctor, err := s.storeDB.GetDoctor(ctx, timeSlot.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("error while getting doctor: %w", err)
	}

	appointment, err := s.storeDB.CreateAppointment(ctx, db.CreateAppointmentParams{
		DoctorID:   pgtype.Int8{Int64: doctor.ID, Valid: true},
		Petid:      pgtype.Int8{Int64: req.petID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.serviceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.timeSlotID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("error while creating appointment: %w", err)
	}

	service, err := s.storeDB.GetService(ctx, appointment.ServiceID.Int64)
	if err != nil {
		return nil, fmt.Errorf("error while getting service: %w", err)
	}

	pet, err := s.storeDB.GetPetByID(ctx, appointment.Petid.Int64)
	if err != nil {
		return nil, fmt.Errorf("error while getting pet: %w", err)
	}

	resTimeSlot := db.Timeslot{
		DoctorID:  timeSlot.DoctorID,
		StartTime: timeSlot.StartTime,
		EndTime:   timeSlot.EndTime,
	}

	return &createAppointmentResponse{
		id:          appointment.AppointmentID,
		doctorName:  doctor.Name,
		petName:     pet.Name,
		serviceName: service.Name,
		timeSlot:    resTimeSlot,
		note:        req.note,
	}, nil

}
