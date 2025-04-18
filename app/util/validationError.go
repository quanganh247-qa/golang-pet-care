package util

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func ErrorValidator(err error) map[string]any {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Error()), fe.Error()}
		}
		return gin.H{"code": "E", "message": "Validation Error", "errors": out}
	}
	return gin.H{"code": "E", "message": fmt.Sprintf("Validation Error: %v", err), "errors": []ApiError{}}
}

func msgForTag(tag string, defaultError string) string {
	switch tag {
	case "required":
		return "Field is required"
	case "min":
		return "Field is not enough length"
	case "max":
		return "Field exceeds maximum length"
	case "alphanum":
		return "Field must contain only letters and numbers"
	case "email":
		return "Email is not valid"
	case "method":
		return "Method must be one of the following: GET, POST, PUT, DELETE"
	default:
		return defaultError
	}
}
