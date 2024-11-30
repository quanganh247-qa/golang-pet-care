package location

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type GoongServiceInterface interface {
<<<<<<< HEAD
	AutoCompleteService(input string, location *Location, limit string) (*AutocompleteResponse, error)
=======
	AutoCompleteService(input string, location *Location) (*AutocompleteResponse, error)
>>>>>>> 4625843 (added goong maps api)
	GetPlaceDetailService(placeID string) (*PlaceDetailResponse, error)
	GetDirectionService(req DirectionRequest) (*DirectionsResponse, error)
	ForwardGeocodeService(address string) (*GeocodeResponse, error)
	ReverseGeocodeService(latlng string) (*GeocodeResponse, error)
	// DistanceMatrixService(req DistanceMatrixRequest) (*DistanceMatrixResponse, error)
}

// Autocomplete searches for places based on input text
<<<<<<< HEAD
func (s *GoongService) AutoCompleteService(input string, location *Location, limit string) (*AutocompleteResponse, error) {
	// Build base URL
	baseURL := fmt.Sprintf("%s/Place/AutoComplete", s.config.BaseURL)
=======
func (s *GoongService) AutoCompleteService(input string, location *Location) (*AutocompleteResponse, error) {
	// Build base URL
	fmt.Println(s.config.BaseURL)

	baseURL := fmt.Sprintf("%s/Place/AutoComplete", s.config.BaseURL)
	fmt.Println(baseURL)
>>>>>>> 4625843 (added goong maps api)
	// Add parameters
	params := url.Values{}
	params.Add("api_key", s.config.APIKey)
	params.Add("input", input)
<<<<<<< HEAD
	params.Add("limit", limit)
=======
>>>>>>> 4625843 (added goong maps api)

	// Add location if provided
	if location != nil {
		params.Add("location", fmt.Sprintf("%f,%f", location.Lat, location.Lng))
	}

	// Build final URL
	requestURL := baseURL + "?" + params.Encode()

	// Make request
	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var result AutocompleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

// GetPlaceDetail gets detailed information about a place
func (s *GoongService) GetPlaceDetailService(placeID string) (*PlaceDetailResponse, error) {
	// Build URL
	url := fmt.Sprintf("%s/Place/Detail?place_id=%s&api_key=%s",
		s.config.BaseURL,
		placeID,
		s.config.APIKey,
	)

	// Make request
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var result PlaceDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

func (s *GoongService) GetDirectionService(req DirectionRequest) (*DirectionsResponse, error) {
	// Build URL
	url := fmt.Sprintf("%s/Direction?origin=%s&destination=%s&api_key=%s",
		s.config.BaseURL,
		req.Origin,
		req.Destination,
		s.config.APIKey,
	)
	// Make request
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	var result DirectionsResponse

	// Parse response using unmashal
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &result, nil
}

func (s *GoongService) ForwardGeocodeService(address string) (*GeocodeResponse, error) {
	// Build URL
	url := fmt.Sprintf("%s/Geocode?address=%s&api_key=%s", s.config.BaseURL, address, s.config.APIKey)

	// Make request
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var result GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &result, nil

}

func (s *GoongService) ReverseGeocodeService(latlng string) (*GeocodeResponse, error) {
	// Build URL
	url := fmt.Sprintf("%s/Geocode?latlng=%s&api_key=%s", s.config.BaseURL, latlng, s.config.APIKey)

	// Make request
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	// Parse response
	var result GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return &result, nil

}

// Distance Matrix
// func (s *GoongService) DistanceMatrixService(req DistanceMatrixRequest) (*DistanceMatrixResponse, error) {
// 	// Build URL
// 	url := fmt.Sprintf("%s/DistanceMatrix?origins=%s&destinations=%s&vehicle=%s&api_key=%s",
// 		s.config.BaseURL,
// 		req.Origins,
// 		req.Destinations,
// 		req.Vehicle,
// 		s.config.APIKey,
// 	)
// 	// Make request
// 	resp, err := s.client.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to make request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Check response status
// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		fmt.Println(string(body))
// 		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
// 	}

// 	// Parse response
// 	var result DistanceMatrixResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return nil, fmt.Errorf("failed to decode response: %v", err)
// 	}
// 	return &result, nil
// }
