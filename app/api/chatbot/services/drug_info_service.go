package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/models"
)

// DrugInfoService handles drug information lookup
type DrugInfoService struct {
	openFDAService *OpenFDAService
	geminiService  *GeminiService
	cache          *CacheService
}

// NewDrugInfoService creates a new instance of DrugInfoService
func NewDrugInfoService(openFDAService *OpenFDAService, geminiService *GeminiService, cache *CacheService) *DrugInfoService {
	return &DrugInfoService{
		openFDAService: openFDAService,
		geminiService:  geminiService,
		cache:          cache,
	}
}

// ExtractDrugInfo processes a user message to extract drug name and get information about it
func (s *DrugInfoService) ExtractDrugInfo(message string) (*models.DrugInfo, error) {
	// Use Gemini to extract drug name from the message
	drugName, err := s.extractDrugName(message)
	if err != nil {
		return nil, fmt.Errorf("could not extract drug name from message: %w", err)
	}

	if drugName == "" {
		return nil, fmt.Errorf("could not identify drug name from the message")
	}

	// Get drug information from OpenFDA
	drugInfo, err := s.getDrugInformation(drugName)
	if err != nil {
		return nil, fmt.Errorf("error getting drug information: %w", err)
	}

	return drugInfo, nil
}

// extractDrugName uses Gemini to analyze the user message and extract drug name
func (s *DrugInfoService) extractDrugName(message string) (string, error) {
	// Construct a prompt for Gemini
	prompt := fmt.Sprintf(`
                Extract the name of the medication or drug from this message: "%s"
                Return only JSON with this field:
                - drugName: the exact drug name mentioned

                Format: Valid JSON only, no explanation.
                `, message)

	// Send the prompt to Gemini
	response, err := s.geminiService.sendRequest(prompt)
	if err != nil {
		return "", err
	}

	// Extract the JSON from the response
	jsonStr := extractJSONFromText(response)
	if jsonStr == "" {
		// Fallback: Try to extract using keyword matching
		return s.extractDrugNameWithKeywords(message)
	}

	// Parse the response
	var result struct {
		DrugName string `json:"drugName"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return s.extractDrugNameWithKeywords(message)
	}

	return result.DrugName, nil
}

// extractDrugNameWithKeywords is a fallback method that extracts drug names using common patterns
func (s *DrugInfoService) extractDrugNameWithKeywords(message string) (string, error) {
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

	// Try to find a drug name in the message
	for _, drug := range commonDrugs {
		if strings.Contains(messageLower, drug) {
			return drug, nil
		}
	}

	// Look for common patterns
	patterns := []string{
		"about (\\w+)",
		"information on (\\w+)",
		"drug (\\w+)",
		"medication (\\w+)",
		"medicine (\\w+)",
		"take (\\w+)",
		"taking (\\w+)",
		"prescribed (\\w+)",
		"(\\w+) pill",
		"(\\w+) tablet",
		"(\\w+) capsule",
	}

	for _, patternStr := range patterns {
		pattern := regexp.MustCompile(patternStr)
		matches := pattern.FindStringSubmatch(messageLower)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	// If no drug name found, try to find the first word after "what is"
	whatIsPattern := regexp.MustCompile(`what is (\w+)`)
	matches := whatIsPattern.FindStringSubmatch(messageLower)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("could not identify drug name from message")
}

// getDrugInformation fetches comprehensive information about a drug from OpenFDA
func (s *DrugInfoService) getDrugInformation(drugName string) (*models.DrugInfo, error) {
	// Check cache first
	if s.cache != nil {
		cacheKey := fmt.Sprintf("drug_info:%s", strings.ToLower(drugName))
		var cachedDrugInfo models.DrugInfo
		if s.cache.GetJSON(cacheKey, &cachedDrugInfo) {
			// Found in cache
			return &cachedDrugInfo, nil
		}
	}

	// Not in cache, create new drug info
	drugInfo := &models.DrugInfo{
		DrugName:      drugName,
		SideEffects:   []string{},
		Recalls:       []string{},
		SourceDetails: "OpenFDA API (FDA.gov)",
	}

	// Get drug label information (indications, warnings, contraindications)
	err := s.getLabelInfo(drugInfo)
	if err != nil {
		// Continue even if there's an error, as we can still get other information
		// log.Printf("Error getting label info for %s: %v", drugName, err)
	}

	// Get side effect information
	err = s.getSideEffectInfo(drugInfo)
	if err != nil {
		// Continue even if there's an error
		// log.Printf("Error getting side effect info for %s: %v", drugName, err)
	}

	// Get recall information
	err = s.getRecallInfo(drugInfo)
	if err != nil {
		// Continue even if there's an error
		// log.Printf("Error getting recall info for %s: %v", drugName, err)
	}

	// If we have no data at all, use Gemini for a fallback response
	if drugInfo.Indications == "" && len(drugInfo.SideEffects) == 0 &&
		len(drugInfo.Recalls) == 0 && drugInfo.Warnings == "" &&
		drugInfo.Contraindications == "" {

		// Try to get information using Gemini
		err := s.getDrugInfoFromGemini(drugInfo)
		if err != nil {
			return nil, fmt.Errorf("no information found for %s in OpenFDA or via Gemini", drugName)
		}
	}

	// Cache the result for 24 hours (86400 seconds)
	if s.cache != nil {
		cacheKey := fmt.Sprintf("drug_info:%s", strings.ToLower(drugName))
		s.cache.Set(cacheKey, drugInfo, 86400)
	}

	return drugInfo, nil
}

// getLabelInfo gets information from drug labels
func (s *DrugInfoService) getLabelInfo(drugInfo *models.DrugInfo) error {
	params := url.Values{}

	// Search for the drug in labels
	query := fmt.Sprintf("openfda.generic_name:\"%s\" OR openfda.brand_name:\"%s\" OR openfda.substance_name:\"%s\"",
		drugInfo.DrugName, drugInfo.DrugName, drugInfo.DrugName)
	params.Add("search", query)
	params.Add("limit", "1")

	// Make the API request
	data, err := s.openFDAService.makeAPIRequest("/drug/label.json", params)
	if err != nil {
		return err
	}

	// Convert to map for processing
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format from OpenFDA")
	}

	// Check if any results were found
	results, ok := dataMap["results"].([]interface{})
	if !ok || len(results) == 0 {
		return fmt.Errorf("no label information found for %s", drugInfo.DrugName)
	}

	// Extract relevant information from the first result
	firstResult, ok := results[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid result format from OpenFDA")
	}

	// Extract indications
	if indications, ok := firstResult["indications_and_usage"].([]interface{}); ok && len(indications) > 0 {
		drugInfo.Indications = indications[0].(string)
	}

	// Extract warnings
	if warnings, ok := firstResult["warnings"].([]interface{}); ok && len(warnings) > 0 {
		drugInfo.Warnings = warnings[0].(string)
	} else if boxedWarnings, ok := firstResult["boxed_warnings"].([]interface{}); ok && len(boxedWarnings) > 0 {
		drugInfo.Warnings = "BOXED WARNING: " + boxedWarnings[0].(string)
	}

	// Extract contraindications
	if contraindications, ok := firstResult["contraindications"].([]interface{}); ok && len(contraindications) > 0 {
		drugInfo.Contraindications = contraindications[0].(string)
	}

	return nil
}

// getSideEffectInfo gets information about side effects
func (s *DrugInfoService) getSideEffectInfo(drugInfo *models.DrugInfo) error {
	params := url.Values{}

	// Search for this drug in adverse event reports
	query := fmt.Sprintf("patient.drug.medicinalproduct:\"%s\" OR patient.drug.openfda.generic_name:\"%s\" OR patient.drug.openfda.brand_name:\"%s\"",
		drugInfo.DrugName, drugInfo.DrugName, drugInfo.DrugName)
	params.Add("search", query)
	params.Add("count", "patient.reaction.reactionmeddrapt.exact")
	params.Add("limit", "10") // Get top 10 side effects

	// Make the API request
	data, err := s.openFDAService.makeAPIRequest("/drug/event.json", params)
	if err != nil {
		return err
	}

	// Convert to map for processing
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format from OpenFDA")
	}

	// Extract side effects from results
	results, ok := dataMap["results"].([]interface{})
	if !ok || len(results) == 0 {
		return fmt.Errorf("no side effect information found for %s", drugInfo.DrugName)
	}

	// Process each side effect
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		if term, ok := resultMap["term"].(string); ok {
			drugInfo.SideEffects = append(drugInfo.SideEffects, term)
		}
	}

	return nil
}

// getRecallInfo gets information about drug recalls
func (s *DrugInfoService) getRecallInfo(drugInfo *models.DrugInfo) error {
	params := url.Values{}

	// Search for recalls related to this drug
	query := fmt.Sprintf("openfda.generic_name:\"%s\" OR openfda.brand_name:\"%s\" OR product_description:\"%s\"",
		drugInfo.DrugName, drugInfo.DrugName, drugInfo.DrugName)
	params.Add("search", query)
	params.Add("limit", "5") // Limit to 5 most recent recalls

	// Make the API request
	data, err := s.openFDAService.makeAPIRequest("/drug/enforcement.json", params)
	if err != nil {
		return err
	}

	// Convert to map for processing
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid response format from OpenFDA")
	}

	// Extract recalls from results
	results, ok := dataMap["results"].([]interface{})
	if !ok || len(results) == 0 {
		// No recalls found, which is actually good news
		drugInfo.Recalls = append(drugInfo.Recalls, "No active recalls found for this medication.")
		return nil
	}

	// Process each recall
	for _, result := range results {
		resultMap, ok := result.(map[string]interface{})
		if !ok {
			continue
		}

		var recallInfo string

		// Get recall reason
		if reason, ok := resultMap["reason_for_recall"].(string); ok {
			recallInfo = "Reason: " + reason
		}

		// Get recall date
		if date, ok := resultMap["recall_initiation_date"].(string); ok {
			recallInfo += " (Date: " + date + ")"
		}

		// Get classification (severity)
		if classification, ok := resultMap["classification"].(string); ok {
			recallInfo += " - " + classification
		}

		drugInfo.Recalls = append(drugInfo.Recalls, recallInfo)
	}

	return nil
}

// getDrugInfoFromGemini uses Gemini to get information when OpenFDA doesn't have data
func (s *DrugInfoService) getDrugInfoFromGemini(drugInfo *models.DrugInfo) error {
	// Construct a prompt for Gemini
	prompt := fmt.Sprintf(`
Provide factual, accurate information about the medication/drug "%s". 
Return only JSON with these fields:
- indications: brief description of what the drug is used for
- sideEffects: array of common side effects (up to 5)
- warnings: important safety warnings
- contraindications: conditions when the drug should not be used

Format: Valid JSON only, no explanation. Do NOT include any disclaimers. Use only factual, medically accurate information.
`, drugInfo.DrugName)

	// Send the prompt to Gemini
	response, err := s.geminiService.sendRequest(prompt)
	if err != nil {
		return err
	}

	// Extract the JSON from the response
	jsonStr := extractJSONFromText(response)
	if jsonStr == "" {
		return fmt.Errorf("could not extract drug information from Gemini response")
	}

	// Parse the response
	var result struct {
		Indications       string   `json:"indications"`
		SideEffects       []string `json:"sideEffects"`
		Warnings          string   `json:"warnings"`
		Contraindications string   `json:"contraindications"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return err
	}

	// Update the drug info with Gemini data
	if drugInfo.Indications == "" {
		drugInfo.Indications = result.Indications
	}

	if len(drugInfo.SideEffects) == 0 {
		drugInfo.SideEffects = result.SideEffects
	}

	if drugInfo.Warnings == "" {
		drugInfo.Warnings = result.Warnings
	}

	if drugInfo.Contraindications == "" {
		drugInfo.Contraindications = result.Contraindications
	}

	// Add note about the data source
	drugInfo.SourceDetails += " and Google Gemini"

	return nil
}

// Sử dụng hàm extractJSONFromText từ gemini_service.go
