package petschedule

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)
<<<<<<< HEAD
=======
import db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
>>>>>>> 272832d (redis cache)
=======
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
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
<<<<<<< HEAD
=======
}

type PetScheduleRequest struct {
>>>>>>> 272832d (redis cache)
	ScheduleType string `json:"schedule_type"`
	EventTime    string `json:"event_time"`
	Duration     string `json:"duration"`
	ActivityType string `json:"activity_type"`
	Frequency    string `json:"frequency"`
	Notes        string `json:"notes"`
<<<<<<< HEAD
>>>>>>> 272832d (redis cache)
=======
	PetID            int64  `json:"pet_id"`
	Title            string `json:"title"`
	ReminderDateTime string `json:"reminder_datetime"`
	EventRepeat      string `json:"event_repeat"`
	EndType          string `json:"end_type"`
	EndDate          string `json:"end_date"`
	Notes            string `json:"notes"`
>>>>>>> 3835eb4 (update pet_schedule api)
=======
>>>>>>> 4c66ef3 (feat: update schedule API)
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
>>>>>>> 272832d (redis cache)
}
