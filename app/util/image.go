package util

import (
	"fmt"
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
}
