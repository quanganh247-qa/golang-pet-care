package clinic_reporting

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/service/llm"
)

// ClinicReportingServiceInterface defines the interface for the analytics service.
type ClinicReportingServiceInterface interface {
	GetClinicAnalytics(ctx *gin.Context) (*ClinicAnalyticsResponse, error)
	GenerateSOAPNote(ctx *gin.Context, req GenerateSOAPNoteRequest) (*Note, error)
}

// ClinicReportingService implements the analytics collection functionality.
func (s *ClinicReportingService) GetClinicAnalytics(ctx *gin.Context) (*ClinicAnalyticsResponse, error) {
	// // Retrieve all appointments.
	// appointments, err := s.storeDB.ListAllAppointments(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to retrieve appointments: %w", err)
	// }

	// // Retrieve all pets.
	// pets, err := s.storeDB.ListPets(ctx, db.ListPetsParams{
	// 	Limit:  1000,
	// 	Offset: 0,
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to retrieve pets: %w", err)
	// }

	// // Retrieve all doctors.
	// doctors, err := s.storeDB.GetAllDoctors(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to retrieve doctors: %w", err)
	// }

	// totalAppointments := len(appointments)
	// upcomingAppointments := 0
	// pastAppointments := 0
	// now := time.Now()

	// // Assume each appointment has an AppointmentTime field of type time.Time.
	// for _, app := range appointments {
	// 	if app.Date.Time.After(now) {
	// 		upcomingAppointments++
	// 	} else {
	// 		pastAppointments++
	// 	}
	// }

	// return &ClinicAnalyticsResponse{
	// 	TotalAppointments:    totalAppointments,
	// 	UpcomingAppointments: upcomingAppointments,
	// 	PastAppointments:     pastAppointments,
	// 	TotalPets:            len(pets),
	// 	TotalDoctors:         len(doctors),
	// }, nil
	return nil, nil
}

func (s *ClinicReportingService) GenerateSOAPNote(ctx *gin.Context, req GenerateSOAPNoteRequest) (*Note, error) {
	prompt := fmt.Sprintf(`As a veterinary AI assistant, analyze the following consultation transcript and generate a SOAP note:

    Transcript:
    %s

    Generate a structured SOAP note with the following sections:
    - Subjective (patient info, chief complaint, duration, history, symptoms)
    - Objective (vital signs, examination findings)
    - Assessment (primary diagnosis, differential diagnoses)
    - Plan (immediate treatment, medications, follow-up, client education)

    Format the response ONLY as a valid JSON structure matching this example:
    {
        "subjective": {
            "patient_info": "species, age, sex",
            "chief_complaint": "main issue",
            "duration": "time period",
            "history": "relevant history",
            "symptoms": ["symptom1", "symptom2"]
			"behavior": {
				"activity_level": "description",
				"appetite": "description",
			}
        },
        "objective": {
            "vital_signs": {
                "temperature": "value",
				"heart_rate": "value",
				"respiratory_rate": "value",
				"weight": "value",
				"general_condition": "description",
            },
            "examination_findings": ["finding1", "finding2"]
        },
        "assessment": {
            "primary_diagnosis": "main diagnosis",
            "differentials": ["differential1", "differential2"]
        },
        "plan": {
            "immediate_treatment": ["treatment1", "treatment2"],
            "medications": ["medication1", "medication2"],
            "follow_up": "follow up plan",
            "client_education": ["education1", "education2"]
        }
    }`, req.Text)

	reqBody := llm.OllamaRequest{
		Model:       "mistral",
		Prompt:      prompt,
		Temperature: 0.7,
		Stream:      false,
	}

	resp, err := llm.CallOllamaAPI(&reqBody)
	if err != nil {
		return nil, fmt.Errorf("error calling Ollama API: %v", err)
	}

	// Parse the JSON response into Note structure
	var note Note
	if err := json.Unmarshal([]byte(resp), &note); err != nil {
		// Try to extract JSON from the response
		start := strings.Index(resp, "{")
		end := strings.LastIndex(resp, "}")
		if start >= 0 && end > start {
			jsonStr := resp[start : end+1]
			if err := json.Unmarshal([]byte(jsonStr), &note); err != nil {
				fmt.Printf("Raw response: %s\n", resp)
				return nil, fmt.Errorf("error parsing SOAP note: %v", err)
			}
		} else {
			fmt.Printf("Raw response: %s\n", resp)
			return nil, fmt.Errorf("error parsing SOAP note: %v", err)
		}
	}

	return &note, nil
}

func (s *ClinicReportingService) CreateSoapNote(ctx *gin.Context, req CreateSoapNoteRequest) (*SoapNoteResponse, error) {
	
	note,err := json.Marshal(req.Note)
	if err != nil {
		return nil, err
	}
	
	var soapNote db.SoapNote

	err = s.storeDB.ExecWithTransaction(ctx, func(q *db.Queries) error {
		soapNote, err = q.CreateSoapNote(ctx, db.CreateSoapNoteParams{
			AppointmentID: int32(req.AppointmentID),
			Note:          note,
		})
		if err != nil {
			return err
		}
		return nil
	})

	return &SoapNoteResponse{
		ID:            int(soapNote.ID.Bytes[0]),
		AppointmentID: int(soapNote.AppointmentID),
		Note:          string(soapNote.Note),
	}, nil
}
