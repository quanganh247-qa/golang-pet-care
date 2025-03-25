package petschedule

import (
	"fmt"
	"log"
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

	const iso8601Format = "2006-01-02T15:04:05Z"

	reminderTime, err := time.Parse(iso8601Format, req.ReminderDateTime)
	if err != nil {
		log.Fatalf("invalid start date format: %v", err)
	}

	var endDate pgtype.Date
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			log.Fatalf("invalid end date format: %v", err)
		}
		endDate = pgtype.Date{Time: parsedEndDate, Valid: true}
	} else {
		endDate = pgtype.Date{Valid: false}
	}
	// Implement logic to create a pet schedule
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		return q.CreatePetSchedule(ctx, db.CreatePetScheduleParams{
			PetID:            pgtype.Int8{Int64: petID, Valid: true},
			Title:            pgtype.Text{String: req.Title, Valid: true},
			ReminderDatetime: pgtype.Timestamp{Time: reminderTime, Valid: true},
			EventRepeat:      pgtype.Text{String: req.EventRepeat, Valid: true},
			EndDate:          endDate,
			Notes:            pgtype.Text{String: req.Notes, Valid: true},
		})
	})
	if err != nil {
		return fmt.Errorf("error creating pet schdule: ", err)
	}
	return nil
}

func (s *PetScheduleService) GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResponse, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	fmt.Println(petID)

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
			ID:               r.ID,
			PetID:            r.PetID.Int64,
			Title:            r.Title.String,
			ReminderDateTime: r.ReminderDatetime.Time.Format(time.RFC3339),
			EventRepeat:      r.EventRepeat.String,
			EndType:          r.EndType.Bool,
			EndDate:          r.EndDate.Time.Format(time.RFC3339),
			Notes:            r.Notes.String,
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
			ID:               schedule.ID,
			PetID:            schedule.PetID.Int64,
			Title:            schedule.Title.String,
			ReminderDateTime: schedule.ReminderDatetime.Time.Format(time.RFC3339),
			EventRepeat:      schedule.EventRepeat.String,
			EndType:          schedule.EndType.Bool,
			EndDate:          schedule.EndDate.Time.Format(time.RFC3339),
			Notes:            schedule.Notes.String,
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
