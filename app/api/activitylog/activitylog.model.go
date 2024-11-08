package activitylog

import (
	"time"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ActivityLogApi struct {
	controller ActivityLogControllerInterface
}

type ActivityLogController struct {
	service ActivityLogServiceInterface
}

type ActivityLogService struct {
	storeDB db.Store
}

type createActivityLogRequest struct {
	PetID        int64     `json:"petID"`
	ActivityType string    `json:"activityType"`
	StartTime    time.Time `json:"startTime"`
	Duration     int64     `json:"duration"`
	Notes        string    `json:"notes"`
}

type createActivityLogResponse struct {
	LogID        int64     `json:"logID"`
	PetID        int64     `json:"petID"`
	ActivityType string    `json:"activityType"`
	StartTime    time.Time `json:"startTime"`
	Duration     int64     `json:"duration"`
	Notes        string    `json:"notes"`
}

type updateActivityLogRequest struct {
	ActivityType string    `json:"activityType"`
	StartTime    time.Time `json:"startTime"`
	Duration     int64     `json:"duration"`
	Notes        string    `json:"notes"`
}
