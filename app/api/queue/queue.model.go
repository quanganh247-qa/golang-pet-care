package queue

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type QueueController struct {
	service QueueServiceInterface
}

type QueueService struct {
	storeDB db.Store
}

type QueueApi struct {
	controller QueueControllerInterface
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
