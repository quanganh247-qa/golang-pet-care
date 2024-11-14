package petschedule

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetScheduleServiceInterface interface {
	CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest) error
}

func (s *PetScheduleService) CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest) error {

	eventTime, _, err := util.ParseStringToTime(req.EventTime, "")
	if err != nil {
		return fmt.Errorf("parse event time")
	}

	// Implement logic to create a pet schedule
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		return q.CreatePetSchedule(ctx, db.CreatePetScheduleParams{
			ScheduleType: req.ScheduleType,
			EventTime:    pgtype.Timestamp{Time: eventTime, Valid: true},
			ActivityType: pgtype.Text{String: req.ActivityType, Valid: true},
			Frequency:    pgtype.Text{String: req.Frequency, Valid: true},
			Notes:        pgtype.Text{String: req.Notes, Valid: true},
		})
	})
	if err != nil {
		return fmt.Errorf("error creating pet schdule: ", err)
	}
	// Use the provided request to create the schedule in the database
	// Return an error if the schedule creation fails

	return nil
}
