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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 7cfffa9 (update dtb and appointment)
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
<<<<<<< HEAD
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
=======
>>>>>>> c7f463c (update dtb)
	petID      int64  `json:"petid"`
	doctorID   int64  `json:"doctorid"`
	date       string `json:"date"`
	timeSlotID int64  `json:"timeSlotID"`
	serviceID  int64  `json:"serviceID"`
	note       string `json:"note"`
}

type createAppointmentResponse struct {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> c7f463c (update dtb)
	id          int64  `json:"id"`
	serviceName string `json:"serviceName"`
	date        string `json:"date"`
	timeSlot    string `json:"timeSlot"`
	doctorName  string `json:"doctorName"`
	note        string `json:"note"`
<<<<<<< HEAD
>>>>>>> c7f463c (update dtb)
=======
=======
>>>>>>> 59d4ef2 (modify type of filed in dtb)
	id          int64       `json:"id"`
	serviceName string      `json:"serviceName"`
	petName     string      `json:"petName"`
	date        string      `json:"date"`
	timeSlot    db.Timeslot `json:"timeSlot"`
	doctorName  string      `json:"doctorName"`
	note        string      `json:"note"`
<<<<<<< HEAD
>>>>>>> 59d4ef2 (modify type of filed in dtb)
=======
=======
>>>>>>> 323513c (appointment api)
	id          int64
	serviceName string
	petName     string
	timeSlot    db.Timeslot
	doctorName  string
	note        string
<<<<<<< HEAD
>>>>>>> 323513c (appointment api)
=======
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
>>>>>>> 7cfffa9 (update dtb and appointment)
=======
>>>>>>> c7f463c (update dtb)
=======
>>>>>>> 59d4ef2 (modify type of filed in dtb)
=======
>>>>>>> 323513c (appointment api)
}
