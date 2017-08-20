package main

import (
	"encoding/json"
	"errors"
)

type JsonData map[string]interface{}
type JsonDataArray map[string][]JsonData

func ValidateJson(data *JsonData) error {
	for _, value := range *data {
		if value == nil {
			return errors.New("nil")
		}
	}
	return nil
}

func readRequstJson(body []byte) (JsonData, error) {
	var data JsonData
	err := json.Unmarshal(body, &data)
	if err != nil {
		err = ValidateJson(&data)
	}
	return data, err
}
