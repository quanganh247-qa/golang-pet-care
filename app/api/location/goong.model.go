package location

import (
	"net/http"
)

// GoongConfig contains configuration for Goong Maps API
type GoongConfig struct {
	APIKey  string
	BaseURL string
}

// GoongService handles interactions with Goong Maps API
type GoongService struct {
	config *GoongConfig
	client *http.Client
}
type GoongApi struct {
	controller GoongControllerInterface
}

type GoongController struct {
	service GoongServiceInterface
}

// alternatives	- > Boolean, if true, Directions service may return several routes in the response. If false, the service will return the best route.
type DirectionRequest struct {
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
	Vehicle      string `json:"vehicle"`
	Alternatives bool   `json:"alternatives,omitempty"` // Optional parameter for multiple routes
}

// DirectionsResponse represents the response from Goong Directions API
type DirectionsResponse struct {
	GeocodedWaypoints []Point `json:"geocoded_waypoints"`
	Routes            []Route `json:"routes"`
}

type Route struct {
	Bounds           Bounds     `json:"bounds"`
	Legs             []RouteLeg `json:"legs"`
	OverviewPolyline struct {
		Points string `json:"points"`
	} `json:"overview_polyline"`
	Warnings      []string `json:"warnings"`
	WaypointOrder []int    `json:"waypoint_order"`
}

type Bounds struct {
	Northeast Point `json:"northeast"`
	Southwest Point `json:"southwest"`
}

type RouteLeg struct {
	Distance Distance    `json:"distance"`
	Duration Duration    `json:"duration"`
	Steps    []RouteStep `json:"steps"`
}

type Distance struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type Duration struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type RouteStep struct {
	Distance    Distance `json:"distance"`
	Duration    Duration `json:"duration"`
	Name        string   `json:"name"`
	Instruction string   `json:"instruction"`
	Geometry    Geometry `json:"geometry"`
}

type Point struct {
	Location []float64 `json:"location"`
	Name     string    `json:"name"`
}

// Place represents a location returned by the Places API// AutocompleteResponse represents the response from Places API
type AutocompleteResponse struct {
	Predictions []Prediction `json:"predictions"`
	Status      string       `json:"status"`
}

type Prediction struct {
	Description          string             `json:"description"`
	MatchedSubstrings    []MatchedSubstring `json:"matched_substrings"`
	PlaceID              string             `json:"place_id"`
	Reference            string             `json:"reference"`
	StructuredFormatting struct {
		MainText      string `json:"main_text"`
		SecondaryText string `json:"secondary_text"`
	} `json:"structured_formatting"`
	Terms []Term   `json:"terms"`
	Types []string `json:"types"`
}

type MatchedSubstring struct {
	Length int `json:"length"`
	Offset int `json:"offset"`
}

type Term struct {
	Offset int    `json:"offset"`
	Value  string `json:"value"`
}

// PlaceDetailResponse represents the detail of a place
type PlaceDetailResponse struct {
	Result PlaceDetail `json:"result"`
	Status string      `json:"status"`
}

type PlaceDetail struct {
	PlaceID    string             `json:"place_id"`
	Name       string             `json:"name"`
	Address    string             `json:"formatted_address"`
<<<<<<< HEAD
	Geometry   Geometry           `json:"geometry"`
=======
	Location   Location           `json:"geometry"`
>>>>>>> 4625843 (added goong maps api)
	Types      []string           `json:"types"`
	Components []AddressComponent `json:"address_components"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

// geocode struct

type Geometry struct {
	Location Location `json:"location"`
}

type GeocodeResponse struct {
	Results []Component `json:"results"`
}

type Component struct {
	AddressComponents []AddressComponent `json:"address_components"`
	FormarttedAddress string             `json:"formatted_address"`
	Geometry          Geometry           `json:"geometry"`
	PlaceID           string             `json:"place_id"`
	Reference         string             `json:"reference"`
	PlusCode          PlusCode           `json:"plus_code"`
	Types             []string           `json:"types"`
}

type PlusCode struct {
	CompoundCode string `json:"compound_code"`
	GlobalCode   string `json:"global_code"`
}

// DistanceMatrixResponse represents the response from the Distance Matrix API

type DistanceMatrixResponse struct {
	Rows []Row `json:"rows"`
}

type Row struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Status   string   `json:"status"`
	Duration Duration `json:"duration"`
	Distance Distance `json:"distance"`
}

// DistanceMatrixRequest represents request parameters for Distance Matrix API
type DistanceMatrixRequest struct {
	Origins      []string `json:"origins"`      // Array of origin points "lat,lng"
	Destinations []string `json:"destinations"` // Array of destination points "lat,lng"
	Vehicle      string   `json:"vehicle"`      // Type of vehicle (car, bike, taxi, truck)
}
