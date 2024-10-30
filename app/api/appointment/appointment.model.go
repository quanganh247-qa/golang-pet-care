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
	petID      int64  `json:"petid"`
	doctorID   int64  `json:"doctorid"`
	date       string `json:"date"`
	timeSlotID int64  `json:"timeSlotID"`
	serviceID  int64  `json:"serviceID"`
	note       string `json:"note"`
}

type createAppointmentResponse struct {
	id          int64
	serviceName string
	petName     string
	timeSlot    db.Timeslot
	doctorName  string
	note        string
}
