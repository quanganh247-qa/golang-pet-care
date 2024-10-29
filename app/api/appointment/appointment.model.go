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
<<<<<<< HEAD
	PetID      int64  `json:"pet_id"`
	DoctorID   int64  `json:"doctor_id"`
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
	ID          int64    `json:"id"`
	ServiceName string   `json:"service_name"`
	PetName     string   `json:"pet_name"`
	TimeSlot    timeslot `json:"time_slot"`
	DoctorName  string   `json:"doctor_name"`
	Note        string   `json:"note"`
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
=======
	petID      int64  `json:"petid"`
	doctorID   int64  `json:"doctorid"`
	date       string `json:"date"`
	timeSlotID int64  `json:"timeSlotID"`
	serviceID  int64  `json:"serviceID"`
	note       string `json:"note"`
}

type createAppointmentResponse struct {
	id          int64  `json:"id"`
	serviceName string `json:"serviceName"`
	date        string `json:"date"`
	timeSlot    string `json:"timeSlot"`
	doctorName  string `json:"doctorName"`
	note        string `json:"note"`
>>>>>>> c7f463c (update dtb)
}
