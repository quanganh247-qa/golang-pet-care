package appointment

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type AppointmentServiceInterface interface {
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error)
<<<<<<< HEAD
<<<<<<< HEAD
	UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error
	GetAppointmentsOfDoctorService(ctx *gin.Context, doctorID int64) ([]AppointmentWithDetails, error)
<<<<<<< HEAD
=======
>>>>>>> 323513c (appointment api)
=======
	UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error
>>>>>>> 7cfffa9 (update dtb and appointment)
=======
>>>>>>> 4b8e9b6 (update appointment api)
}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
// creating an appointment by time slot available of doctor
func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error) {
=======
// func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (createAppointmentResponse, error) {
>>>>>>> c7f463c (update dtb)
=======
// func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (createAppointmentResponse, error) {
>>>>>>> c7f463c (update dtb)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	println("CreateAppointment", req.TimeSlotID)
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.TimeSlotID)
=======
<<<<<<< HEAD
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.timeSlotID)
>>>>>>> a5cefab (modify type of filed in dtb)
=======
	println("CreateAppointment", req.TimeSlotID)
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.TimeSlotID)
>>>>>>> 9c7c7b8 (update dtb and appointment)
=======
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.timeSlotID)
>>>>>>> adc2e22 (modify type of filed in dtb)
=======
	println("CreateAppointment", req.TimeSlotID)
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.TimeSlotID)
>>>>>>> f96fe68 (update dtb and appointment)
=======
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.timeSlotID)
>>>>>>> 7833094 (modify type of filed in dtb)
=======
	println("CreateAppointment", req.TimeSlotID)
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.TimeSlotID)
>>>>>>> 430a2a2 (update dtb and appointment)
=======
// creating an appointment by time slot available of doctor
func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest) (*createAppointmentResponse, error) {

<<<<<<< HEAD
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.timeSlotID)
>>>>>>> 59d4ef2 (modify type of filed in dtb)
<<<<<<< HEAD
>>>>>>> a1c3177 (modify type of filed in dtb)
=======
=======
	println("CreateAppointment", req.TimeSlotID)
	timeSlot, err := s.storeDB.GetTimeSlotByID(ctx, req.TimeSlotID)
>>>>>>> 7cfffa9 (update dtb and appointment)
>>>>>>> 940e5bf (update dtb and appointment)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> a1c3177 (modify type of filed in dtb)
=======
>>>>>>> 940e5bf (update dtb and appointment)
		Petid:      pgtype.Int8{Int64: req.PetID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.ServiceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
=======
		Petid:      pgtype.Int8{Int64: req.petID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.serviceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.timeSlotID, Valid: true},
<<<<<<< HEAD
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
=======
		Petid:      pgtype.Int8{Int64: req.PetID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.ServiceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
>>>>>>> f96fe68 (update dtb and appointment)
=======
		Petid:      pgtype.Int8{Int64: req.petID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.serviceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.timeSlotID, Valid: true},
>>>>>>> 7833094 (modify type of filed in dtb)
=======
		Petid:      pgtype.Int8{Int64: req.PetID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.ServiceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
>>>>>>> 430a2a2 (update dtb and appointment)
=======
>>>>>>> 59d4ef2 (modify type of filed in dtb)
<<<<<<< HEAD
>>>>>>> a1c3177 (modify type of filed in dtb)
=======
=======
		Petid:      pgtype.Int8{Int64: req.PetID, Valid: true},
		ServiceID:  pgtype.Int8{Int64: req.ServiceID, Valid: true},
		TimeSlotID: pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
>>>>>>> 7cfffa9 (update dtb and appointment)
>>>>>>> 940e5bf (update dtb and appointment)
	})
	if err != nil {
		return nil, fmt.Errorf("error while creating appointment: %w", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
=======
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
>>>>>>> c4ee544 (update service api)
=======
	service, err := s.storeDB.GetService(ctx, appointment.ServiceID.Int64)
>>>>>>> 7833094 (modify type of filed in dtb)
=======
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
>>>>>>> c9d6049 (update service api)
=======
=======
>>>>>>> 860385e (update service api)
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
=======
	service, err := s.storeDB.GetService(ctx, appointment.ServiceID.Int64)
>>>>>>> 59d4ef2 (modify type of filed in dtb)
<<<<<<< HEAD
>>>>>>> a1c3177 (modify type of filed in dtb)
=======
=======
	service, err := s.storeDB.GetServiceByID(ctx, appointment.ServiceID.Int64)
>>>>>>> 6e40c8e (update service api)
>>>>>>> 860385e (update service api)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> a1c3177 (modify type of filed in dtb)
=======
>>>>>>> 940e5bf (update dtb and appointment)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> a26c42a (update appointment api)

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
=======
>>>>>>> 836bb95 (update dtb)
// }
=======
<<<<<<< HEAD
=======
>>>>>>> adc2e22 (modify type of filed in dtb)
=======
>>>>>>> 7833094 (modify type of filed in dtb)
=======
>>>>>>> a1c3177 (modify type of filed in dtb)
	resTimeSlot := db.Timeslot{
		DoctorID:  timeSlot.DoctorID,
		StartTime: timeSlot.StartTime,
		EndTime:   timeSlot.EndTime,
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 940e5bf (update dtb and appointment)
=======
	resTimeSlot := timeslot{
		StartTime: timeSlot.StartTime.Time.Format("2006-01-02 15:04:05"),
		EndTime:   timeSlot.EndTime.Time.Format("2006-01-02 15:04:05"),
<<<<<<< HEAD
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
=======
	resTimeSlot := timeslot{
		StartTime: timeSlot.StartTime.Time.Format("2006-01-02 15:04:05"),
		EndTime:   timeSlot.EndTime.Time.Format("2006-01-02 15:04:05"),
>>>>>>> f96fe68 (update dtb and appointment)
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
>>>>>>> adc2e22 (modify type of filed in dtb)
=======

func (s *AppointmentService) UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error {
	fmt.Println("UpdateAppointmentStatus", req.Status, id)
=======
>>>>>>> aa5c8ab (update appointment api)

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
>>>>>>> f96fe68 (update dtb and appointment)
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
>>>>>>> aa5c8ab (update appointment api)
=======
=======
	resTimeSlot := timeslot{
		StartTime: timeSlot.StartTime.Time.Format("2006-01-02 15:04:05"),
		EndTime:   timeSlot.EndTime.Time.Format("2006-01-02 15:04:05"),
>>>>>>> 430a2a2 (update dtb and appointment)
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
>>>>>>> 7833094 (modify type of filed in dtb)
=======

func (s *AppointmentService) UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error {
	fmt.Println("UpdateAppointmentStatus", req.Status, id)
=======
>>>>>>> 7697b39 (update appointment api)

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
>>>>>>> 430a2a2 (update dtb and appointment)
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
>>>>>>> 7697b39 (update appointment api)
=======
=======
>>>>>>> 7cfffa9 (update dtb and appointment)
>>>>>>> 940e5bf (update dtb and appointment)
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
>>>>>>> 59d4ef2 (modify type of filed in dtb)
<<<<<<< HEAD
>>>>>>> a1c3177 (modify type of filed in dtb)
=======
=======

func (s *AppointmentService) UpdateAppointmentStatus(ctx *gin.Context, req updateAppointmentStatusRequest, id int64) error {
	fmt.Println("UpdateAppointmentStatus", req.Status, id)
=======
>>>>>>> 4b8e9b6 (update appointment api)

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
>>>>>>> 7cfffa9 (update dtb and appointment)
<<<<<<< HEAD
>>>>>>> 940e5bf (update dtb and appointment)
=======
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
>>>>>>> 4b8e9b6 (update appointment api)
>>>>>>> a26c42a (update appointment api)
