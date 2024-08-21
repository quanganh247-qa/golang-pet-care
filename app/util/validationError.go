package util

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/quanganh247-qa/go-blog-be/app/util/enums"
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
		message:=""
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Error()), fe.Error()}
			message+=fmt.Sprintf("%s : %s, ",fe.Field(),msgForTag(fe.Tag(), fe.Error()))
		}
		errRes := gin.H{"code": "E", "message": "Validation Error > "+message, "errors": out}
		return errRes
	}
	errRes := gin.H{"code": "E", "message": fmt.Sprintf("Validation Error Input Request Body!!!, %v",err), "errors": []ApiError{}}
	return errRes
}

func ParseSqlErr(err error) error {
	var sqlErr *pgconn.PgError
	if errors.As(err, &sqlErr) {
		if msg, ok := enums.ErrCodeSql[sqlErr.Code]; ok {
			message := fmt.Sprintf("%s  [Chi tiết] : %s", msg, sqlErr.Message)
			return errors.New(message)
		} else {
			return errors.New("Lỗi không xác định từ cơ sở dữ liệu!!!")
		}
	}
	return err
}

func ParseSqlErrByCode(err error, mapErrCodeMessages map[string]string) (errSql error) {
	var sqlErr *pgconn.PgError
	if errors.As(err, &sqlErr) {
		if msg, ok := mapErrCodeMessages[sqlErr.Code]; ok {
			return errors.New(msg)
		} else if msg, ok := enums.ErrCodeSql[sqlErr.Code]; ok {
			message := fmt.Sprintf("%s  [Chi tiết] : %s", msg, sqlErr.Message)
			return errors.New(message)
		} else {
			return fmt.Errorf("Lỗi không xác định từ cơ sở dữ liệu!!! : %s", sqlErr.Message)
		}
	}
	return fmt.Errorf("Error From Other Resource : %s", err.Error())
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
