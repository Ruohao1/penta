package utils

import (
	"encoding/json"
)

func FromJSON[T any](data []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func ToJSON[T any](v T) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
