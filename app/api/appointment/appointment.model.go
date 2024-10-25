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
	Username string `json:"username"`
	Petid    int64  `json:"petid"`
	Doctorid int64  `json:"doctorid"`
	Date     string `json:"date"`
	Time     string `json:"time"`
}
