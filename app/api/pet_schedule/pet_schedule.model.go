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
}

type PetScheduleResponse struct {
	ID               int64  `json:"id"`
	PetID            int64  `json:"pet_id"`
	Title            string `json:"title"`
	ReminderDateTime string `json:"reminder_datetime"`
	EventRepeat      string `json:"event_repeat"`
	EndType          string `json:"end_type"`
	EndDate          string `json:"end_date"`
	Notes            string `json:"notes"`
	CreatedAt        string `json:"created_at"`
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
