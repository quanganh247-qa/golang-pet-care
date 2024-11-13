package util

import (
	"fmt"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func HandleImageUpload(ctx *gin.Context, formFieldName string) ([]byte, string, error) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse multipart form: %w", err)
	}

	// Handle image file
	file, header, err := ctx.Request.FormFile(formFieldName)
	if err != nil {
		return nil, "", fmt.Errorf("image is required")
	}
	defer file.Close()

	// Get the original image name
	originalImageName := header.Filename

	// Read the file content into a byte array
	dataImage, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read data image: %w", err)
	}

	return dataImage, originalImageName, nil
=======
	"io"
	"mime/multipart"
	"path/filepath"
=======
	"io/ioutil"

	"github.com/gin-gonic/gin"
>>>>>>> 473cd1d (uplaod image method)
)

func HandleImageUpload(ctx *gin.Context, formFieldName string) ([]byte, string, error) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse multipart form: %w", err)
	}

<<<<<<< HEAD
	return buffer, originalFilename, nil
>>>>>>> 0fb3f30 (user images)
=======
	// Handle image file
	file, header, err := ctx.Request.FormFile(formFieldName)
	if err != nil {
		return nil, "", fmt.Errorf("image is required")
	}
	defer file.Close()

	// Get the original image name
	originalImageName := header.Filename

	// Read the file content into a byte array
	dataImage, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read data image: %w", err)
	}

	return dataImage, originalImageName, nil
>>>>>>> 473cd1d (uplaod image method)
=======
	"io"
	"mime/multipart"
	"path/filepath"
)

// Helper functions
func IsValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	return validTypes[contentType]
}

func ProcessImage(file multipart.File, fileHeader *multipart.FileHeader) ([]byte, string, error) {
	// Get original filename
	originalFilename := filepath.Base(fileHeader.Filename)

	// Read file into byte array
	buffer := make([]byte, fileHeader.Size)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, "", fmt.Errorf("error reading file: %v", err)
	}

	return buffer, originalFilename, nil
>>>>>>> 0fb3f30 (user images)
}
