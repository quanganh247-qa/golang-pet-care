package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/models"
)

// OpenFDAService handles interactions with the openFDA API
type OpenFDAService struct {
	apiBaseURL string
	apiKey     string
	client     *http.Client
	cache      *CacheService
}

// NewOpenFDAService creates a new instance of OpenFDAService
func NewOpenFDAService(apiKey string, cache *CacheService) *OpenFDAService {
	return &OpenFDAService{
		apiBaseURL: "https://api.fda.gov",
		apiKey:     apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: cache,
	}
}

// FetchDataBasedOnIntent fetches data from OpenFDA based on the intent and parameters
func (s *OpenFDAService) FetchDataBasedOnIntent(geminiResponse *models.GeminiQueryResponse) (interface{}, string, error) {
	// Map the query intent to the appropriate OpenFDA endpoint
	endpoint, queryType, err := s.mapIntentToEndpoint(geminiResponse.QueryIntent)
	if err != nil {
		return nil, "", err
	}

	// Build the query parameters based on the intent and search terms
	queryParams, err := s.buildQueryParams(geminiResponse)
	if err != nil {
		return nil, "", err
	}

	// Make the API request
	data, err := s.makeAPIRequest(endpoint, queryParams)
	if err != nil {
		return nil, "", err
	}

	return data, queryType, nil
}

// mapIntentToEndpoint maps a query intent to the appropriate OpenFDA API endpoint
func (s *OpenFDAService) mapIntentToEndpoint(intent string) (string, string, error) {
	intent = strings.ToLower(intent)

	if strings.Contains(intent, "drug") && strings.Contains(intent, "event") {
		return "/drug/event.json", "drugEvent", nil
	} else if strings.Contains(intent, "drug") && strings.Contains(intent, "label") {
		return "/drug/label.json", "drugLabel", nil
	} else if strings.Contains(intent, "device") && strings.Contains(intent, "event") {
		return "/device/event.json", "deviceEvent", nil
	} else if strings.Contains(intent, "food") && strings.Contains(intent, "event") {
		return "/food/event.json", "foodEvent", nil
	} else if strings.Contains(intent, "drug") && strings.Contains(intent, "enforcement") {
		return "/drug/enforcement.json", "drugEnforcement", nil
	} else if strings.Contains(intent, "device") && strings.Contains(intent, "enforcement") {
		return "/device/enforcement.json", "deviceEnforcement", nil
	} else if strings.Contains(intent, "food") && strings.Contains(intent, "enforcement") {
		return "/food/enforcement.json", "foodEnforcement", nil
	} else if strings.Contains(intent, "drug") {
		return "/drug/event.json", "drugEvent", nil
	} else if strings.Contains(intent, "device") {
		return "/device/event.json", "deviceEvent", nil
	} else if strings.Contains(intent, "food") {
		return "/food/event.json", "foodEvent", nil
	}

	// Default to drug events if intent is unclear
	return "/drug/event.json", "drugEvent", nil
}

// buildQueryParams constructs the query parameters for the OpenFDA API request
func (s *OpenFDAService) buildQueryParams(response *models.GeminiQueryResponse) (url.Values, error) {
	params := url.Values{}

	// Process search terms
	searchTerms := getStringValue(response.SearchTerms)
	if searchTerms != "" {
		// Construct a search query based on the intent and search terms
		searchQuery := s.constructSearchQuery(response.QueryIntent, searchTerms)
		if searchQuery != "" {
			params.Add("search", searchQuery)
		}
	}

	// Process time frame
	timeFrame := getStringValue(response.TimeFrame)
	if timeFrame != "" {
		timeQuery, err := s.parseTimeFrame(timeFrame)
		if err != nil {
			return nil, err
		}
		if timeQuery != "" {
			// If search already exists, append with AND
			if search := params.Get("search"); search != "" {
				params.Set("search", fmt.Sprintf("(%s) AND %s", search, timeQuery))
			} else {
				params.Add("search", timeQuery)
			}
		}
	}

	// Process demographic filters
	demographics := getStringValue(response.DemographicFilters)
	if demographics != "" {
		demoQuery := s.parseDemographics(demographics, response.QueryIntent)
		if demoQuery != "" {
			// If search already exists, append with AND
			if search := params.Get("search"); search != "" {
				params.Set("search", fmt.Sprintf("(%s) AND %s", search, demoQuery))
			} else {
				params.Add("search", demoQuery)
			}
		}
	}

	// Process data aggregation
	aggregation := getStringValue(response.DataAggregation)
	if aggregation != "" {
		countParam := s.parseAggregation(aggregation, response.QueryIntent)
		if countParam != "" {
			params.Add("count", countParam)
		}
	}

	// Set default limit if not aggregating
	if params.Get("count") == "" {
		params.Add("limit", "100")
	}

	return params, nil
}

// constructSearchQuery builds a search query string based on intent and search terms
func (s *OpenFDAService) constructSearchQuery(intent, terms string) string {
	intent = strings.ToLower(intent)
	terms = strings.ToLower(terms)

	var fieldPrefix string
	if strings.Contains(intent, "drug") && strings.Contains(intent, "event") {
		fieldPrefix = "patient.drug.medicinalproduct:"
		if strings.Contains(terms, "side effect") || strings.Contains(terms, "adverse") {
			fieldPrefix = "patient.reaction.reactionmeddrapt:"
			terms = strings.ReplaceAll(terms, "side effect", "")
			terms = strings.ReplaceAll(terms, "adverse", "")
		}
	} else if strings.Contains(intent, "drug") && strings.Contains(intent, "label") {
		fieldPrefix = "openfda.generic_name:"
		if strings.Contains(terms, "warning") {
			fieldPrefix = "warnings:"
			terms = strings.ReplaceAll(terms, "warning", "")
		}
	} else if strings.Contains(intent, "device") && strings.Contains(intent, "event") {
		fieldPrefix = "device.generic_name:"
	} else if strings.Contains(intent, "food") && strings.Contains(intent, "event") {
		fieldPrefix = "products.name_brand:"
	} else {
		fieldPrefix = ""
	}

	// Clean up search terms
	terms = strings.TrimSpace(terms)
	if terms == "" {
		return ""
	}

	// If multiple terms, wrap them in quotes
	if strings.Contains(terms, " ") {
		return fmt.Sprintf("%s\"%s\"", fieldPrefix, terms)
	}
	return fmt.Sprintf("%s%s", fieldPrefix, terms)
}

// parseTimeFrame converts a time frame description to an API query parameter
func (s *OpenFDAService) parseTimeFrame(timeFrame string) (string, error) {
	timeFrame = strings.ToLower(timeFrame)

	// Handle "last X years"
	if strings.Contains(timeFrame, "last") && strings.Contains(timeFrame, "year") {
		var years int
		_, err := fmt.Sscanf(timeFrame, "last %d years", &years)
		if err != nil {
			// Try alternative format
			_, err = fmt.Sscanf(timeFrame, "last %d year", &years)
		}
		if err == nil && years > 0 {
			startDate := time.Now().AddDate(-years, 0, 0).Format("20060102")
			return fmt.Sprintf("receivedate:[%s TO 99991231]", startDate), nil
		}
	}

	// Handle specific years range (2020-2022)
	if strings.Contains(timeFrame, "-") {
		var startYear, endYear int
		_, err := fmt.Sscanf(timeFrame, "%d-%d", &startYear, &endYear)
		if err == nil && startYear > 1900 && endYear > startYear {
			return fmt.Sprintf("receivedate:[%d0101 TO %d1231]", startYear, endYear), nil
		}
	}

	// Handle specific year
	var year int
	_, err := fmt.Sscanf(timeFrame, "%d", &year)
	if err == nil && year > 1900 {
		return fmt.Sprintf("receivedate:[%d0101 TO %d1231]", year, year), nil
	}

	// Default: no time frame filter
	return "", nil
}

// parseDemographics converts demographic information to API query parameters
func (s *OpenFDAService) parseDemographics(demographics, intent string) string {
	demographics = strings.ToLower(demographics)

	var queries []string

	// Handle age ranges
	if strings.Contains(demographics, "age") {
		if strings.Contains(demographics, "children") || strings.Contains(demographics, "pediatric") {
			queries = append(queries, "patient.patientonsetage:[0 TO 17]")
		} else if strings.Contains(demographics, "adult") {
			queries = append(queries, "patient.patientonsetage:[18 TO 64]")
		} else if strings.Contains(demographics, "elderly") || strings.Contains(demographics, "senior") {
			queries = append(queries, "patient.patientonsetage:[65 TO 120]")
		}
	}

	// Handle gender
	if strings.Contains(demographics, "male") && !strings.Contains(demographics, "female") {
		queries = append(queries, "patient.patientsex:1")
	} else if strings.Contains(demographics, "female") && !strings.Contains(demographics, "male") {
		queries = append(queries, "patient.patientsex:2")
	}

	// Combine all queries with AND
	if len(queries) > 0 {
		return strings.Join(queries, " AND ")
	}
	return ""
}

// parseAggregation creates a count parameter based on aggregation request
func (s *OpenFDAService) parseAggregation(aggregation, intent string) string {
	aggregation = strings.ToLower(aggregation)
	intent = strings.ToLower(intent)

	if strings.Contains(aggregation, "year") {
		return "receivedate"
	} else if strings.Contains(aggregation, "drug") && strings.Contains(intent, "event") {
		return "patient.drug.medicinalproduct.exact"
	} else if strings.Contains(aggregation, "reaction") || strings.Contains(aggregation, "side effect") {
		return "patient.reaction.reactionmeddrapt.exact"
	} else if strings.Contains(aggregation, "manufacturer") {
		return "companynumb.exact"
	} else if strings.Contains(aggregation, "seriousness") {
		return "serious"
	} else if strings.Contains(aggregation, "country") {
		return "occurcountry.exact"
	}

	return ""
}

// makeAPIRequest sends a request to the OpenFDA API and returns the response
func (s *OpenFDAService) makeAPIRequest(endpoint string, params url.Values) (interface{}, error) {
	// Add API key to parameters if available
	if s.apiKey != "" {
		params.Add("api_key", s.apiKey)
	}

	// Build the full URL
	fullURL := fmt.Sprintf("%s%s?%s", s.apiBaseURL, endpoint, params.Encode())

	// Create cache key from the URL
	cacheKey := fmt.Sprintf("openfda:%s", fullURL)

	// Try to get data from cache first
	if s.cache != nil {
		var cachedResult interface{}
		if s.cache.GetJSON(cacheKey, &cachedResult) {
			return cachedResult, nil
		}
	}

	// Create and send the request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openFDA API returned status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Store in cache for 1 hour (3600 seconds)
	if s.cache != nil {
		s.cache.Set(cacheKey, result, 3600)
	}

	return result, nil
}

// getStringValue extracts a string value from different types of data
func getStringValue(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			// Convert the first array element to string if possible
			if str, ok := v[0].(string); ok {
				return str
			}
			// Try converting to string via fmt.Sprint
			return fmt.Sprint(v[0])
		}
	case []string:
		if len(v) > 0 {
			return v[0]
		}
	case map[string]interface{}:
		// Try to marshal map to JSON string
		if bytes, err := json.Marshal(v); err == nil {
			return string(bytes)
		}
	default:
		// For all other types, try simple string conversion
		return fmt.Sprint(v)
	}

	return ""
}
