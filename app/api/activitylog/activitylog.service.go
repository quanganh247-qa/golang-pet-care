package activitylog

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

type ActivityLogServiceInterface interface {
	CreateActivityLog(ctx *gin.Context, req createActivityLogRequest) (*createActivityLogResponse, error)
	GetActivityLogByID(ctx *gin.Context, logID int64) (*createActivityLogResponse, error)
	UpdateActivityLog(ctx *gin.Context, logID int64, req updateActivityLogRequest) error
	DeleteActivityLog(ctx *gin.Context, logID int64) error
	ListActivityLogs(ctx *gin.Context, petID int64, limit, offset int32) ([]createActivityLogResponse, error)
}

func (s *ActivityLogService) CreateActivityLog(ctx *gin.Context, req createActivityLogRequest) (*createActivityLogResponse, error) {
	var log createActivityLogResponse

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.CreateActivityLog(ctx, db.CreateActivityLogParams{
			Petid:        pgtype.Int8{Int64: req.PetID, Valid: true},
			Activitytype: req.ActivityType,
			Starttime:    pgtype.Timestamp{Time: req.StartTime, Valid: true},
			Duration:     pgtype.Interval{Microseconds: int64(req.Duration), Valid: true},
			Notes:        pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to create activity log: %w", err)
		}

		log = createActivityLogResponse{
			LogID:        res.Logid,
			PetID:        res.Petid.Int64,
			ActivityType: res.Activitytype,
			StartTime:    res.Starttime.Time,
			Duration:     res.Duration.Microseconds,
			Notes:        res.Notes.String,
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return &log, nil
}

func (s *ActivityLogService) GetActivityLogByID(ctx *gin.Context, logID int64) (*createActivityLogResponse, error) {
	var log createActivityLogResponse

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		res, err := q.GetActivityLogByID(ctx, logID)
		if err != nil {
			return fmt.Errorf("failed to get activity log: %w", err)
		}

		log = createActivityLogResponse{
			LogID:        res.Logid,
			PetID:        res.Petid.Int64,
			ActivityType: res.Activitytype,
			StartTime:    res.Starttime.Time,
			Duration:     res.Duration.Microseconds,
			Notes:        res.Notes.String,
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return &log, nil
}

func (s *ActivityLogService) ListActivityLogs(ctx *gin.Context, petID int64, limit, offset int32) ([]createActivityLogResponse, error) {
	var logs []createActivityLogResponse

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		params := db.ListActivityLogsParams{
			Petid:  pgtype.Int8{Int64: petID, Valid: true},
			Limit:  limit,
			Offset: offset,
		}

		res, err := q.ListActivityLogs(ctx, params)
		if err != nil {
			return fmt.Errorf("failed to list activity logs: %w", err)
		}

		for _, r := range res {
			logs = append(logs, createActivityLogResponse{
				LogID:        r.Logid,
				PetID:        r.Petid.Int64,
				ActivityType: r.Activitytype,
				StartTime:    r.Starttime.Time,
				Duration:     r.Duration.Microseconds,
				Notes:        r.Notes.String,
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return logs, nil
}

func (s *ActivityLogService) UpdateActivityLog(ctx *gin.Context, logID int64, req updateActivityLogRequest) error {
	params := db.UpdateActivityLogParams{
		Logid:        logID,
		Activitytype: req.ActivityType,
		Starttime:    pgtype.Timestamp{Time: req.StartTime, Valid: true},
		Duration:     pgtype.Interval{Microseconds: int64(req.Duration), Valid: true},
		Notes:        pgtype.Text{String: req.Notes, Valid: true},
	}
	return s.storeDB.UpdateActivityLog(ctx, params)
}

func (s *ActivityLogService) DeleteActivityLog(ctx *gin.Context, logID int64) error {
	return s.storeDB.DeleteActivityLog(ctx, logID)
}
