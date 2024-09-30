package util

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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

// func ParseInterfaceToStringChannel(options ParseOptions) string {
// 	jsonData, err := json.Marshal(options.Data)
// 	if err != nil {
// 		log.Println("[ ERROR ] ParseInterfaceToString : ", options.Message, err)
// 		return ""
// 	}
// 	return string(jsonData)
// }

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
