package appointment

import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"

type AppointmentController struct {
	service AppointmentServiceInterface
}

type AppointmentService struct {
	storeDB db.Store
}

type AppointmentApi struct {
	controller AppointmentControllerInterface
}

type createAppointmentRequest struct {
	PetID      int64  `json:"pet_id"`
	DoctorID   string `json:"doctor_id"`
	Date       string `json:"date"`
	TimeSlotID int64  `json:"time_slot_id"`
	ServiceID  int64  `json:"service_id"`
	Note       string `json:"note"`
}

type timeslot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	// Date      string `json:"date"`
}

type createAppointmentResponse struct {
	ID          int64  `json:"id"`
	ServiceName string `json:"service_name"`
	PetName     string `json:"pet_name"`
	DoctorName  string `json:"doctor_name"`
	Note        string `json:"note"`
}

type AppointmentWithDetails struct {
	AppointmentID int64  `json:"appointment_id"`
	PetName       string `json:"pet_name"`
	ServiceName   string `json:"service_name"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
}

type updateAppointmentStatusRequest struct {
	Status string `json:"status"`
}
