package appointment

import (
	"sync"
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

// NotificationManager xử lý việc lưu trữ và phân phối thông báo
type NotificationManager struct {
	notifications       map[string][]AppointmentNotification // map[username][]Notification
	pendingClients      map[string]chan AppointmentNotification
	userRoles           map[string]string                    // map[username]role
	missedNotifications map[string][]AppointmentNotification // Lưu trữ thông báo khi client không online
	notificationMutex   sync.RWMutex
	clientMutex         sync.RWMutex
	roleMutex           sync.RWMutex
	missedMutex         sync.RWMutex
}

// NewNotificationManager tạo một đối tượng quản lý thông báo mới
func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		notifications:       make(map[string][]AppointmentNotification),
		pendingClients:      make(map[string]chan AppointmentNotification),
		userRoles:           make(map[string]string),
		missedNotifications: make(map[string][]AppointmentNotification),
	}
}

type AppointmentController struct {
	service AppointmentServiceInterface
}

type AppointmentService struct {
	storeDB             db.Store
	taskDistributor     worker.TaskDistributor
	ws                  *websocket.WSClientManager
	notificationManager *NotificationManager
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
	OwnerPhone   string `json:"owner_number"`
	OwnerEmail   string `json:"owner_email"`
	OwnerAddress string `json:"owner_address"`
}

type Serivce struct {
	ServiceName     string  `json:"service_name"`
	ServiceDuration int16   `json:"service_duration"`
	ServiceAmount   float64 `json:"service_amount"`
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
	ArrivalTime  string   `json:"arrival_time"`
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

type historyAppointmentResponse struct {
	ID          int64  `json:"id"`
	PetName     string `json:"pet_name"`
	Reason      string `json:"reason"`
	Date        string `json:"date"`
	ServiceName string `json:"service_name"`
	ArrivalTime string `json:"arrival_time"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
	DoctorName  string `json:"doctor_name"`
	Room        string `json:"room"`
}

type createAppointmentResponse struct {
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
}

type updateAppointmentRequest struct {
	PaymentStatus     *string `json:"payment_status"`
	StateID           *int32  `json:"state_id"`
	RoomID            *int32  `json:"room_id"`
	Notes             *string `json:"notes"`
	AppointmentReason *string `json:"appointment_reason"`
	ReminderSend      *bool   `json:"reminder_send"`
	ArrivalTime       *string `json:"arrival_time"`
	Priority          *string `json:"priority"`
}

type CreateSOAPRequest struct {
	DoctorID   int64  `json:"doctor_id"`
	Subjective string `json:"subjective"`
}

type UpdateSOAPRequest struct {
	Subjective string        `json:"subjective"`
	Objective  ObjectiveData `json:"objective"`
	Assessment string        `json:"assessment"`
	Plan       int8          `json:"plan"`
}

type SOAPResponse struct {
	ConsultationID int64         `json:"consultation_id"`
	AppointmentID  int64         `json:"appointment_id"`
	Subjective     string        `json:"subjective"`
	Objective      ObjectiveData `json:"objective"`
	Assessment     string        `json:"assessment"`
	// Plan           string `json:"plan"`
	Notes string `json:"notes"`
}

// ObjectiveData đại diện cho dữ liệu trong phần Objective
type ObjectiveData struct {
	VitalSigns VitalSignsData `json:"vital_signs"` // Dấu hiệu sinh tồn
	Systems    SystemsData    `json:"systems"`     // Hệ thống cơ thể
}

// VitalSignsData chứa các dấu hiệu sinh tồn
type VitalSignsData struct {
	Weight          string `json:"weight"`           // Cân nặng (kg)
	Temperature     string `json:"temperature"`      // Nhiệt độ (°C)
	HeartRate       string `json:"heart_rate"`       // Nhịp tim (bpm)
	RespiratoryRate string `json:"respiratory_rate"` // Nhịp thở (rpm)
	GeneralNotes    string `json:"general_notes"`    // Ghi chú chung
}

// SystemsData chứa dữ liệu khám theo hệ thống
type SystemsData struct {
	Cardiovascular   string `json:"cardiovascular"`   // Tim mạch
	Respiratory      string `json:"respiratory"`      // Hô hấp
	Gastrointestinal string `json:"gastrointestinal"` // Tiêu hóa
	Musculoskeletal  string `json:"musculoskeletal"`  // Cơ xương
	Neurological     string `json:"neurological"`     // Thần kinh
	Skin             string `json:"skin"`             // Da/lông
	Eyes             string `json:"eyes"`             // Mắt
	Ears             string `json:"ears"`             // Tai
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

type CancelOrderRequest struct {
	OrderID int    `json:"order_id" binding:"required"`
	Reason  string `json:"reason"`
}

type UpdateTestStatusRequest struct {
	TestID int    `json:"test_id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=pending processing completed cancelled"`
}

type AppointmentDistribution struct {
	ServiceID        int64   `json:"service_id"`
	ServiceName      string  `json:"service_name"`
	AppointmentCount int64   `json:"appointment_count"`
	Percentage       float64 `json:"percentage"`
}

type AppointmentDistributionRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request parameters"`
}

// SuccessResponse represents a success response
// @Description Success response structure
type SuccessResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

// CreateAppointmentRequest represents the request payload for creating an appointment
// @Description Appointment creation request
type CreateAppointmentRequest struct {
	PetID     string    `json:"pet_id" example:"pet-123"`
	DoctorID  string    `json:"doctor_id" example:"doctor-456"`
	StartTime time.Time `json:"start_time" example:"2023-04-01T10:00:00Z"`
	EndTime   time.Time `json:"end_time" example:"2023-04-01T11:00:00Z"`
	Notes     string    `json:"notes" example:"Regular check-up"`
	Status    string    `json:"status" example:"scheduled"`
}

// AppointmentResponse represents the response payload for appointment data
// @Description Appointment response data
type AppointmentResponse struct {
	ID        string    `json:"id" example:"appt-789"`
	PetID     string    `json:"pet_id" example:"pet-123"`
	DoctorID  string    `json:"doctor_id" example:"doctor-456"`
	StartTime time.Time `json:"start_time" example:"2023-04-01T10:00:00Z"`
	EndTime   time.Time `json:"end_time" example:"2023-04-01T11:00:00Z"`
	Notes     string    `json:"notes" example:"Regular check-up"`
	Status    string    `json:"status" example:"scheduled"`
	CreatedAt time.Time `json:"created_at" example:"2023-03-30T09:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-03-30T09:00:00Z"`
}

type createWalkInAppointmentRequest struct {
	PetID     int64   `json:"pet_id"`
	DoctorID  int64   `json:"doctor_id"`
	ServiceID int64   `json:"service_id"`
	Reason    string  `json:"reason"`
	Priority  string  `json:"priority"`
	Owner     *Owner  `json:"owner,omitempty"` // Optional owner information for new users
	Pet       *NewPet `json:"pet,omitempty"`   // Optional pet information for new pets
}

// NewPet represents the information needed to create a new pet
type NewPet struct {
	Name      string  `json:"name"`
	Breed     string  `json:"breed"`
	Species   string  `json:"species"`
	BirthDate string  `json:"birth_date"`
	Gender    string  `json:"gender"`
	Weight    float64 `json:"weight"`
	Age       int     `json:"age"`
}

type AppointmentNotification struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	AppointmentID int64    `json:"appointment_id"`
	Pet           Pet      `json:"pet"`
	Doctor        Doctor   `json:"doctor"`
	Reason        string   `json:"reason"`
	Date          string   `json:"date"`
	TimeSlot      timeslot `json:"time_slot"`
	ServiceName   string   `json:"service_name"`
}

// DatabaseNotification đại diện cho thông báo được lưu trong cơ sở dữ liệu
type DatabaseNotification struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	NotifyType  string    `json:"notify_type"`
	RelatedID   int64     `json:"related_id"`
	RelatedType string    `json:"related_type"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
