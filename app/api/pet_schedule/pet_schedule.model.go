package petschedule

<<<<<<< HEAD
import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)
=======
import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
>>>>>>> 272832d (redis cache)

type PetScheduleApi struct {
	controller PetScheduleControllerInterface
}

type PetScheduleController struct {
	service PetScheduleServiceInterface
}

type PetScheduleService struct {
	storeDB db.Store
<<<<<<< HEAD
	config  *util.Config
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
=======
}

type PetScheduleRequest struct {
	ScheduleType string `json:"schedule_type"`
	EventTime    string `json:"event_time"`
	Duration     string `json:"duration"`
	ActivityType string `json:"activity_type"`
	Frequency    string `json:"frequency"`
	Notes        string `json:"notes"`
>>>>>>> 272832d (redis cache)
}

type PetScheduleResponse struct {
	ID           int64  `json:"id"`
	PetName      string `json:"pet_name"`
	ScheduleType string `json:"schedule_type"`
	EventTime    string `json:"event_time"`
	Duration     string `json:"duration"`
	ActivityType string `json:"activity_type"`
	Frequency    string `json:"frequency"`
	Notes        string `json:"notes"`
}

type PetSchedules struct {
	PetID     int64                 `json:"petid"`
	Schedules []PetScheduleResponse `json:"schedules"`
}
