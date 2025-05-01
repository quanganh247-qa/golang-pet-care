package util

import (
	"encoding/json"
)

// ToJSON converts a value to JSON bytes
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// FromJSON parses JSON bytes into the given value
func FromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
