package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/models"
)

// SideEffectService handles side effect reporting related functionality
type SideEffectService struct {
	openFDAService *OpenFDAService
	geminiService  *GeminiService
	cache          *CacheService
}

// NewSideEffectService creates a new instance of SideEffectService
func NewSideEffectService(openFDAService *OpenFDAService, geminiService *GeminiService) *SideEffectService {
	return &SideEffectService{
		openFDAService: openFDAService,
		geminiService:  geminiService,
		cache:          nil, // We'll use the cache from openFDAService which already implements it
	}
}

// ExtractSideEffectReport processes a user message to extract drug and side effect information
// and prepares a report on the reported side effect
func (s *SideEffectService) ExtractSideEffectReport(message string) (*models.SideEffectReport, error) {
	// Use Gemini to extract drug name and side effect from user message
	drugName, sideEffect, err := s.extractDrugAndSideEffect(message)
	if err != nil {
		return nil, fmt.Errorf("could not extract drug and side effect from message: %w", err)
	}

	if drugName == "" || sideEffect == "" {
		return nil, fmt.Errorf("could not identify both drug name and side effect from the message")
	}

	// Check if the side effect is common for this drug
	isCommon, frequency, err := s.checkSideEffectFrequency(drugName, sideEffect)
	if err != nil {
		// Even if there's an error checking frequency, we can still provide a report
		// So we'll log the error but continue
		// log.Printf("Error checking side effect frequency: %v", err)
	}

	// Create the side effect report
	report := &models.SideEffectReport{
		DrugName:       drugName,
		SideEffectName: sideEffect,
		IsCommon:       isCommon,
		Frequency:      frequency,
		ReportingSteps: s.getReportingSteps(),
		ReportingLink:  "https://www.fda.gov/safety/reporting-serious-problems-fda/how-consumers-can-report-adverse-event-or-serious-problem-fda",
	}

	return report, nil
}

// extractDrugAndSideEffect uses Gemini to analyze the user message and extract drug and side effect
func (s *SideEffectService) extractDrugAndSideEffect(message string) (string, string, error) {
	// Construct a prompt for Gemini
	prompt := fmt.Sprintf(`
Extract drug name and side effect from this message: "%s"
Return only JSON with these fields:
- drugName: the exact drug name mentioned
- sideEffectName: the specific side effect or adverse reaction mentioned

Format: Valid JSON only, no explanation.
`, message)

	// Send the prompt to Gemini
	response, err := s.geminiService.sendRequest(prompt)
	if err != nil {
		return "", "", err
	}

	// Extract the JSON from the response
	jsonStr := extractJSONFromText(response)
	if jsonStr == "" {
		// If no valid JSON is found, try to extract using regex
		return s.extractDrugAndSideEffectWithRegex(message)
	}

	// Parse the response
	var result struct {
		DrugName       string `json:"drugName"`
		SideEffectName string `json:"sideEffectName"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return s.extractDrugAndSideEffectWithRegex(message)
	}

	return result.DrugName, result.SideEffectName, nil
}

// extractDrugAndSideEffectWithRegex is a fallback method that uses regex patterns to extract drug and side effect
func (s *SideEffectService) extractDrugAndSideEffectWithRegex(message string) (string, string, error) {
	// Convert to lowercase for case-insensitive matching
	messageLower := strings.ToLower(message)

	// Common drug names to look for
	commonDrugs := []string{
		"acetaminophen", "paracetamol", "ibuprofen", "aspirin", "naproxen",
		"lisinopril", "atorvastatin", "lipitor", "metformin", "amlodipine",
		"gabapentin", "omeprazole", "levothyroxine", "simvastatin", "metoprolol",
		"losartan", "albuterol", "amoxicillin", "hydrochlorothiazide", "sertraline",
		"zoloft", "fluoxetine", "prozac", "montelukast", "singulair",
	}

	// Common side effects to look for
	commonSideEffects := []string{
		"headache", "nausea", "dizziness", "fatigue", "drowsiness",
		"insomnia", "vomiting", "diarrhea", "constipation", "rash",
		"itching", "swelling", "fever", "pain", "anxiety",
		"depression", "dry mouth", "blurred vision", "confusion", "weakness",
		"cough", "shortness of breath", "heart palpitations", "muscle pain", "joint pain",
	}

	// Try to find a drug name
	var drugName string
	for _, drug := range commonDrugs {
		if strings.Contains(messageLower, drug) {
			drugName = drug
			break
		}
	}

	// Try to find a side effect
	var sideEffect string
	for _, effect := range commonSideEffects {
		if strings.Contains(messageLower, effect) {
			sideEffect = effect
			break
		}
	}

	// Use regex to look for patterns like "after taking X" or "X causes Y"
	if drugName == "" {
		afterTakingPattern := regexp.MustCompile(`after taking (\w+)`)
		causesPattern := regexp.MustCompile(`(\w+) (causes|caused|is causing|gave me)`)

		matches := afterTakingPattern.FindStringSubmatch(messageLower)
		if len(matches) > 1 {
			drugName = matches[1]
		} else {
			matches = causesPattern.FindStringSubmatch(messageLower)
			if len(matches) > 1 {
				drugName = matches[1]
			}
		}
	}

	// If we still couldn't find both, return an error
	if drugName == "" || sideEffect == "" {
		return "", "", fmt.Errorf("could not identify drug name and side effect from message")
	}

	return drugName, sideEffect, nil
}

// checkSideEffectFrequency checks if a given side effect is commonly reported for a drug
func (s *SideEffectService) checkSideEffectFrequency(drugName, sideEffect string) (bool, string, error) {
	// Try to get from cache first (cache key combines drug name and side effect)
	cacheKey := fmt.Sprintf("side_effect_frequency:%s:%s", strings.ToLower(drugName), strings.ToLower(sideEffect))

	// If we have access to the cache through openFDAService
	if s.openFDAService != nil && s.openFDAService.cache != nil {
		var cachedResult struct {
			IsCommon  bool
			Frequency string
		}

		if s.openFDAService.cache.GetJSON(cacheKey, &cachedResult) {
			return cachedResult.IsCommon, cachedResult.Frequency, nil
		}
	}

	// Not found in cache, proceed with API calls
	// Create query parameters for OpenFDA
	params := url.Values{}

	// Search for this drug and side effect combination
	drugQuery := fmt.Sprintf("patient.drug.medicinalproduct:\"%s\" OR patient.drug.openfda.generic_name:\"%s\" OR patient.drug.openfda.brand_name:\"%s\"",
		drugName, drugName, drugName)

	effectQuery := fmt.Sprintf("patient.reaction.reactionmeddrapt:\"%s\"", sideEffect)

	// Combine queries
	searchQuery := fmt.Sprintf("(%s) AND (%s)", drugQuery, effectQuery)
	params.Add("search", searchQuery)
	params.Add("count", "patient.reaction.reactionmeddrapt.exact")

	// Make the API request
	data, err := s.openFDAService.makeAPIRequest("/drug/event.json", params)
	if err != nil {
		return false, "Unknown", err
	}

	// Convert to map for processing
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return false, "Unknown", fmt.Errorf("invalid response format from OpenFDA")
	}

	// Check if meta information exists
	_, metaOk := dataMap["meta"].(map[string]interface{})
	if !metaOk {
		return false, "Unknown", fmt.Errorf("meta information missing in OpenFDA response")
	}

	// Get the total number of reports for this drug-side effect pair
	var totalReports float64
	if results, ok := dataMap["results"].([]interface{}); ok && len(results) > 0 {
		for _, result := range results {
			resultMap, ok := result.(map[string]interface{})
			if !ok {
				continue
			}

			if term, ok := resultMap["term"].(string); ok {
				if strings.EqualFold(term, sideEffect) {
					if count, ok := resultMap["count"].(float64); ok {
						totalReports = count
						break
					}
				}
			}
		}
	}

	// Now get the total number of reports for this drug (to calculate frequency)
	paramsTotal := url.Values{}
	paramsTotal.Add("search", drugQuery)
	paramsTotal.Add("limit", "1")

	totalData, err := s.openFDAService.makeAPIRequest("/drug/event.json", paramsTotal)
	if err != nil {
		return false, "Unknown", err
	}

	totalDataMap, ok := totalData.(map[string]interface{})
	if !ok {
		return false, "Unknown", fmt.Errorf("invalid response format from OpenFDA")
	}

	totalMeta, ok := totalDataMap["meta"].(map[string]interface{})
	if !ok {
		return false, "Unknown", fmt.Errorf("meta information missing in OpenFDA response")
	}

	var drugTotalReports float64
	if results, ok := totalMeta["results"].(map[string]interface{}); ok {
		if total, ok := results["total"].(float64); ok {
			drugTotalReports = total
		}
	}

	// Determine commonality and frequency based on the data
	isCommon := false
	frequency := "Rare"

	if drugTotalReports > 0 {
		// Calculate frequency percentage
		percentage := (totalReports / drugTotalReports) * 100

		if percentage >= 5 {
			isCommon = true
			frequency = "Very Common (≥5%)"
		} else if percentage >= 1 {
			isCommon = true
			frequency = "Common (1-5%)"
		} else if percentage >= 0.1 {
			frequency = "Uncommon (0.1-1%)"
		} else if percentage >= 0.01 {
			frequency = "Rare (0.01-0.1%)"
		} else {
			frequency = "Very Rare (<0.01%)"
		}

		// If we have at least some reports but can't calculate percentage
		if totalReports > 100 && percentage == 0 {
			isCommon = true
			frequency = "Common (based on number of reports)"
		} else if totalReports > 10 && percentage == 0 {
			frequency = "Reported multiple times"
		}
	} else if totalReports > 0 {
		// If we have reports but couldn't get total drug reports
		if totalReports > 100 {
			isCommon = true
			frequency = "Likely Common (based on number of reports)"
		} else if totalReports > 10 {
			frequency = "Reported multiple times"
		} else {
			frequency = "Reported in a few cases"
		}
	} else {
		frequency = "Not commonly reported"
	}

	// Store result in cache for 1 week (604800 seconds) if we have access to the cache
	if s.openFDAService != nil && s.openFDAService.cache != nil {
		cachedResult := struct {
			IsCommon  bool
			Frequency string
		}{
			IsCommon:  isCommon,
			Frequency: frequency,
		}

		s.openFDAService.cache.Set(cacheKey, cachedResult, 604800)
	}

	return isCommon, frequency, nil
}

// getReportingSteps returns a formatted string with instructions for reporting side effects
func (s *SideEffectService) getReportingSteps() string {
	return `
1. Visit the FDA MedWatch website: https://www.fda.gov/safety/medwatch-fda-safety-information-and-adverse-event-reporting-program
2. Click on "Report a Problem" 
3. Choose either the online reporting form or downloadable form
4. Complete the form with your information, drug details, and description of the side effect
5. Submit the form to the FDA

You can also call FDA at 1-800-FDA-1088 to report by phone.
`
}

// Sử dụng hàm extractJSONFromText từ gemini_service.go
