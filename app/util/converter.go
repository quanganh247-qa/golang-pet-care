package util

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/pgtype"
)

func ConvertedQuery(query map[string][]string) map[string]string {
	convertedQuery := make(map[string]string)
	for key, value := range query {
		if len(value) > 0 {
			convertedQuery[key] = value[0]
		}
	}
	return convertedQuery
}

type ParseOptions struct {
	Data    interface{}
	Message string
}

type ParseByteOptions struct {
	Data    []byte
	Message string
}

func ParseInterfaceToString(options ParseOptions) string {
	jsonData, err := json.Marshal(options.Data)
	if err != nil {
		log.Println("[ ERROR ] ParseInterfaceToString : ", options.Message, err)
		return ""
	}
	return string(jsonData)
}

func ParseStringToInterface[T any](options ParseOptions) *T {
	var jsonData T
	err := json.Unmarshal([]byte(options.Data.(string)), &jsonData)
	if err != nil {
		log.Println("[ ERROR ] ParseStringToInterface : ", options.Message, err)
		json.Unmarshal([]byte("{}"), &jsonData)
		return &jsonData
	}
	return &jsonData
}

func ParseByteToMapInterface(options ParseByteOptions) *map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal((options.Data), &result)
	if err != nil {
		log.Println("[ ERROR ] ParseByteToMapInterface : ", options.Message, err)
		json.Unmarshal([]byte("{}"), &result)
		return &result
	}
	return &result
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// Slugify converts a string to a slug format suitable for URLs.
func Slugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace all non-alphanumeric characters (except for hyphens) with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")
	return s
}

func GetRequiredNumber(ctx *gin.Context, fieldName string) (int64, error) {
	id := ctx.Query(fieldName)
	if id == "" {
		return 0, fmt.Errorf("%s is required in query", fieldName)
	}
	idNum, err := strconv.ParseInt(id, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("%s is not valid, Detail error : %v", fieldName, err)
	}
	return idNum, nil
}

func ParseMapInterfaceToString(data map[string]interface{}) string {
	contenJson, err := json.Marshal(data)
	if err != nil {
		log.Println("[ ERROR ] ParseInterfaceToString : ", err)
		return ""
	}
	return string(contenJson)
}

func ParseStringToTime(startDate string, endDate string) (time.Time, time.Time, error) {
	// string to time.Time
	parsedStartDate, err := time.Parse("2006-01-02 15:04:05", startDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start date format: %w", err)
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e01abc5 (pet schedule api)
	// If endDate is provided, parse it; otherwise, return a zero time for parsedEndDate
	var parsedEndDate time.Time
	if endDate != "" {
		parsedEndDate, err = time.Parse("2006-01-02 15:04:05", endDate)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end date format: %w", err)
		}
	}
	return parsedStartDate, parsedEndDate, nil
}

// StringToInterval converts a string duration (like "00:30:00") to a pgtype.Interval
func StringToInterval(duration string) (pgtype.Interval, error) {
	var interval pgtype.Interval

	// If duration is empty, return null interval
	if duration == "" {
		interval.Status = pgtype.Null
		return interval, nil
	}

	// Check if the format is HH:MM:SS
	parts := strings.Split(duration, ":")
	if len(parts) == 3 {
		// Parse hours, minutes, and seconds
		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return interval, fmt.Errorf("invalid hours in duration: %w", err)
		}

		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			return interval, fmt.Errorf("invalid minutes in duration: %w", err)
		}

		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return interval, fmt.Errorf("invalid seconds in duration: %w", err)
		}

		// Convert to microseconds
		totalMicroseconds := (int64(hours)*3600 + int64(minutes)*60 + int64(seconds)) * 1e6

		interval.Microseconds = totalMicroseconds
		interval.Status = pgtype.Present
		return interval, nil
	}

	// Otherwise, handle default time.ParseDuration format (e.g., "1h", "30m")
	parsed, err := time.ParseDuration(duration)
	if err != nil {
		return interval, fmt.Errorf("invalid duration format: %w", err)
	}

	// Convert to microseconds for PostgreSQL interval
	microseconds := parsed.Microseconds()

	interval.Microseconds = microseconds
	interval.Status = pgtype.Present

	return interval, nil
}

// Helper function to convert pgtype.Interval to string
func IntervalToString(interval pgtype.Interval) string {
	if interval.Status != pgtype.Present {
		return ""
	}

	// Convert microseconds back to duration
	duration := time.Duration(interval.Microseconds) * time.Microsecond

	// Add days (assuming 24 hours per day)
	duration += time.Duration(interval.Days) * 24 * time.Hour

	// Add months (approximate - assuming 30 days per month)
	duration += time.Duration(interval.Months) * 30 * 24 * time.Hour

	return duration.String()
}

// Generic function to map fields from one struct to another
func MapStructs(source interface{}, destination interface{}) error {
	srcValue := reflect.ValueOf(source)
	destValue := reflect.ValueOf(destination)

	// Ensure source and destination are pointers to structs
	if srcValue.Kind() != reflect.Ptr || srcValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("source must be a pointer to a struct")
	}
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a pointer to a struct")
	}

	// Dereference the pointers
	srcValue = srcValue.Elem()
	destValue = destValue.Elem()

	// Iterate over fields in the source struct
	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		srcFieldName := srcValue.Type().Field(i).Name

		// Find the corresponding field in the destination struct
		destField := destValue.FieldByName(srcFieldName)
		if !destField.IsValid() || !destField.CanSet() {
			continue // Skip if the field doesn't exist or can't be set
		}

		// Handle nested structs
		if srcField.Kind() == reflect.Struct && destField.Kind() == reflect.Struct {
			err := MapStructs(srcField.Addr().Interface(), destField.Addr().Interface())
			if err != nil {
				return err
			}
		} else if srcField.Type() == destField.Type() {
			// Copy the value if the types match
			destField.Set(srcField)
		}
	}

	return nil
}
<<<<<<< HEAD
=======
=======
>>>>>>> 79a3bcc (medicine api)
	// Parse the EndDate string to time.Time (if provided)
	parsedEndDate, err := time.Parse("2006-01-02 15:04:05", endDate)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end date format: %w", err)
	}

	return parsedStartDate, parsedEndDate, nil
}
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 79a3bcc (medicine api)
=======

// StringToInterval converts a string duration (like "00:30:00") to a pgtype.Interval
func StringToInterval(duration string) (pgtype.Interval, error) {
	var interval pgtype.Interval

	// If duration is empty, return null interval
	if duration == "" {
		interval.Status = pgtype.Null
		return interval, nil
	}

	// Check if the format is HH:MM:SS
	parts := strings.Split(duration, ":")
	if len(parts) == 3 {
		// Parse hours, minutes, and seconds
		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return interval, fmt.Errorf("invalid hours in duration: %w", err)
		}

		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			return interval, fmt.Errorf("invalid minutes in duration: %w", err)
		}

		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return interval, fmt.Errorf("invalid seconds in duration: %w", err)
		}

		// Convert to microseconds
		totalMicroseconds := (int64(hours)*3600 + int64(minutes)*60 + int64(seconds)) * 1e6

		interval.Microseconds = totalMicroseconds
		interval.Status = pgtype.Present
		return interval, nil
	}

	// Otherwise, handle default time.ParseDuration format (e.g., "1h", "30m")
	parsed, err := time.ParseDuration(duration)
	if err != nil {
		return interval, fmt.Errorf("invalid duration format: %w", err)
	}

	// Convert to microseconds for PostgreSQL interval
	microseconds := parsed.Microseconds()

	interval.Microseconds = microseconds
	interval.Status = pgtype.Present

	return interval, nil
}

// Helper function to convert pgtype.Interval to string
func IntervalToString(interval pgtype.Interval) string {
	if interval.Status != pgtype.Present {
		return ""
	}

	// Convert microseconds back to duration
	duration := time.Duration(interval.Microseconds) * time.Microsecond

	// Add days (assuming 24 hours per day)
	duration += time.Duration(interval.Days) * 24 * time.Hour

	// Add months (approximate - assuming 30 days per month)
	duration += time.Duration(interval.Months) * 30 * 24 * time.Hour

	return duration.String()
}
>>>>>>> 272832d (redis cache)
=======
>>>>>>> 883d5b3 (update treatment)
=======
>>>>>>> 79a3bcc (medicine api)
