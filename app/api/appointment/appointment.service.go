package appointment

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type AppointmentServiceInterface interface {
	CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error)
	ConfirmPayment(ctx context.Context, appointmentID int64) error
	GetAppointmentByID(ctx *gin.Context, id int64) (*createAppointmentResponse, error)
	GetAppointmentsByUser(ctx *gin.Context, username string) ([]createAppointmentResponse, error)
	GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error)
	GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error)
}

func (s *AppointmentService) CreateAppointment(ctx *gin.Context, req createAppointmentRequest, username string) (*createAppointmentResponse, error) {

	// Fetch doctor details
	doctor, err := s.storeDB.GetDoctor(ctx, req.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch doctor: %w", err)
	}

	// Parse and validate the date
	dateTime, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	var startTimeFormatted string
	var endTimeFormatted string

	// Create the appointment within a transaction
	var appointment db.Appointment
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {

		// Lấy thông tin khung giờ và khóa bản ghi
		timeSlot, err := q.GetTimeSlot(ctx, db.GetTimeSlotParams{
			ID:       req.TimeSlotID,
			Date:     pgtype.Date{Time: dateTime, Valid: true},
			DoctorID: int32(req.DoctorID),
		})
		if err != nil {
			return fmt.Errorf("failed to get time slot: %w", err)
		}

		startTime := time.UnixMicro(timeSlot.StartTime.Microseconds).UTC()
		startTimeFormatted = startTime.Format("15:04:05")

		// Format the end time
		endTime := time.UnixMicro(timeSlot.EndTime.Microseconds).UTC()
		endTimeFormatted = endTime.Format("15:04:05")

		// Kiểm tra xem khung giờ còn chỗ trống không
		if timeSlot.BookedPatients.Int32 >= timeSlot.MaxPatients.Int32 {
			return fmt.Errorf("time slot is fully booked")
		}

		// Tạo cuộc hẹn với trạng thái "pending"
		appointment, err = q.CreateAppointment(ctx, db.CreateAppointmentParams{
			DoctorID:      pgtype.Int8{Int64: doctor.ID, Valid: true},
			Petid:         pgtype.Int8{Int64: req.PetID, Valid: true},
			ServiceID:     pgtype.Int8{Int64: req.ServiceID, Valid: true},
			Date:          pgtype.Timestamp{Time: dateTime, Valid: true},
			TimeSlotID:    pgtype.Int8{Int64: req.TimeSlotID, Valid: true},
			Notes:         pgtype.Text{String: req.Note, Valid: true},
			PaymentStatus: pgtype.Text{String: "pending", Valid: true},
			ReminderSend:  pgtype.Bool{Bool: false, Valid: true},
			Username:      pgtype.Text{String: username, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
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

	// Prepare the response
	return &createAppointmentResponse{
		ID:          appointment.AppointmentID,
		DoctorName:  doctor.Name,
		PetName:     pet.Name,
		Date:        appointment.Date.Time.Format(time.RFC3339),
		ServiceName: service.Name.String,
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		PaymentStatus: appointment.PaymentStatus.String,
		Notes:         appointment.Notes.String,
		ReminderSend:  appointment.ReminderSend.Bool,
		CreatedAt:     appointment.CreatedAt.Time.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *AppointmentService) ConfirmPayment(ctx context.Context, appointmentID int64) error {
	// Bắt đầu transaction
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		// Lấy thông tin cuộc hẹn
		appointment, err := q.GetAppointmentDetailById(ctx, appointmentID)
		if err != nil {
			return fmt.Errorf("failed to get appointment: %w", err)
		}

		// Kiểm tra xem cuộc hẹn đã được thanh toán chưa
		if appointment.PaymentStatus.String == "paid" {
			return fmt.Errorf("appointment is already paid")
		}

		// Lấy thông tin khung giờ và khóa bản ghi
		timeSlot, err := q.GetTimeSlot(ctx, db.GetTimeSlotParams{
			ID:       appointment.TimeSlotID.Int64,
			Date:     pgtype.Date{Time: appointment.Date.Time, Valid: true},
			DoctorID: int32(appointment.DoctorID.Int64),
		})
		if err != nil {
			return fmt.Errorf("failed to get time slot: %w", err)
		}

		if timeSlot.BookedPatients.Int32 >= timeSlot.MaxPatients.Int32 {
			return fmt.Errorf("time slot is fully booked")
		}

		// Cập nhật trạng thái thanh toán và cuộc hẹn
		err = q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
			AppointmentID: appointmentID,
			PaymentStatus: pgtype.Text{String: "paid", Valid: true},
		})
		if err != nil {
			log.Printf("Failed to update appointment payment status: %v", err)
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		// Tăng số lượng bệnh nhân đã đặt lịch trong khung giờ
		err = q.UpdateTimeSlotBookedPatients(ctx, db.UpdateTimeSlotBookedPatientsParams{
			ID:       appointment.TimeSlotID.Int64,
			DoctorID: int32(appointment.DoctorID.Int64),
		})
		if err != nil {
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
		ServiceName: service.Name.String,
		TimeSlot: timeslot{
			StartTime: startTimeFormatted,
			EndTime:   endTimeFormatted,
		},
		Notes:         appointment.Notes.String,
		ReminderSend:  appointment.ReminderSend.Bool,
		PaymentStatus: appointment.PaymentStatus.String,
		CreatedAt:     appointment.CreatedAt.Time.Format(time.RFC3339),
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
				ServiceName: service.Name.String,
				TimeSlot: timeslot{
					StartTime: startTimeFormatted,
					EndTime:   endTimeFormatted,
				},
				PaymentStatus: row.PaymentStatus.String,
				Date:          row.Date.Time.Format(time.RFC3339),
				CreatedAt:     row.CreatedAt.Time.Format(time.RFC3339),
			})
			mu.Unlock()
		}(row)
	}

	wg.Wait()
	return a, nil
}

func (s *AppointmentService) GetAppointmentsByDoctor(ctx *gin.Context, doctorID int64) ([]createAppointmentResponse, error) {
	// Lấy danh sách lịch hẹn theo doctor_id
	appointments, err := s.storeDB.GetAppointmentsByDoctor(ctx, pgtype.Int8{Int64: doctorID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}

	// Format danh sách lịch hẹn
	var response []createAppointmentResponse
	for _, appointment := range appointments {

		doc, err := s.storeDB.GetDoctor(ctx, appointment.DoctorID.Int64)
		if err != nil {
			log.Printf("Failed to get doctor for appointment %d: %v", appointment.AppointmentID, err)
			return nil, fmt.Errorf("failed to get doctor: %w", err)
		}

		// Format start and end times
		startTime := time.UnixMicro(appointment.StartTime.Microseconds).UTC()
		startTimeFormatted := startTime.Format("15:04:05")

		endTime := time.UnixMicro(appointment.EndTime.Microseconds).UTC()
		endTimeFormatted := endTime.Format("15:04:05")

		response = append(response, createAppointmentResponse{
			ID:          appointment.AppointmentID,
			DoctorName:  doc.Name,
			PetName:     appointment.PetName,
			Date:        appointment.Date.Time.Format(time.RFC3339),
			ServiceName: appointment.ServiceName.String,
			TimeSlot: timeslot{
				StartTime: startTimeFormatted,
				EndTime:   endTimeFormatted,
			},
			Notes:         appointment.Notes.String,
			PaymentStatus: appointment.PaymentStatus.String,
			ReminderSend:  appointment.ReminderSend.Bool,
			CreatedAt:     appointment.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *AppointmentService) GetAvailableTimeSlots(ctx *gin.Context, doctorID int64, date string) ([]timeSlotResponse, error) {
	// Parse ngày
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// Lấy danh sách khung giờ của bác sĩ trong ngày cụ thể
	timeSlots, err := s.storeDB.GetTimeSlotsByDoctorAndDate(ctx, db.GetTimeSlotsByDoctorAndDateParams{
		DoctorID: int32(doctorID),
		Date:     pgtype.Date{Time: dateTime, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get time slots: %w", err)
	}

	// Lọc ra các khung giờ còn chỗ trống
	var availableTimeSlots []timeSlotResponse
	var slotRes timeSlotResponse
	for _, slot := range timeSlots {
		// Format start and end times
		startTime := time.UnixMicro(slot.StartTime.Microseconds).UTC()
		startTimeFormatted := startTime.Format("15:04:05")

		endTime := time.UnixMicro(slot.EndTime.Microseconds).UTC()
		endTimeFormatted := endTime.Format("15:04:05")

		if slot.BookedPatients.Int32 < slot.MaxPatients.Int32 {
			slotRes = timeSlotResponse{
				ID:        int32(slot.ID),
				StartTime: startTimeFormatted,
				EndTime:   endTimeFormatted,
				Status:    "available",
			}
		} else {
			slotRes = timeSlotResponse{
				ID:        int32(slot.ID),
				StartTime: startTimeFormatted,
				EndTime:   endTimeFormatted,
				Status:    "full",
			}
		}

		availableTimeSlots = append(availableTimeSlots, slotRes)
	}

	return availableTimeSlots, nil
}

// func (s *AppointmentService) CancelAppointment(ctx *gin.Context, appointmentID int64, username string) error {
//     // Lấy thông tin cuộc hẹn
//     appointment, err := s.storeDB.GetAppointmentByID(ctx, appointmentID)
//     if err != nil {
//         return fmt.Errorf("failed to get appointment: %w", err)
//     }

//     // Kiểm tra xem cuộc hẹn thuộc về người dùng
//     if appointment.Username != username {
//         return fmt.Errorf("unauthorized to cancel this appointment")
//     }

//     // Bắt đầu transaction
//     err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
//         // Cập nhật trạng thái lịch hẹn thành "cancelled"
//         err := q.UpdateAppointmentStatus(ctx, db.UpdateAppointmentStatusParams{
//             ID:     appointmentID,
//             Status: pgtype.Text{String: "cancelled", Valid: true},
//         })
//         if err != nil {
//             return fmt.Errorf("failed to update appointment status: %w", err)
//         }

//         // Nếu lịch hẹn đã được xác nhận, giảm số lượng booked_patients
//         if appointment.Status.String == "confirmed" {
//             err = q.DecreaseTimeSlotBookedPatients(ctx, db.DecreaseTimeSlotBookedPatientsParams{
//                 ID:       appointment.TimeSlotID.Int64,
//                 Date:     pgtype.Date{Time: appointment.Date.Time, Valid: true},
//                 DoctorID: int32(appointment.DoctorID.Int64),
//             })
//             if err != nil {
//                 return fmt.Errorf("failed to decrease booked patients: %w", err)
//             }
//         }

//         return nil
//     })
//     if err != nil {
//         return fmt.Errorf("transaction failed: %w", err)
//     }

//     return nil
// }
