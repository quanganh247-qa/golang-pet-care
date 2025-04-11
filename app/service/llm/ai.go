package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/quanganh247-qa/go-blog-be/app/models"
	"github.com/sashabaranov/go-openai"
)

// AIService handles interactions with the OpenAI API
type AIService struct {
	client *openai.Client
}

// NewAIService creates a new AI service
func NewAIService(apiKey string) (*AIService, error) {
	if apiKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}

	client := openai.NewClient(apiKey)
	return &AIService{
		client: client,
	}, nil
}

// IsInitialized checks if the OpenAI client is properly initialized
func (s *AIService) IsInitialized() bool {
	return s.client != nil
}

// AnalyzeSymptoms analyzes symptoms and generates treatment recommendations
func (s *AIService) AnalyzeSymptoms(patient models.Patient, symptoms string) (models.AITreatmentPlan, error) {
	// Check if client is initialized
	if s.client == nil {
		return models.AITreatmentPlan{}, errors.New("OpenAI client not initialized, symptom analysis unavailable")
	}

	// Format the prompt for OpenAI
	prompt := s.formatSymptomAnalysisPrompt(patient, symptoms)

	// Call OpenAI API
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo, // Using GPT-4 since the client may not support gpt-4o yet
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a veterinary AI assistant that helps analyze symptoms and suggest treatment plans for pets. Provide accurate, medically sound recommendations based on the provided information.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.3,
			MaxTokens:   1500,
		},
	)

	if err != nil {
		return models.AITreatmentPlan{}, fmt.Errorf("OpenAI API error: %v", err)
	}

	if len(resp.Choices) == 0 {
		return models.AITreatmentPlan{}, errors.New("no response from OpenAI API")
	}

	// Parse the response to extract the treatment plan
	treatmentPlan, err := s.parseTreatmentResponse(resp.Choices[0].Message.Content)
	if err != nil {
		return models.AITreatmentPlan{}, fmt.Errorf("failed to parse treatment plan: %v", err)
	}

	return treatmentPlan, nil
}

// formatSymptomAnalysisPrompt creates a prompt for the OpenAI API based on patient data and symptoms
func (s *AIService) formatSymptomAnalysisPrompt(patient models.Patient, symptoms string) string {
	prompt := fmt.Sprintf(`
Patient Information:
- Name: %s
- Species: %s
- Breed: %s
- Age: %.1f years
- Gender: %s
- Weight: %.1f kg
- Current Status: %s

Medical Notes: %s

Presenting Symptoms:
%s

Based on the above information, please provide:
1. A detailed diagnosis based on the symptoms and patient information
2. A complete treatment plan with distinct phases
3. Each phase should have a name, description, duration (in days), and specific medications

Format your response as a JSON object with the following structure:
{
  "diagnosis": "Your detailed diagnosis here",
  "confidence": 0.8,
  "treatment_phases": [
    {
      "name": "Initial Treatment",
      "description": "Detailed description of this phase",
      "duration": 7,
      "order": 1,
      "medicines": [
        {
          "medicine_id": 0,
          "dosage": "Detailed dosage instructions",
          "frequency": "Frequency of administration",
          "duration": 7,
          "notes": "Any additional notes about this medication"
        }
      ]
    },
    {
      "name": "Follow-up Phase",
      "description": "Description of follow-up treatment",
      "duration": 14,
      "order": 2,
      "medicines": [
        {
          "medicine_id": 0,
          "dosage": "Detailed dosage instructions",
          "frequency": "Frequency of administration",
          "duration": 14,
          "notes": "Any additional notes"
        }
      ]
    }
  ]
}

IMPORTANT: Be very specific about dosages, which should be based on the animal's weight and condition. Always include at least two treatment phases. The first should be for immediate treatment and the second for follow-up care.
`,
		patient.Name,
		patient.Species,
		patient.Breed,
		patient.Age,
		patient.Gender,
		patient.Weight,
		patient.Status,
		patient.Notes,
		symptoms,
	)

	return prompt
}

// parseTreatmentResponse extracts structured treatment data from the API response
func (s *AIService) parseTreatmentResponse(response string) (models.AITreatmentPlan, error) {
	// Find JSON content in the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonEnd <= jsonStart {
		// If proper JSON is not found, try to parse the text into a structured format
		return s.extractTreatmentFromText(response)
	}

	jsonContent := response[jsonStart : jsonEnd+1]

	// First parse into a map to handle flexible types
	var rawMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonContent), &rawMap)
	if err != nil {
		fmt.Println("JSON parse error (map):", err)
		return models.AITreatmentPlan{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Create our treatment plan
	treatmentPlan := models.AITreatmentPlan{
		Diagnosis:       getStringValue(rawMap, "diagnosis", "No diagnosis provided"),
		Confidence:      getFloatValue(rawMap, "confidence", 0.7),
		TreatmentPhases: []models.TreatmentPhase{},
		CreatedAt:       time.Now().Format(time.RFC3339),
		UpdatedAt:       time.Now().Format(time.RFC3339),
	}

	// Process treatment phases if present
	if phasesRaw, ok := rawMap["treatment_phases"]; ok {
		phases, ok := phasesRaw.([]interface{})
		if ok {
			for i, phaseRaw := range phases {
				phase, ok := phaseRaw.(map[string]interface{})
				if !ok {
					continue
				}

				newPhase := models.TreatmentPhase{
					Name:        getStringValue(phase, "name", fmt.Sprintf("Phase %d", i+1)),
					Description: getStringValue(phase, "description", "No description provided"),
					Duration:    getIntValue(phase, "duration", 7),
					Order:       getIntValue(phase, "order", i+1),
					Medicines:   []models.PhaseMedicine{},
				}

				// Process medicines
				if medsRaw, ok := phase["medicines"]; ok {
					meds, ok := medsRaw.([]interface{})
					if ok {
						for _, medRaw := range meds {
							med, ok := medRaw.(map[string]interface{})
							if !ok {
								continue
							}

							newMedicine := models.PhaseMedicine{
								MedicineID: int64(getIntValue(med, "medicine_id", 0)),
								Dosage:     getStringValue(med, "dosage", "As prescribed"),
								Frequency:  getStringValue(med, "frequency", "As directed"),
								Duration:   getIntValue(med, "duration", newPhase.Duration),
								Notes:      getStringValue(med, "notes", ""),
							}
							newPhase.Medicines = append(newPhase.Medicines, newMedicine)
						}
					}
				}

				treatmentPlan.TreatmentPhases = append(treatmentPlan.TreatmentPhases, newPhase)
			}
		}
	}

	// If we still have no treatment phases, create a default one
	if len(treatmentPlan.TreatmentPhases) == 0 {
		defaultPhase := models.TreatmentPhase{
			Name:        "General Treatment",
			Description: "Based on the diagnosis: " + treatmentPlan.Diagnosis,
			Duration:    14,
			Order:       1,
			Medicines:   []models.PhaseMedicine{},
		}

		defaultMedicine := models.PhaseMedicine{
			Dosage:    "As prescribed by veterinarian",
			Frequency: "As directed",
			Duration:  14,
			Notes:     "Please consult with your veterinarian for specific medication details",
		}
		defaultPhase.Medicines = append(defaultPhase.Medicines, defaultMedicine)

		treatmentPlan.TreatmentPhases = append(treatmentPlan.TreatmentPhases, defaultPhase)
	}

	return treatmentPlan, nil
}

// Helper functions to extract values with type conversion
func getStringValue(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		case float64:
			return fmt.Sprintf("%.f", v)
		case int:
			return fmt.Sprintf("%d", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return defaultValue
}

func getIntValue(m map[string]interface{}, key string, defaultValue int) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case string:
			// Try to parse the string as an integer
			// Remove any non-numeric characters first
			numStr := regexp.MustCompile(`\d+`).FindString(v)
			if numStr != "" {
				if num, err := strconv.Atoi(numStr); err == nil {
					return num
				}
			}
		}
	}
	return defaultValue
}

func getFloatValue(m map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			// Try to parse the string as a float
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return defaultValue
}

// extractTreatmentFromText attempts to extract treatment information from non-JSON text
func (s *AIService) extractTreatmentFromText(text string) (models.AITreatmentPlan, error) {
	treatmentPlan := models.AITreatmentPlan{
		TreatmentPhases: []models.TreatmentPhase{},
		Confidence:      0.7, // Default confidence when parsing from text
	}

	// Look for diagnosis section
	if diagnosisStart := strings.Index(strings.ToLower(text), "diagnosis"); diagnosisStart != -1 {
		diagnosisText := text[diagnosisStart:]
		if nextSection := strings.Index(diagnosisText[10:], "\n\n"); nextSection != -1 {
			treatmentPlan.Diagnosis = strings.TrimSpace(diagnosisText[10 : 10+nextSection])
		} else {
			treatmentPlan.Diagnosis = "Diagnosis could not be extracted"
		}
	} else {
		treatmentPlan.Diagnosis = "Unknown diagnosis"
	}

	// Create a default treatment phase
	defaultPhase := models.TreatmentPhase{
		Name:        "Initial Treatment",
		Description: "Treatment based on initial assessment",
		Duration:    14, // Default 2 weeks
		Order:       1,
		Medicines:   []models.PhaseMedicine{},
	}

	// Look for treatment section and add to description
	if treatmentStart := strings.Index(strings.ToLower(text), "treatment"); treatmentStart != -1 {
		treatmentText := text[treatmentStart:]
		if nextSection := strings.Index(treatmentText[10:], "\n\n"); nextSection != -1 {
			defaultPhase.Description = strings.TrimSpace(treatmentText[10 : 10+nextSection])
		} else {
			defaultPhase.Description = strings.TrimSpace(treatmentText[10:])
		}
	} else {
		defaultPhase.Description = "No specific treatment plan generated"
	}

	// Look for medications and add as PhaseMedicine objects
	if medStart := strings.Index(strings.ToLower(text), "medication"); medStart != -1 {
		medText := text[medStart:]
		lines := strings.Split(medText, "\n")
		for i := 1; i < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if line == "" || strings.Contains(strings.ToLower(line), "follow") {
				break
			}
			// Remove bullet points or numbers
			line = strings.TrimLeft(line, "- *•1234567890. ")
			if line != "" {
				// Create a medicine entry from the text line
				medicine := models.PhaseMedicine{
					Dosage:    line,
					Frequency: "As directed",
					Duration:  defaultPhase.Duration,
					Notes:     "Extracted from AI response",
				}
				defaultPhase.Medicines = append(defaultPhase.Medicines, medicine)
			}
		}
	}

	// Add the phase to the treatment plan
	treatmentPlan.TreatmentPhases = append(treatmentPlan.TreatmentPhases, defaultPhase)

	// Look for follow-up timeframe and update the phase duration if found
	if followUpStart := strings.Index(strings.ToLower(text), "follow"); followUpStart != -1 {
		followUpText := text[followUpStart:]
		// Try to extract a number
		for _, word := range strings.Fields(followUpText) {
			var days int
			if _, err := fmt.Sscanf(word, "%d", &days); err == nil && days > 0 && days < 365 {
				if len(treatmentPlan.TreatmentPhases) > 0 {
					treatmentPlan.TreatmentPhases[0].Duration = days
					// Update medicines duration as well
					for i := range treatmentPlan.TreatmentPhases[0].Medicines {
						treatmentPlan.TreatmentPhases[0].Medicines[i].Duration = days
					}
				}
				break
			}
		}
	}

	return treatmentPlan, nil
}
