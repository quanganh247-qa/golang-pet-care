package petschedule

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/llm"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetScheduleServiceInterface interface {
	CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) (*PetScheduleResponse, error)
	GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResponse, error)
	ListPetSchedulesByUsernameService(ctx *gin.Context, username string) ([]PetSchedules, error)
	ActivePetScheduleService(ctx *gin.Context, scheduleID int64, req ActiceRemider) error
	DeletePetScheduleService(ctx *gin.Context, scheduleID int64) error
	UpdatePetScheduleService(ctx *gin.Context, scheduleID int64, req UpdatePetScheduleRequest) (*PetScheduleResponse, error)
	ProcessSuggestionGemini(ctx *gin.Context, description string) (*llm.BaseResponse, error)
}

func (s *PetScheduleService) CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) (*PetScheduleResponse, error) {

	if pet, err := s.storeDB.GetPetByID(ctx, petID); err != nil {
		return nil, fmt.Errorf("Cannot find pet with ID %s: %w", pet.Name, err)
	}

	const iso8601Format = "2006-01-02T15:04:05Z"

	reminderTime, err := time.Parse(iso8601Format, req.ReminderDateTime)
	if err != nil {
		return nil, fmt.Errorf("invalid reminder date time format: %v", err)
	}

	var endDate pgtype.Date
	if req.EndDate != "" {
		// Try ISO 8601 format first (with time part)
		parsedEndDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
		endDate = pgtype.Date{Time: parsedEndDate, Valid: true}
	} else {
		endDate = pgtype.Date{Valid: false}
	}

	var res *PetScheduleResponse
	// Implement logic to create a pet schedule
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		schedule, err := q.CreatePetSchedule(ctx, db.CreatePetScheduleParams{
			PetID:            pgtype.Int8{Int64: petID, Valid: true},
			Title:            pgtype.Text{String: req.Title, Valid: true},
			EndType:          pgtype.Bool{Bool: req.EndType, Valid: true},
			ReminderDatetime: pgtype.Timestamp{Time: reminderTime, Valid: true},
			EventRepeat:      pgtype.Text{String: req.EventRepeat, Valid: true},
			EndDate:          endDate,
			Notes:            pgtype.Text{String: req.Notes, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error creating pet schedule: %w", err)
		}
		res = &PetScheduleResponse{
			ID:               schedule.ID,
			PetID:            schedule.PetID.Int64,
			Title:            schedule.Title.String,
			ReminderDateTime: schedule.ReminderDatetime.Time.Format(time.RFC3339),
			EventRepeat:      schedule.EventRepeat.String,
			EndType:          schedule.EndType.Bool,
			EndDate:          schedule.EndDate.Time.Format(time.RFC3339),
			Notes:            schedule.Notes.String,
			IsActive:         schedule.IsActive.Bool,
		}
		return nil

	})
	if err != nil {
		return nil, fmt.Errorf("error creating pet schdule: ", err)
	}

	// Invalidate cache after creating a new schedule
	middleware.InvalidateCache("pet_schedules")

	// Also clear the user's schedule cache if Redis is available
	if s.redis != nil {
		// Get the pet to find the username
		pet, err := s.storeDB.GetPetByID(ctx, petID)
		if err == nil {
			s.redis.ClearUserSchedulesCache(pet.Username)
		}
	}

	return res, nil
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
			ID:               r.ID,
			PetID:            r.PetID.Int64,
			Title:            r.Title.String,
			ReminderDateTime: r.ReminderDatetime.Time.Format(time.RFC3339),
			EventRepeat:      r.EventRepeat.String,
			EndType:          r.EndType.Bool,
			EndDate:          r.EndDate.Time.Format(time.RFC3339),
			Notes:            r.Notes.String,
			IsActive:         r.IsActive.Bool,
		})
	}

	return petSchedules, nil
}

func (s *PetScheduleService) ListPetSchedulesByUsernameService(ctx *gin.Context, username string) ([]PetSchedules, error) {
	// Try to get from cache first if Redis client is available
	if s.redis != nil {
		cachedSchedules, err := s.redis.PetSchedulesByUsernameLoadCache(username)
		if err == nil {
			// Cache hit, convert to response format
			var response []PetSchedules
			for petID, schedules := range cachedSchedules {
				// Get pet name (we could enhance by storing the pet name in the cache)
				pet, err := s.storeDB.GetPetByID(ctx, petID)
				if err != nil {
					// Skip this pet if we can't get its info
					continue
				}

				petSchedules := make([]PetScheduleResponse, 0, len(schedules))
				for _, schedule := range schedules {
					petSchedules = append(petSchedules, PetScheduleResponse{
						ID:               schedule.ID,
						PetID:            schedule.PetID,
						Title:            schedule.Title,
						ReminderDateTime: schedule.ReminderDateTime.Format(time.RFC3339),
						EventRepeat:      schedule.EventRepeat,
						EndType:          schedule.EndType,
						EndDate:          schedule.EndDate.Format(time.RFC3339),
						Notes:            schedule.Notes,
						IsActive:         schedule.IsActive,
					})
				}

				response = append(response, PetSchedules{
					PetID:     petID,
					PetName:   pet.Name,
					Schedules: petSchedules,
				})
			}
			// Log cache hit
			ctx.Set("cache_status", "HIT")
			ctx.Set("cache_source", "redis:schedules")
			return response, nil
		}
	}

	// Cache miss or no Redis client, get from database
	// Log cache miss
	ctx.Set("cache_status", "MISS")
	ctx.Set("cache_source", "db")

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
			PetName: schedule.PetName.String,
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
			IsActive:         schedule.IsActive.Bool,
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

// Active Pet Schedule
func (s *PetScheduleService) ActivePetScheduleService(ctx *gin.Context, scheduleID int64, req ActiceRemider) error {

	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.ActiveReminder(ctx, db.ActiveReminderParams{
			ID:       scheduleID,
			IsActive: pgtype.Bool{Bool: req.IsActive, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("error activating reminder: ", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error activating reminder: ", err)
	}
	return nil
}

// Delete Pet Schedule
func (s *PetScheduleService) DeletePetScheduleService(ctx *gin.Context, scheduleID int64) error {
	var schedule db.PetSchedule
	var err error
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		schedule, err = q.DeletePetSchedule(ctx, scheduleID)

		if err != nil {
			return fmt.Errorf("error deleting reminder: ", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting reminder: ", err)
	}
	s.redis.ClearPetSchedulesByPetCache(schedule.PetID.Int64)
	return nil
}

// Update Pet Schedule
func (s *PetScheduleService) UpdatePetScheduleService(ctx *gin.Context, scheduleID int64, req UpdatePetScheduleRequest) (*PetScheduleResponse, error) {

	var r db.UpdatePetScheduleParams

	schedule, err := s.storeDB.GetPetScheduleById(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("error getting reminder: ", err)
	}

	// validate req data
	if req.Title == "" {
		r.Title = schedule.Title
	} else {
		r.Title = pgtype.Text{String: req.Title, Valid: true}
	}

	if req.ReminderDateTime == "" {
		r.ReminderDatetime = schedule.ReminderDatetime
	} else {
		reminderTime, err := time.Parse("2006-01-02T15:04:05Z", req.ReminderDateTime)
		if err != nil {
			return nil, fmt.Errorf("invalid reminder date time format: %v", err)
		}
		r.ReminderDatetime = pgtype.Timestamp{Time: reminderTime, Valid: true}
	}

	if req.EventRepeat == "" {
		req.EventRepeat = schedule.EventRepeat.String
	} else {
		r.EventRepeat = pgtype.Text{String: req.EventRepeat, Valid: true}
	}

	// No need for empty string check for boolean values
	r.EndType = pgtype.Bool{Bool: req.EndType, Valid: true}

	if req.Notes == "" {
		r.Notes = schedule.Notes
	} else {
		r.Notes = pgtype.Text{String: req.Notes, Valid: true}
	}

	r.ID = scheduleID
	r.IsActive = pgtype.Bool{Bool: schedule.IsActive.Bool, Valid: true}
	if req.EndDate == "" {
		r.EndDate = schedule.EndDate
	} else {
		// Try ISO 8601 format first (with time part)
		parsedEndDate, err := time.Parse(time.RFC3339, req.EndDate)
		if err != nil {
			// If that fails, try simple date format
			parsedEndDate, err = time.Parse("2006-01-02", req.EndDate)
			if err != nil {
				return nil, fmt.Errorf("invalid end date format: %v", err)
			}
		}
		r.EndDate = pgtype.Date{Time: parsedEndDate, Valid: true}
	}

	var res *PetScheduleResponse
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		_, err := q.UpdatePetSchedule(ctx, r)
		if err != nil {
			return fmt.Errorf("error updating reminder: ", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error updating reminder: ", err)
	}

	s.redis.ClearPetSchedulesByPetCache(schedule.PetID.Int64)
	s.redis.ClearPetSchedulesCache()

	updateSchedule, err := s.storeDB.GetPetScheduleById(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("error getting reminder: ", err)
	}

	res = &PetScheduleResponse{
		ID:               updateSchedule.ID,
		PetID:            updateSchedule.PetID.Int64,
		Title:            updateSchedule.Title.String,
		ReminderDateTime: updateSchedule.ReminderDatetime.Time.Format(time.RFC3339),
		EventRepeat:      updateSchedule.EventRepeat.String,
		EndType:          updateSchedule.EndType.Bool,
		EndDate:          updateSchedule.EndDate.Time.Format(time.RFC3339),
		Notes:            updateSchedule.Notes.String,
		IsActive:         updateSchedule.IsActive.Bool,
	}
	return res, nil
}

func (s *PetScheduleService) ProcessSuggestionGemini(ctx *gin.Context, description string) (*llm.BaseResponse, error) {
	actionResponse, err := llm.DetermineActionGemini(ctx, s.config, description)
	if err != nil {
		return nil, fmt.Errorf("error determining action: %v", err)
	}

	var res *llm.BaseResponse

	switch actionResponse.Action {
	case llm.ActionAppointment:
		res, err = llm.GenerateSuggestionGemini(ctx, s.config, actionResponse.Action, description)
		if err != nil {
			return nil, fmt.Errorf("error generating suggestion: %v", err)
		}
		return res, nil
	case llm.ActionPetLog:
		res, err = llm.GenerateSuggestionGemini(ctx, s.config, actionResponse.Action, description)
		if err != nil {
			return nil, fmt.Errorf("error generating suggestion: %v", err)
		}
		return res, nil
	case llm.ActionPetSchedule:
		res, err = llm.GenerateSuggestionGemini(ctx, s.config, actionResponse.Action, description)
		if err != nil {
			return nil, fmt.Errorf("error generating suggestion: %v", err)
		}
		return res, nil
	}
	return res, nil
}
