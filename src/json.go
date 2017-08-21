package main

import (
	"encoding/json"
	"errors"
	// "fmt"
)

func validateJson(data *JsonData) error {
	for _, value := range *data {
		if value == nil {
			return errors.New("nil")
		}
	}
	return nil
}

func parseJson(body []byte) (*JsonData, bool) {
	var data JsonData
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, false
		// 	err = ValidateJson(&data)
	}
	return &data, true
}
