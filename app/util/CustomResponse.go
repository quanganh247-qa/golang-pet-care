package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(message string, data interface{}) gin.H {
	successRes := gin.H{
		"code":    "S",
		"message": message,
		"data":    data,
	}

	return successRes
}
func SuccessResponseConnector(message string, data interface{}, debug interface{}) gin.H {
	successRes := gin.H{
		"code":    "S",
		"message": message,
		"data":    data,
		"debug":   debug,
	}

	return successRes
}

func ErrorResponse(err error) gin.H {
	errRes := gin.H{"message": err.Error(), "code": "E"}

	return errRes
}

func CustomError(errType string, err error) error {
	return fmt.Errorf("%s : %w", errType, err)
}

// SendErrorResponse sends an error response with the given status code, message, and error details.
