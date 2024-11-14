package petschedule

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetScheduleServiceInterface interface {
	CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) error
	GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResonse, error)
}

func (s *PetScheduleService) CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) error {

	if pet, err := s.storeDB.GetPetByID(ctx, petID); err != nil {
		return fmt.Errorf("Cannot find pet with ID %s: %w", pet.Name, err)
	}

	eventTime, _, err := util.ParseStringToTime(req.EventTime, "")
	if err != nil {
		return fmt.Errorf("parse event time%w", err)
	}

	// Implement logic to create a pet schedule
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		return q.CreatePetSchedule(ctx, db.CreatePetScheduleParams{
			PetID:        pgtype.Int8{Int64: petID, Valid: true},
			ScheduleType: req.ScheduleType,
			Duration:     pgtype.Text{String: req.Duration, Valid: true},
			EventTime:    pgtype.Timestamp{Time: eventTime, Valid: true},
			ActivityType: pgtype.Text{String: req.ActivityType, Valid: true},
			Frequency:    pgtype.Text{String: req.Frequency, Valid: true},
			Notes:        pgtype.Text{String: req.Notes, Valid: true},
		})
	})
	if err != nil {
		return fmt.Errorf("error creating pet schdule: ", err)
	}
	return nil
}

func (s *PetScheduleService) GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResonse, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	res, err := s.storeDB.GetAllSchedulesByPet(ctx, db.GetAllSchedulesByPetParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
		PetID:  pgtype.Int8{Int64: petID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching pet schedule: ", err)
	}

	var petSchedules []PetScheduleResonse
	for _, r := range res {

		petSchedules = append(petSchedules, PetScheduleResonse{
			ID:           r.PetID.Int64,
			ScheduleType: r.ScheduleType,
			Duration:     r.Duration.String,
			EventTime:    r.EventTime.Time.Format(time.RFC3339),
			ActivityType: r.ActivityType.String,
			Frequency:    r.Frequency.String,
			Notes:        r.Notes.String,
		})
	}

	return petSchedules, nil
}
