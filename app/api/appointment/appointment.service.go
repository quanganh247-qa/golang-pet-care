package appointment

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type AppointmentServiceInterface interface {
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error)
	UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error
	GetAppointmentsOfDoctorService(ctx *gin.Context, doctorID int64) ([]AppointmentWithDetails, error)
}

// creating an appointment by time slot available of doctor
func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error) {

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	println("CreateAppointment", req.TimeSlotID)
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.TimeSlotID)
=======
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.timeSlotID)
>>>>>>> a5cefab (modify type of filed in dtb)
=======
	println("CreateAppointment", req.TimeSlotID)
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.TimeSlotID)
>>>>>>> 9c7c7b8 (update dtb and appointment)
=======
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.timeSlotID)
>>>>>>> adc2e22 (modify type of filed in dtb)
	if err != nil {
		return nil, fmt.Errorf("error while getting time slot: %w", err)
	}

	doctor, err := s.storeDB.GetDoctor(ctx, timeSlot.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("error while getting doctor: %w", err)
	}

	appointment, err := s.storeDB.CreateAppointment(ctx, db.CreateAppointmentParams{
		DoctorID:   pgtype.Int8{Int64: doctor.ID, Valid: true},
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		Petid:      pgtype.Int8{Int64: req.PetID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.ServiceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
=======
		Petid:      pgtype.Int8{Int64: req.petID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.serviceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.timeSlotID, Valid: true},
>>>>>>> a5cefab (modify type of filed in dtb)
=======
		Petid:      pgtype.Int8{Int64: req.PetID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.ServiceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
>>>>>>> 9c7c7b8 (update dtb and appointment)
=======
		Petid:      pgtype.Int8{Int64: req.petID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.serviceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.timeSlotID, Valid: true},
>>>>>>> adc2e22 (modify type of filed in dtb)
	})
	if err != nil {
		return nil, fmt.Errorf("error while creating appointment: %w", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
=======
	service, err := s.storeDB.GetService(ctx, appointment.ServiceID.Int64)
>>>>>>> a5cefab (modify type of filed in dtb)
=======
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
>>>>>>> 1c933f4 (update service api)
=======
	service, err := s.storeDB.GetService(ctx, appointment.ServiceID.Int64)
>>>>>>> adc2e22 (modify type of filed in dtb)
	if err != nil {
		return nil, fmt.Errorf("error while getting service: %w", err)
	}

	pet, err := s.storeDB.GetPetByID(ctx, appointment.Petid.Int64)
	if err != nil {
		return nil, fmt.Errorf("error while getting pet: %w", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	resTimeSlot := timeslot{
		StartTime: timeSlot.StartTime.Time.Format("2006-01-02 15:04:05"),
		EndTime:   timeSlot.EndTime.Time.Format("2006-01-02 15:04:05"),
	}

	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     pet.Name,
		ServiceName: service.Name,
		TimeSlot:    resTimeSlot,
		Note:        req.Note,
	}, nil

}

func (s *AppointmentService) UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error {
<<<<<<< HEAD

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		return q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			Status:        pgtype.Text{String: req.Status, Valid: true},
			AppointmentID: id,
		})
	})
	if err != nil {
		return fmt.Errorf("error while updating appointment status: %w", err)
	}
	return nil
}

func (s *AppointmentService) GetAppointmentsOfDoctorService(ctx *gin.Context, doctorID int64) ([]AppointmentWithDetails, error) {

	// Get all appointments with related data in a single query
	appointments, err := s.storeDB.GetAppointmentsOfDoctorWithDetails(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("fetching appointments: %w", err)
	}

	// Pre-allocate slice with known capacity
	result := make([]AppointmentWithDetails, 0, len(appointments))

	for _, appt := range appointments {
		ts := timeslot{
			StartTime: appt.StartTime.Time.Format("2006-01-02 15:04:05"),
			EndTime:   appt.EndTime.Time.Format("2006-01-02 15:04:05"),
		}

		result = append(result, AppointmentWithDetails{
			AppointmentID: appt.AppointmentID,
			PetName:       appt.PetName.String,
			ServiceName:   appt.ServiceName.String,
			StartTime:     ts.StartTime,
			EndTime:       ts.EndTime,
		})
	}

	return result, nil

}

// func (s *AppointmentService) UpdateAppointmentNotificationService(ctx *gin.Context, id int64) error {
// 	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		return q.UpdateAppointmentNotification(ctx, id)
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error while updating appointment notification: %w", err)
// 	}
// 	return nil
=======
>>>>>>> 9472ea3 (update dtb)
=======
>>>>>>> bb4b692 (update dtb)
// }
=======
=======
>>>>>>> adc2e22 (modify type of filed in dtb)
	resTimeSlot := db.Timeslot{
		DoctorID:  timeSlot.DoctorID,
		StartTime: timeSlot.StartTime,
		EndTime:   timeSlot.EndTime,
<<<<<<< HEAD
=======
	resTimeSlot := timeslot{
		StartTime: timeSlot.StartTime.Time.Format("2006-01-02 15:04:05"),
		EndTime:   timeSlot.EndTime.Time.Format("2006-01-02 15:04:05"),
>>>>>>> 9c7c7b8 (update dtb and appointment)
	}

	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     pet.Name,
		ServiceName: service.Name,
		TimeSlot:    resTimeSlot,
		Note:        req.Note,
	}, nil

}
<<<<<<< HEAD
>>>>>>> a5cefab (modify type of filed in dtb)
=======

func (s *AppointmentService) UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error {
	fmt.Println("UpdateAppointmentStatus", req.Status, id)
=======
>>>>>>> 18fce3d (update appointment api)

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		return q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			Status:        pgtype.Text{String: req.Status, Valid: true},
			AppointmentID: id,
		})
	})
	if err != nil {
		return fmt.Errorf("error while updating appointment status: %w", err)
	}
	return nil
}
<<<<<<< HEAD
>>>>>>> 9c7c7b8 (update dtb and appointment)
=======

func (s *AppointmentService) GetAppointmentsOfDoctorService(ctx *gin.Context, doctorID int64) ([]AppointmentWithDetails, error) {

	// Get all appointments with related data in a single query
	appointments, err := s.storeDB.GetAppointmentsOfDoctorWithDetails(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("fetching appointments: %w", err)
	}

	// Pre-allocate slice with known capacity
	result := make([]AppointmentWithDetails, 0, len(appointments))

	for _, appt := range appointments {
		ts := timeslot{
			StartTime: appt.StartTime.Time.Format("2006-01-02 15:04:05"),
			EndTime:   appt.EndTime.Time.Format("2006-01-02 15:04:05"),
		}

		result = append(result, AppointmentWithDetails{
			AppointmentID: appt.AppointmentID,
			PetName:       appt.PetName.String,
			ServiceName:   appt.ServiceName.String,
			StartTime:     ts.StartTime,
			EndTime:       ts.EndTime,
		})
	}

	return result, nil

}

// func (s *AppointmentService) UpdateAppointmentNotificationService(ctx *gin.Context, id int64) error {
// 	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
// 		return q.UpdateAppointmentNotification(ctx, id)
// 	})
// 	if err != nil {
// 		return fmt.Errorf("error while updating appointment notification: %w", err)
// 	}
// 	return nil
// }
>>>>>>> 18fce3d (update appointment api)
=======
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
>>>>>>> adc2e22 (modify type of filed in dtb)
