package petschedule

import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"

type PetScheduleApi struct {
	controller PetScheduleControllerInterface
}

type PetScheduleController struct {
	service PetScheduleServiceInterface
}

type PetScheduleService struct {
	storeDB db.Store
}

type PetScheduleRequest struct {
	PetID            int64  `json:"pet_id"`
	Title            string `json:"title"`
	ReminderDateTime string `json:"reminder_datetime"`
	EventRepeat      string `json:"event_repeat"`
	EndType          string `json:"end_type"`
	EndDate          string `json:"end_date"`
	Notes            string `json:"notes"`
	IsActive         bool   `json:"is_active"`
}

type PetScheduleResponse struct {
	ID               int64  `json:"id"`
	PetID            int64  `json:"pet_id"`
	Title            string `json:"title"`
	ReminderDateTime string `json:"reminder_datetime"`
	EventRepeat      string `json:"event_repeat"`
	EndType          bool   `json:"end_type"`
	EndDate          string `json:"end_date"`
	Notes            string `json:"notes"`
	CreatedAt        string `json:"created_at"`
	IsActive         bool   `json:"is_active"`
}
type PetKey struct {
	PetID   int64  `json:"pet_id"`
	PetName string `json:"pet_name"`
}

type PetSchedules struct {
	PetID     int64                 `json:"pet_id"`
	PetName   string                `json:"pet_name"`
	Schedules []PetScheduleResponse `json:"schedules"`
}

type ActiceRemider struct {
	IsActive bool `json:"is_active"`
}

type ScheduleSuggestion struct {
	Voice string `json:"voice"`
}


// type ScheduleSuggestionResponse struct {
// 	PetName        string `json:"pet_name"`
// 	PetType        string `json:"pet_type"`
// 	Activity       string `json:"activity"`
// 	ReminderTime   string `json:"reminder_time"`
// 	EventRepeat    string `json:"event_repeat"`
// 	Notes          string `json:"notes"`
// }

// - Response:
// {
// 	"action": "appointment",
// 	"parameters": {
// 	"pet_name": "Buddy",
// 	"appointment_type": "vet",
// 	"date": "2023-10-15",
// 	"time": "14:00"
// 	}


type ActionResponse struct {
	Action string `json:"action"`
}

const (
	ActionAppointment = "appointment"
	ActionPetLog      = "pet_log"
	ActionPetSchedule = "pet_schedule"
)

type AppointmentResponse struct {
	Action string `json:"action"`
	Parameters struct {
		PetName        string `json:"pet_name"`
		AppointmentType string `json:"appointment_type"`
		Date           string `json:"date"`
		Time           string `json:"time"`
	} `json:"parameters"`
}

type LogResponse struct {
	Action string `json:"action"`
	Parameters struct {
		PetName string `json:"pet_name"`
		Activity string `json:"activity"`
		Date     string `json:"date"`
		Time     string `json:"time"`
		Notes    string `json:"notes"`
	} `json:"parameters"`
}

type ScheduleResponse struct {
	Action string `json:"action"`
	Parameters struct {
		PetName string `json:"pet_name"`
		Activity string `json:"activity"`
		Date     string `json:"date"`
		Time     string `json:"time"`
		Notes    string `json:"notes"`
	} `json:"parameters"`
}

// SuggestionResponse is the interface that all response types must implement
type SuggestionResponse interface {
    GetAction() string
    Validate() error
}

// BaseResponse contains common fields for all response types
type BaseResponse struct {
    Action     string                 `json:"action"`
    Parameters map[string]interface{} `json:"parameters"`
}

func (b BaseResponse) GetAction() string {
    return b.Action
}