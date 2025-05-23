package petschedule

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetScheduleApi struct {
	controller PetScheduleControllerInterface
}

type PetScheduleController struct {
	service PetScheduleServiceInterface
}

type PetScheduleService struct {
	storeDB db.Store
	config  *util.Config
	redis   *redis.ClientType
}

type PetScheduleRequest struct {
	PetID            int64  `json:"pet_id"`
	Title            string `json:"title"`
	ReminderDateTime string `json:"reminder_datetime"`
	EventRepeat      string `json:"event_repeat"`
	EndType          bool   `json:"end_type"`
	EndDate          string `json:"end_date"`
	Notes            string `json:"notes"`
	IsActive         bool   `json:"is_active"`
}

type UpdatePetScheduleRequest struct {
	Title            string `json:"title"`
	ReminderDateTime string `json:"reminder_datetime"`
	EventRepeat      string `json:"event_repeat"`
	EndType          bool   `json:"end_type"`
	EndDate          string `json:"end_date"`
	Notes            string `json:"notes"`
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
