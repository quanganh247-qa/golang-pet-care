package petschedule

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetScheduleServiceInterface interface {
	CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) error
	GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResponse, error)
	ListPetSchedulesByUsernameService(ctx *gin.Context, username string) ([]PetSchedules, error)
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

func (s *PetScheduleService) GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResponse, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	res, err := s.storeDB.GetAllSchedulesByPet(ctx, db.GetAllSchedulesByPetParams{
		Limit:  int32(pagination.PageSize),
		Offset: int32(offset),
		PetID:  pgtype.Int8{Int64: petID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching pet schedule: ", err)
	}

	var petSchedules []PetScheduleResponse
	for _, r := range res {

		petSchedules = append(petSchedules, PetScheduleResponse{
			ID:           r.ID,
			PetID:        r.PetID.Int64,
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

func (s *PetScheduleService) ListPetSchedulesByUsernameService(ctx *gin.Context, username string) ([]PetSchedules, error) {
	schedules, err := s.storeDB.ListPetSchedulesByUsername(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list pet schedules"})
		return nil, err
	}

	// Group schedules by pet ID
	groupedSchedules := make(map[PetKey][]PetScheduleResponse)
	for _, schedule := range schedules {
		petKey := PetKey{
			PetID:   schedule.PetID.Int64,
			PetName: schedule.Name.String,
		}
		groupedSchedules[petKey] = append(groupedSchedules[petKey], PetScheduleResponse{
			ID:           schedule.ID,
			PetID:        schedule.PetID.Int64,
			EventTime:    schedule.EventTime.Time.Format(time.RFC3339),
			ScheduleType: schedule.ScheduleType,
			ActivityType: schedule.ActivityType.String,
			Duration:     schedule.Duration.String,
			Frequency:    schedule.Frequency.String,
			Notes:        schedule.Notes.String,
		})

	}

	// Convert the map to a slice of responses
	var response []PetSchedules
	for petKey, schedules := range groupedSchedules {
		response = append(response, PetSchedules{
			PetID:     petKey.PetID,
			PetName:   petKey.PetName,
			Schedules: schedules,
		})
	}

	return response, nil
}
