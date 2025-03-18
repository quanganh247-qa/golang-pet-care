package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

type AppointmentController struct {
	service AppointmentServiceInterface
}

type AppointmentService struct {
	storeDB         db.Store
	taskDistributor worker.TaskDistributor
}

type AppointmentApi struct {
	controller AppointmentControllerInterface
}

type createAppointmentRequest struct {
	PetID      int64  `json:"pet_id"`
	DoctorID   int64  `json:"doctor_id"`
	Date       string `json:"date"`
	TimeSlotID int64  `json:"time_slot_id"`
	ServiceID  int64  `json:"service_id"`
	Reason     string `json:"reason"`
	StateID    int64  `json:"state_id"`
}

type timeslot struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type createAppointmentResponse struct {
	ID           int64    `json:"id"`
	PetName      string   `json:"pet_name"`
	ServiceName  string   `json:"service_name"`
	DoctorName   string   `json:"doctor_name"`
	Date         string   `json:"date"`
	TimeSlot     timeslot `json:"time_slot"`
	State        string   `json:"state"`
	Reason       string   `json:"reason"`
	ReminderSend bool     `json:"reminder_send"`
	CreatedAt    string   `json:"created_at"`
}
type timeSlotResponse struct {
	ID        int32  `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	// BookedPatients int32  `json:"booked_patients"`
	// MaxPatients    int32  `json:"max_patients"`
}

type updateAppointmentRequest struct {
	PaymentStatus string `json:"payment_status"`
	StateID       string `json:"state_id"`
	Notes         string `json:"notes"`
	ReminderSend  bool   `json:"reminder_send"`
}

type CreateSOAPRequest struct {
	DoctorID   int64  `json:"doctor_id"`
	Subjective string `json:"subjective"`
	Objective  string `json:"objective"`
	Assessment string `json:"assessment"`
	Plan       string `json:"plan"`
}

type UpdateSOAPRequest struct {
	Subjective string `json:"subjective"`
	Objective  string `json:"objective"`
	Assessment string `json:"assessment"`
	Plan       string `json:"plan"`
}

type SOAPResponse struct {
	ConsultationID int64  `json:"consultation_id"`
	AppointmentID  int64  `json:"appointment_id"`
	Subjective     string `json:"subjective"`
	Objective      string `json:"objective"`
	Assessment     string `json:"assessment"`
	Plan           string `json:"plan"`
	Notes          string `json:"notes"`
}
