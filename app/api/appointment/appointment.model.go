package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
=======
>>>>>>> e30b070 (Get list appoinment by user)
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

type Pet struct {
	PetID    int64  `json:"pet_id"`
	PetName  string `json:"pet_name"`
	PetBreed string `json:"pet_breed"`
}

type Owner struct {
	OwnerName    string `json:"owner_name"`
	OwnerPhone   string `json:"owner_phone"`
	OwnerEmail   string `json:"owner_email"`
	OwnerAddress string `json:"owner_address"`
}

type Serivce struct {
	ServiceName     string `json:"service_name"`
	ServiceDuration int16  `json:"service_duration"`
}

type Doctor struct {
	DoctorID   int64  `json:"doctor_id"`
	DoctorName string `json:"doctor_name"`
}

type Appointment struct {
	ID           int64    `json:"id"`
	Pet          Pet      `json:"pet"`
	Owner        Owner    `json:"owner"`
	Serivce      Serivce  `json:"service"`
	Doctor       Doctor   `json:"doctor"`
	Room         string   `json:"room"`
	Date         string   `json:"date"`
	TimeSlot     timeslot `json:"time_slot"`
	State        string   `json:"state"`
	Reason       string   `json:"reason"`
	ReminderSend bool     `json:"reminder_send"`
	CreatedAt    string   `json:"created_at"`
}

type createAppointmentRequest struct {
	PetID      int64  `json:"pet_id"`
	DoctorID   string `json:"doctor_id"`
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

type historyAppointmentResponse struct {
	ID          int64  `json:"id"`
	PetName     string `json:"pet_name"`
	Reason      string `json:"reason"`
	Date        string `json:"date"`
	ServiceName string `json:"service_name"`
	ArrivalTime string `json:"arrival_time"`
	DoctorName  string `json:"doctor_name"`
	Room        string `json:"room"`
}

type createAppointmentResponse struct {
<<<<<<< HEAD
	ID           int64    `json:"id"`
	DoctorName   string   `json:"doctor_name"`
	PetName      string   `json:"pet_name"`
	Reason       string   `json:"reason"`
	Date         string   `json:"date"`
	ServiceName  string   `json:"service_name"`
	TimeSlot     timeslot `json:"time_slot"`
	State        string   `json:"state"`
	ReminderSend bool     `json:"reminder_send"`
	CreatedAt    string   `json:"created_at"`
	RoomType     string   `json:"room_type"`
}
type timeSlotResponse struct {
	ID        int32  `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
=======
	ID          int64  `json:"id"`
	ServiceName string `json:"service_name"`
	PetName     string `json:"pet_name"`
	Date        string `json:"date"`
	DoctorName  string `json:"doctor_name"`
	Note        string `json:"note"`
>>>>>>> cfbe865 (updated service response)
}

<<<<<<< HEAD
type updateAppointmentRequest struct {
	PaymentStatus     *string `json:"payment_status"`
	StateID           *int32  `json:"state_id"`
	RoomID            *int32  `json:"room_id"`
	Notes             *string `json:"notes"`
	AppointmentReason *string `json:"appointment_reason"`
	ReminderSend      *bool   `json:"reminder_send"`
	ArrivalTime       *string `json:"arrival_time"`
	Priority          *string `json:"priority"`
=======
type AppointmentWithDetails struct {
	AppointmentID int64  `json:"appointment_id"`
	PetName       string `json:"pet_name"`
	ServiceName   string `json:"service_name"`
<<<<<<< HEAD
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
	DoctorID      int64  `json:"doctor_id"`
	ServiceID     int64  `json:"service_id"`
	Date          string `json:"date"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
	ReminderSend  bool   `json:"reminder_send"`
	CreatedAt     string `json:"created_at"`
>>>>>>> e30b070 (Get list appoinment by user)
=======
	// StartTime     string `json:"start_time"`
	// EndTime       string `json:"end_time"`
	DoctorName   string `json:"doctor_name"`
	Date         string `json:"date"`
	Status       string `json:"status"`
	Notes        string `json:"notes"`
	ReminderSend bool   `json:"reminder_send"`
	CreatedAt    string `json:"created_at"`
>>>>>>> 98e9e45 (ratelimit and recovery function)
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
	// Plan           string `json:"plan"`
	Notes string `json:"notes"`
}

type QueueItem struct {
	ID              int64  `json:"id"`
	PatientName     string `json:"patientName"`
	Status          string `json:"status"`
	Priority        string `json:"priority"`
	AppointmentType string `json:"appointmentType"`
	Doctor          string `json:"doctor"`
	WaitingSince    string `json:"waitingSince"`
	ActualWaitTime  string `json:"actualWaitTime"`
}
