package llm

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/quanganh247-qa/go-blog-be/app/util"
<<<<<<< HEAD
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/option"
)

// Giới hạn 15 request đồng thời (điều chỉnh dựa trên rate limit của Gemini)
var sem = semaphore.NewWeighted(15)

=======
	"google.golang.org/api/option"
)

>>>>>>> e859654 (Elastic search)
type SuggestionResponse interface {
	GetAction() string
	Validate() error
}

// BaseResponse contains common fields for all response types
type BaseResponse struct {
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
}

func (b BaseResponse) GetAction() string {
	return b.Action
}

type GeminiRequest struct {
	Prompt string `json:"prompt"`
}

const (
	ActionAppointment = "appointment"
	ActionPetLog      = "pet_log"
	ActionPetSchedule = "pet_schedule"
)

type AppointmentResponse struct {
	Action     string `json:"action"`
	Parameters struct {
		Title           string `json:"title"`
		PetName         string `json:"pet_name"`
		AppointmentType string `json:"appointment_type"`
		Date            string `json:"date"`
		Time            string `json:"time"`
	} `json:"parameters"`
}

type LogResponse struct {
	Action     string `json:"action"`
	Parameters struct {
		Title    string `json:"title"`
		PetName  string `json:"pet_name"`
		Activity string `json:"activity"`
		Date     string `json:"date"`
		Time     string `json:"time"`
		Notes    string `json:"notes"`
	} `json:"parameters"`
}

type ScheduleResponse struct {
	Action     string `json:"action"`
	Parameters struct {
		Title    string `json:"title"	`
		PetName  string `json:"pet_name"`
		Activity string `json:"activity"`
		Date     string `json:"date"`
		Time     string `json:"time"`
		Notes    string `json:"notes"`
	} `json:"parameters"`
}

type ActionResponse struct {
	Action string `json:"action"`
}

func DetermineActionGemini(ctx *gin.Context, config *util.Config, description string) (*ActionResponse, error) {
<<<<<<< HEAD
	// Acquire semaphore để giới hạn đồng thời
	if err := sem.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("server busy: %v", err)
	}
	defer sem.Release(1)

=======
>>>>>>> e859654 (Elastic search)
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

	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GoogleAPIKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

<<<<<<< HEAD
	// model := client.GenerativeModel("gemini-1.5-pro-latest")
	model := client.GenerativeModel("gemini-2.0-flash")
=======
	model := client.GenerativeModel("gemini-1.5-pro-latest")
>>>>>>> e859654 (Elastic search)

	model.ResponseMIMEType = "application/json"

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(prompt)},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}

	// Get the response text
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText = string(txt)

		}
	}

	var actionResponse ActionResponse
	if err := json.Unmarshal([]byte(responseText), &actionResponse); err != nil {
		// Try to extract JSON from the response if it's wrapped in other text
		start := strings.Index(responseText, "{")
		end := strings.LastIndex(responseText, "}")
		if start >= 0 && end > start {
			jsonStr := responseText[start : end+1]
			if err := json.Unmarshal([]byte(jsonStr), &actionResponse); err != nil {
				return nil, fmt.Errorf("error parsing action response: %v", err)
			}
		} else {
			return nil, fmt.Errorf("error parsing action response: %v", err)
		}
	}

	return &actionResponse, nil
}

func GenerateSuggestionGemini(ctx *gin.Context, config *util.Config, action, description string) (*BaseResponse, error) {
<<<<<<< HEAD

	// Acquire semaphore để giới hạn đồng thời
	if err := sem.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("server busy: %v", err)
	}
	defer sem.Release(1)

=======
>>>>>>> e859654 (Elastic search)
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
		  Generate a title in format: "{pet_name}'s {appointment_type} Appointment"
		- For "pet_log", extract "pet_name", "activity", "date" (in YYYY-MM-DD format), and "time" (in HH:MM format if mentioned).
		  Generate a title in format: "{pet_name}'s {activity} Log"
		- For "pet_schedule", extract "pet_name", "activity", "frequency" (e.g., "daily", "weekly"), and "start_date" (in YYYY-MM-DD format).
		  Generate a title in format: "{pet_name}'s {frequency} {activity}"
		- If the action is "unknown", provide a brief explanation in "parameters" under "reason".
		- Always use the current time %s for any dates mentioned without a year.
		- Always include a generated title in the parameters based on the rules above.

		Example:
		- User: "Log that Max had his walk today at 2:30 PM"
		- Response:
		{
			"action": "pet_log",
			"parameters": {
				"title": "Max's walk Log",
				"pet_name": "Max",
				"activity": "walk",
				"date": "%s",
				"time": "14:30"
			}
		}

		Analyze the following request and provide the JSON response:

		Request: %s with action: %s
	`, time.Now(), time.Now().Format("2006-01-02"), description, action)

	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GoogleAPIKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

<<<<<<< HEAD
	// model := client.GenerativeModel("gemini-1.5-pro-latest")
	model := client.GenerativeModel("gemini-2.0-flash")

=======
	model := client.GenerativeModel("gemini-1.5-pro-latest")
>>>>>>> e859654 (Elastic search)
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(prompt)},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, err
	}
	// Get the response text
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText = string(txt)

		}
	}
	var baseResponse BaseResponse

	switch action {
	case ActionAppointment:

		var v AppointmentResponse
		if err := json.Unmarshal([]byte(responseText), &v); err != nil {
			// Try to extract JSON from the response
			start := strings.Index(responseText, "{")
			end := strings.LastIndex(responseText, "}")
			if start >= 0 && end > start {
				jsonStr := responseText[start : end+1]
				if err := json.Unmarshal([]byte(jsonStr), &v); err != nil {
					fmt.Printf("Raw response: %s\n", responseText)
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
				"title":            v.Parameters.Title,
				"pet_name":         v.Parameters.PetName,
				"appointment_type": v.Parameters.AppointmentType,
				"date":             v.Parameters.Date,
				"time":             v.Parameters.Time,
			},
		}

	case ActionPetLog:
		var v LogResponse
		if err := json.Unmarshal([]byte(responseText), &v); err != nil {
			// Try to extract JSON from the response
			start := strings.Index(responseText, "{")
			end := strings.LastIndex(responseText, "}")
			if start >= 0 && end > start {
				jsonStr := responseText[start : end+1]
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
				"title":    v.Parameters.Title,
				"pet_name": v.Parameters.PetName,
				"activity": v.Parameters.Activity,
				"date":     v.Parameters.Date,
				"time":     v.Parameters.Time,
				"notes":    v.Parameters.Notes,
			},
		}
	case ActionPetSchedule:
		var v ScheduleResponse
		if err := json.Unmarshal([]byte(responseText), &v); err != nil {
			// Try to extract JSON from the response
			start := strings.Index(responseText, "{")
			end := strings.LastIndex(responseText, "}")
			if start >= 0 && end > start {
				jsonStr := responseText[start : end+1]
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
				"title":    v.Parameters.Title,
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
