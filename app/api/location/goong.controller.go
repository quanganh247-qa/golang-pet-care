package location

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type GoongControllerInterface interface {
	Autocomplete(c *gin.Context)
	GetPlaceDetail(c *gin.Context)
	GetDirection(c *gin.Context)
	ForwardGeocode(c *gin.Context)
	ReverseGeocode(c *gin.Context)
	// DistanceMatrix(c *gin.Context)
}

func (controller *GoongController) Autocomplete(c *gin.Context) {

	// Parse the input query parameter
	input := c.DefaultQuery("input", "") // Default to an empty string if not provided
	limit := c.DefaultQuery("limit", "")
	// Parse the location query parameters if provided
	locationParam := c.DefaultQuery("location", "")
	var location *Location
	if locationParam != "" {
		// Split the location string into lat and lng
		parts := strings.Split(locationParam, ",")
		if len(parts) == 2 {
			lat, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
				return
			}
			lng, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
				return
			}
			location = &Location{Lat: lat, Lng: lng}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "location must be in 'lat,lng' format"})
			return
		}
	}

	// Validate input
	if input == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input text cannot be empty"})
		return
	}

	// Call the GoongService's Autocomplete method
	result, err := controller.service.AutoCompleteService(input, location, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, result)
}

func (controller *GoongController) GetPlaceDetail(c *gin.Context) {
	// Parse the place_id query parameter
	placeID := c.Query("place_id")

	// Validate place_id
	if placeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "place_id cannot be empty"})
		return
	}

	// Call the GoongService's GetPlaceDetail method
	result, err := controller.service.GetPlaceDetailService(placeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, result)
}

func (controller *GoongController) GetDirection(c *gin.Context) {

	var req DirectionRequest
	// Parse the origin and destination query parameters
	origin := c.Query("origin")
	destination := c.Query("destination")

	// Validate origin and destination
	if origin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "origin cannot be empty"})
		return
	}
	if destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "destination cannot be empty"})
		return
	}

	req.Origin = origin
	req.Destination = destination

	// Call the GoongService's GetDirection method
	result, err := controller.service.GetDirectionService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, result)
}

func (controller *GoongController) ForwardGeocode(c *gin.Context) {
	// Parse the address query parameter
	address := c.Query("address")

	// Validate address
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address cannot be empty"})
		return
	}

	// Call the GoongService's ForwardGeocode method
	result, err := controller.service.ForwardGeocodeService(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, result)
}

func (controller *GoongController) ReverseGeocode(c *gin.Context) {
	// Parse the address query parameter
	location := c.Query("latlng")

	// Validate address
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "location cannot be empty"})
		return
	}

	// Call the GoongService's ForwardGeocode method
	result, err := controller.service.ReverseGeocodeService(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, result)
}

// func (controller *GoongController) DistanceMatrix(c *gin.Context) {

// 	var req DistanceMatrixRequest
// 	// Parse the origins and destinations query parameters
// 	origins := c.QueryArray("origins")
// 	destinations := c.QueryArray("destinations")
// 	vehicle := c.DefaultQuery("vehicle", "car") // Default to "car" if vehicle is not provided

// 	// Validate origins and destinations
// 	if len(origins) == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "origins cannot be empty"})
// 		return
// 	}
// 	if len(destinations) == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "destinations cannot be empty"})
// 		return
// 	}

// 	// Ensure that the coordinates are properly formatted and URL-encoded
// 	encodedOrigins := encodeCoordinates(origins)
// 	encodedDestinations := encodeCoordinates(destinations)

// 	req.Origins = encodedOrigins
// 	req.Destinations = encodedDestinations
// 	req.Vehicle = vehicle

// 	// Call the GoongService's DistanceMatrix method
// 	result, err := controller.service.DistanceMatrixService(req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Return the result as JSON
// 	c.JSON(http.StatusOK, result)
// }
// func encodeCoordinates(coords []string) string {
// 	// Join the coordinates into a single string, separated by the pipe character (|), and URL-encode them
// 	encoded := url.QueryEscape(strings.Join(coords, "|"))
// 	return encoded
// }
