package petschedule

import (
	"encoding/json"
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"log"
>>>>>>> 3835eb4 (update pet_schedule api)
=======
>>>>>>> eb8d761 (updated pet schedule)
	"net/http"
	"strconv"
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"net/http"
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> 4c66ef3 (feat: update schedule API)
=======
	"strings"
>>>>>>> ffc9071 (AI suggestion)
	"time"
=======
>>>>>>> 272832d (redis cache)
=======
	"time"
>>>>>>> e01abc5 (pet schedule api)

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/llm"
=======
>>>>>>> 272832d (redis cache)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/llm"
>>>>>>> ffc9071 (AI suggestion)
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type PetScheduleServiceInterface interface {
<<<<<<< HEAD
<<<<<<< HEAD
	CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) error
	GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResponse, error)
	ListPetSchedulesByUsernameService(ctx *gin.Context, username string) ([]PetSchedules, error)
	ActivePetScheduleService(ctx *gin.Context, scheduleID int64, req ActiceRemider) error
	DeletePetScheduleService(ctx *gin.Context, scheduleID int64) error
<<<<<<< HEAD
<<<<<<< HEAD
	UpdatePetScheduleService(ctx *gin.Context, scheduleID int64, req PetScheduleRequest) error
<<<<<<< HEAD
	ProcessSuggestionGemini(ctx *gin.Context, description string) (*llm.BaseResponse, error)
=======
	UpdatePetScheduleService(ctx *gin.Context, scheduleID int64, req PetScheduleRequest) error
>>>>>>> 4c66ef3 (feat: update schedule API)
=======
	ProcessSuggestion(ctx *gin.Context, description string) (*BaseResponse, error)
>>>>>>> ffc9071 (AI suggestion)
}

func (s *PetScheduleService) CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) error {

	if pet, err := s.storeDB.GetPetByID(ctx, petID); err != nil {
		return fmt.Errorf("Cannot find pet with ID %s: %w", pet.Name, err)
	}

	const iso8601Format = "2006-01-02T15:04:05Z"

	reminderTime, err := time.Parse(iso8601Format, req.ReminderDateTime)
	if err != nil {
		return fmt.Errorf("invalid reminder date time format: %v", err)
	}

	var endDate pgtype.Date
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date format: %v", err)
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
=======
	CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest) error
=======
	CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) error
<<<<<<< HEAD
	GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResonse, error)
>>>>>>> e01abc5 (pet schedule api)
=======
	GetAllSchedulesByPetService(ctx *gin.Context, petID int64, pagination *util.Pagination) ([]PetScheduleResponse, error)
	ListPetSchedulesByUsernameService(ctx *gin.Context, username string) ([]PetSchedules, error)
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> eb8d761 (updated pet schedule)
}

func (s *PetScheduleService) CreatePetScheduleService(ctx *gin.Context, req PetScheduleRequest, petID int64) error {

	if pet, err := s.storeDB.GetPetByID(ctx, petID); err != nil {
		return fmt.Errorf("Cannot find pet with ID %s: %w", pet.Name, err)
	}

	const iso8601Format = "2006-01-02T15:04:05Z"

	reminderTime, err := time.Parse(iso8601Format, req.ReminderDateTime)
	if err != nil {
		return fmt.Errorf("invalid reminder date time format: %v", err)
	}

	var endDate pgtype.Date
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end date format: %v", err)
		}
		endDate = pgtype.Date{Time: parsedEndDate, Valid: true}
	} else {
		endDate = pgtype.Date{Valid: false}
	}
	// Implement logic to create a pet schedule
	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		return q.CreatePetSchedule(ctx, db.CreatePetScheduleParams{
<<<<<<< HEAD
			PetID:        pgtype.Int8{Int64: petID, Valid: true},
			ScheduleType: req.ScheduleType,
			Duration:     pgtype.Text{String: req.Duration, Valid: true},
			EventTime:    pgtype.Timestamp{Time: eventTime, Valid: true},
			ActivityType: pgtype.Text{String: req.ActivityType, Valid: true},
			Frequency:    pgtype.Text{String: req.Frequency, Valid: true},
			Notes:        pgtype.Text{String: req.Notes, Valid: true},
>>>>>>> 272832d (redis cache)
=======
			PetID:            pgtype.Int8{Int64: petID, Valid: true},
			Title:            pgtype.Text{String: req.Title, Valid: true},
			ReminderDatetime: pgtype.Timestamp{Time: reminderTime, Valid: true},
			EventRepeat:      pgtype.Text{String: req.EventRepeat, Valid: true},
			EndDate:          endDate,
			Notes:            pgtype.Text{String: req.Notes, Valid: true},
>>>>>>> 3835eb4 (update pet_schedule api)
		})
	})
	if err != nil {
		return fmt.Errorf("error creating pet schdule: ", err)
	}
<<<<<<< HEAD
<<<<<<< HEAD
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
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.DeletePetSchedule(ctx, scheduleID)

		if err != nil {
			return fmt.Errorf("error deleting reminder: ", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting reminder: ", err)
	}
	return nil
}

// Update Pet Schedule
func (s *PetScheduleService) UpdatePetScheduleService(ctx *gin.Context, scheduleID int64, req PetScheduleRequest) error {

	var r db.UpdatePetScheduleParams

	schedule, err := s.storeDB.GetPetScheduleById(ctx, scheduleID)
	if err != nil {
		return fmt.Errorf("error getting reminder: ", err)
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
			return fmt.Errorf("invalid reminder date time format: %v", err)
		}
		r.ReminderDatetime = pgtype.Timestamp{Time: reminderTime, Valid: true}
	}

	if req.EventRepeat == "" {
		req.EventRepeat = schedule.EventRepeat.String
	} else {
		r.EventRepeat = pgtype.Text{String: req.EventRepeat, Valid: true}
	}

	if req.EndType == "" {
		r.EndType = schedule.EndType
	} else {
		endType, err := strconv.ParseBool(req.EndType)
		if err != nil {
			return fmt.Errorf("invalid end type format: %v", err)
		}
		r.EndType = pgtype.Bool{Bool: endType, Valid: true}
	}

	if req.Notes == "" {
		r.Notes = schedule.Notes
	} else {
		r.Notes = pgtype.Text{String: req.Notes, Valid: true}
	}

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.UpdatePetSchedule(ctx, r)
		if err != nil {
			return fmt.Errorf("error updating reminder: ", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating reminder: ", err)
	}
	return nil
}
<<<<<<< HEAD
<<<<<<< HEAD

func (s *PetScheduleService) ProcessSuggestionGemini(ctx *gin.Context, description string) (*llm.BaseResponse, error) {
	actionResponse, err := llm.DetermineActionGemini(ctx, s.config, description)
=======

func GenerateSuggestion(ctx *gin.Context,action, description string) (*BaseResponse, error) {
	prompt := fmt.Sprintf(`As an AI assistant, your task is to interpret the user's request and determine the appropriate action to take. The possible actions are:
		- "appointment": For scheduling a one-time appointment or event.
		- "pet_log": For logging a one-time activity or event related to a pet.
		- "pet_schedule": For setting up a recurring schedule (e.g., daily, weekly).
		- "unknown": If the request doesn't match any of the above or is unclear.

		Based on the user's description, identify the action type and extract the necessary parameters. Respond with a JSON object containing:

		- "action": The determined action type.
		- "parameters": A dictionary of extracted parameters relevant to the action.

		Rules:
		- For "appointment", extract "pet_name", "appointment_type", "date" (in YYYY-MM-DD format, using current year for relative dates), and "time" (in HH:MM format).
		- For "pet_log", extract "pet_name", "activity", and "date" (in YYYY-MM-DD format).
		- For "pet_schedule", extract "pet_name", "activity", "frequency" (e.g., "daily", "weekly"), and "start_date" (in YYYY-MM-DD format).
		- If the action is "unknown", provide a brief explanation in "parameters" under "reason".
		- Always use the current time %s for any dates mentioned without a year.

		Example:
		- User: "Schedule a vet appointment for Buddy today at 2 PM."
		- Response:
		{
			"action": "appointment",
			"parameters": {
				"pet_name": "Buddy",
				"appointment_type": "vet",
				"date": "%s",
				"time": "14:00"
			}
		}

		Analyze the following request and provide the JSON response:

		Request: %s	with action: %s

`, time.Now(), time.Now().Format("2006-01-02"), description, action)

	reqBody := llm.OllamaRequest{
		Model:       "phi3",
		Prompt:      prompt,
		Temperature: 0.7,
		Stream:      false,
	}

	resp, err := llm.CallOllamaAPI(&reqBody)
	if err != nil {
		return nil, fmt.Errorf("error calling Ollama API: %v", err)
	}

	var baseResponse BaseResponse

	switch action {
	case ActionAppointment:	

		var v AppointmentResponse
		if err := json.Unmarshal([]byte(resp), &v); err != nil {
			// Try to extract JSON from the response
			start := strings.Index(resp, "{")
			end := strings.LastIndex(resp, "}")
			if start >= 0 && end > start {
				jsonStr := resp[start : end+1]
				if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
					fmt.Printf("Raw response: %s\n", resp)
					return nil, fmt.Errorf("error parsing schedule suggestion: %v", err)
				}
			} else {
				fmt.Printf("Raw response: %s\n", resp)
				return nil, fmt.Errorf("error parsing schedule suggestion: %v", err)
			}
		}
		baseResponse = BaseResponse{
			Action: v.Action,
			Parameters: map[string]interface{}{
				"pet_name": v.Parameters.PetName,
				"appointment_type": v.Parameters.AppointmentType,
				"date":     v.Parameters.Date,
				"time":     v.Parameters.Time,
			},
		}

	case ActionPetLog:
		// Clean the response by removing comments
		
		var v LogResponse
		if err := json.Unmarshal([]byte(resp), &v); err != nil {
			// Try to extract JSON from the response
			start := strings.Index(resp, "{")
			end := strings.LastIndex(resp, "}")
			if start >= 0 && end > start {
				jsonStr := resp[start : end+1]
				if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
					fmt.Printf("Raw response: %s\n", resp)
					return nil, fmt.Errorf("error parsing schedule suggestion: %v", err)
				}
			} else {
				fmt.Printf("Raw response: %s\n", resp)
				return nil, fmt.Errorf("error parsing schedule suggestion: %v", err)
			}
		}
		
		// If date is empty, use current date
		if v.Parameters.Date == "" {
			v.Parameters.Date = time.Now().Format("2006-01-02")
		}
		
		baseResponse = BaseResponse{
			Action: v.Action,
			Parameters: map[string]interface{}{
				"pet_name": v.Parameters.PetName,
				"activity": v.Parameters.Activity,
				"date":     v.Parameters.Date,
				"time":     v.Parameters.Time,
				"notes":    v.Parameters.Notes,
			},
		}
	case ActionPetSchedule:
		var v ScheduleResponse
		if err := json.Unmarshal([]byte(resp), &v); err != nil {
			// Try to extract JSON from the response
			start := strings.Index(resp, "{")
			end := strings.LastIndex(resp, "}")
			if start >= 0 && end > start {
				jsonStr := resp[start : end+1]
				if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
					fmt.Printf("Raw response: %s\n", resp)
					return nil, fmt.Errorf("error parsing schedule suggestion: %v", err)
				}
			} else {
				fmt.Printf("Raw response: %s\n", resp)
				return nil, fmt.Errorf("error parsing schedule suggestion: %v", err)
				}
			}
			baseResponse = BaseResponse{
				Action: v.Action,
				Parameters: map[string]interface{}{
					"pet_name": v.Parameters.PetName,
					"activity": v.Parameters.Activity,
					"date":     v.Parameters.Date,
					"time":     v.Parameters.Time,
					"notes":    v.Parameters.Notes,
				},
			}
		}

		// return response


		return &baseResponse, nil
	}

	func DetermineAction(ctx *gin.Context, description string) (*ActionResponse, error) {
	prompt := fmt.Sprintf(`As an AI assistant, your task is to interpret the user's request and determine the appropriate action to take. The possible actions are:

	- "appointment": For scheduling a one-time appointment or event.
	- "pet_log": For logging a one-time activity or event related to a pet.
	- "pet_schedule": For setting up a recurring schedule (e.g., daily, weekly).

	Task:
	Analyze the following user description and determine the most suitable action type.

	User Input:
	%s

	Response Format:
	Return only a valid JSON object, structured as follows:
	{
		"action": "<appointment | pet_log | pet_schedule>"
	}`, description)

	reqBody := llm.OllamaRequest{
		Model:       "phi3",
		Prompt:      prompt,
		Temperature: 0.1,
		Stream:      false,
	}

	resp, err := llm.CallOllamaAPI(&reqBody)
	if err != nil {
		return nil, fmt.Errorf("error calling Ollama API: %v", err)
	}

	var actionResponse ActionResponse
	if err := json.Unmarshal([]byte(resp), &actionResponse); err != nil {
		// Try to extract JSON from the response if it's wrapped in other text
		start := strings.Index(resp, "{")
		end := strings.LastIndex(resp, "}")
		if start >= 0 && end > start {
			jsonStr := resp[start : end+1]
			if err := json.Unmarshal([]byte(jsonStr), &actionResponse); err != nil {
				return nil, fmt.Errorf("error parsing action response: %v", err)
			}
		} else {
			return nil, fmt.Errorf("error parsing action response: %v", err)
		}
	}

	return &actionResponse, nil
}

func (s *PetScheduleService) ProcessSuggestion(ctx *gin.Context, description string) (*BaseResponse, error) {
	actionResponse, err := DetermineAction(ctx, description)
>>>>>>> ffc9071 (AI suggestion)
	if err != nil {
		return nil, fmt.Errorf("error determining action: %v", err)
	}

<<<<<<< HEAD
	var res *llm.BaseResponse

	switch actionResponse.Action {
	case llm.ActionAppointment:
		res, err = llm.GenerateSuggestionGemini(ctx, s.config, actionResponse.Action, description)
=======
	var res *BaseResponse

	switch actionResponse.Action {
	case ActionAppointment:
		res, err = GenerateSuggestion(ctx, actionResponse.Action, description)
>>>>>>> ffc9071 (AI suggestion)
		if err != nil {
			return nil, fmt.Errorf("error generating suggestion: %v", err)
		}
		return res, nil
<<<<<<< HEAD
	case llm.ActionPetLog:
		res, err = llm.GenerateSuggestionGemini(ctx, s.config, actionResponse.Action, description)
=======
	case ActionPetLog:
		res, err = GenerateSuggestion(ctx, actionResponse.Action, description)
>>>>>>> ffc9071 (AI suggestion)
		if err != nil {
			return nil, fmt.Errorf("error generating suggestion: %v", err)
		}
		return res, nil
<<<<<<< HEAD
	case llm.ActionPetSchedule:
		res, err = llm.GenerateSuggestionGemini(ctx, s.config, actionResponse.Action, description)
=======
	case ActionPetSchedule:
		res, err = GenerateSuggestion(ctx, actionResponse.Action, description)
>>>>>>> ffc9071 (AI suggestion)
		if err != nil {
			return nil, fmt.Errorf("error generating suggestion: %v", err)
		}
		return res, nil
	}
	return res, nil
}
<<<<<<< HEAD
=======
	// Use the provided request to create the schedule in the database
	// Return an error if the schedule creation fails

	return nil
}
>>>>>>> 272832d (redis cache)
=======
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
<<<<<<< HEAD
>>>>>>> e01abc5 (pet schedule api)
=======

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
<<<<<<< HEAD
>>>>>>> 6610455 (feat: redis queue)
=======

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
	err := s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		err := q.DeletePetSchedule(ctx, scheduleID)

		if err != nil {
			return fmt.Errorf("error deleting reminder: ", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting reminder: ", err)
	}
	return nil
}
>>>>>>> eb8d761 (updated pet schedule)
=======
>>>>>>> 4c66ef3 (feat: update schedule API)
=======

>>>>>>> ffc9071 (AI suggestion)
