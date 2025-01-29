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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Error()), fe.Error()}
		}
		return gin.H{"code": "E", "message": "Validation Error", "errors": out}
	}
	return gin.H{"code": "E", "message": fmt.Sprintf("Validation Error: %v", err), "errors": []ApiError{}}
=======
		message := ""
=======
>>>>>>> e859654 (Elastic search)
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Error()), fe.Error()}
		}
		return gin.H{"code": "E", "message": "Validation Error", "errors": out}
	}
<<<<<<< HEAD
	errRes := gin.H{"code": "E", "message": fmt.Sprintf("Validation Error Input Request Body!!!, %v", err), "errors": []ApiError{}}
	return errRes
>>>>>>> 3bf345d (happy new year)
=======
	return gin.H{"code": "E", "message": fmt.Sprintf("Validation Error: %v", err), "errors": []ApiError{}}
>>>>>>> e859654 (Elastic search)
=======
		message := ""
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Error()), fe.Error()}
			message += fmt.Sprintf("%s : %s, ", fe.Field(), msgForTag(fe.Tag(), fe.Error()))
		}
		errRes := gin.H{"code": "E", "message": "Validation Error > " + message, "errors": out}
		return errRes
	}
	errRes := gin.H{"code": "E", "message": fmt.Sprintf("Validation Error Input Request Body!!!, %v", err), "errors": []ApiError{}}
	return errRes
>>>>>>> 3bf345d (happy new year)
}

func msgForTag(tag string, defaultError string) string {
	switch tag {
	case "required":
		return "Trường này không được để trống"
	case "min":
		return "Trường này chưa đủ độ dài tối thiểu"
	case "max":
		return "Trường này vượt quá độ dài tối đa"
	case "alphanum":
		return "Trường này chỉ được chứa chữ và số"
	case "email":
		return "Email không hợp lệ"
	case "currency":
		return "Loại tiền phải thuộc các giá trị sau : USD, EUR,"
	case "method":
		return "Method phải thuộc các giá trị sau : GET, POST, PUT, DELETE"
	default:
		return defaultError
	}
}
