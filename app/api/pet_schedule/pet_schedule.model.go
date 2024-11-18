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
	ScheduleType string `json:"schedule_type"`
	EventTime    string `json:"event_time"`
	Duration     string `json:"duration"`
	ActivityType string `json:"activity_type"`
	Frequency    string `json:"frequency"`
	Notes        string `json:"notes"`
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
