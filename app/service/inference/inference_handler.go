package inference

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

// InferenceHandler handles pet breed identification using Roboflow API
type InferenceHandler struct {
	roboflowClient *RoboflowClient
	// Default model IDs to use if none specified
	defaultCatModelID string
	defaultDogModelID string
}

// NewInferenceHandler creates a new inference handler
func NewInferenceHandler(apiKey string) *InferenceHandler {
	if apiKey == "" {
		// Log warning if no API key provided, but still create the handler
		// It will return errors when actually called
		fmt.Println("WARNING: No Roboflow API key provided in configuration")
	}

	return &InferenceHandler{
		roboflowClient:    NewRoboflowClient(apiKey, ""),
		defaultCatModelID: "cat-breeds-cbvra/1",
		defaultDogModelID: "dogsdetector/2",
	}
}

// RegisterRoutes registers all routes related to inference
func (h *InferenceHandler) RegisterRoutes(router *gin.RouterGroup) {
	infGroup := router.Group("/inference")
	{
		infGroup.POST("/detect", h.DetectFromFile)
		infGroup.POST("/detect-base64", h.DetectFromBase64)
	}
}

// DetectFromFile handles inference requests with file upload
// @Summary Detect objects in an uploaded image
// @Description Upload an image to detect objects using Roboflow models
// @Tags inference
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image to analyze"
// @Param model_id formData string false "Roboflow model ID (default: oxford-pets/3)"
// @Success 200 {object} inference.InferenceResult
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /inference/detect [post]
func (h *InferenceHandler) DetectFromFile(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image uploaded"})
		return
	}

	// Get model ID from request or use default
	modelID := c.PostForm("model_id")
	if modelID == "" {
		modelID = "oxford-pets/3" // Default model from the JS example
	}

	// Get breed type
	breed := c.PostForm("breed")

	// Save uploaded file temporarily
	tempPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
		return
	}
	defer os.Remove(tempPath) // Clean up temp file

	// Call Roboflow API
	result, err := h.roboflowClient.Infer(tempPath, modelID, breed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Inference failed: %v", err)})
		return
	}

	// Return result
	c.JSON(http.StatusOK, result)
}

// DetectFromBase64 handles inference requests with base64-encoded images
// @Summary Detect objects in a base64-encoded image
// @Description Send a base64-encoded image to detect objects using Roboflow models
// @Tags inference
// @Accept application/json
// @Produce json
// @Param model_id query string false "Roboflow model ID"
// @Param breed query string false "Animal breed (cat or dog)"
// @Param image body string true "Base64-encoded image"
// @Success 200 {object} inference.InferenceResult
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /inference/detect-base64 [post]
func (h *InferenceHandler) DetectFromBase64(c *gin.Context) {
	// Get model ID and breed from query params
	breed := c.Query("breed")
	modelID := c.Query("model_id")

	// Set default model ID based on breed if not provided
	if modelID == "" {
		if breed == "cat" {
			modelID = h.defaultCatModelID
		} else if breed == "dog" {
			modelID = h.defaultDogModelID
		} else {
			// If no valid breed, use a default
			modelID = "oxford-pets/3"
		}
	}

	dataImage, _, err := util.HandleImageUpload(c, "image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to handle image upload: %v", err)})
		return
	}

	bodyData := base64.StdEncoding.EncodeToString(dataImage)

	// Call Roboflow API with the base64 data
	result, err := h.roboflowClient.InferFromBytes([]byte(bodyData), modelID, breed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Inference failed: %v", err)})
		return
	}

	// Return result
	c.JSON(http.StatusOK, result)
}

// DetectPet handles unified pet detection (both cats and dogs)
// @Summary Detect both cats and dogs in an uploaded image
// @Description Upload an image to detect cat or dog breeds
// @Tags inference
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Pet image to analyze"
// @Param breed formData string true "Animal type (cat or dog)"
// @Success 200 {object} inference.InferenceResult
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /inference/detect-pet [post]
func (h *InferenceHandler) DetectPet(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image uploaded"})
		return
	}

	// Get breed type (cat or dog)
	breed := c.PostForm("breed")
	if breed != "cat" && breed != "dog" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid breed type. Must be 'cat' or 'dog'"})
		return
	}

	// Determine which model to use based on the breed
	var modelID string
	if breed == "cat" {
		modelID = h.defaultCatModelID
	} else {
		modelID = h.defaultDogModelID
	}

	// Save uploaded file temporarily
	tempPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
		return
	}
	defer os.Remove(tempPath) // Clean up temp file

	// Read the file data
	imageData, err := os.ReadFile(tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read file: %v", err)})
		return
	}

	// Convert to base64 for the API
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	// Call Roboflow API with the image data
	result, err := h.roboflowClient.InferFromBytes([]byte(base64Data), modelID, breed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Inference failed: %v", err)})
		return
	}

	// Return the full result object which contains all fields
	c.JSON(http.StatusOK, result)
}
