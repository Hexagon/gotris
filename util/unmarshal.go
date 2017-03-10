package util

import "encoding/json"

func Unmarshal(b []byte) (error, map[string]interface{}) {

	// Create a generic interface to unmarshal into
	var jsonDataInterface interface{}

	// This is brigde (http://i.imgur.com/W0N40BK.jpg)
	err := json.Unmarshal(b, &jsonDataInterface)
	if err != nil {
		return err, nil
	}

	// Type assertion to get a map of string => interface
	jsonData := jsonDataInterface.(map[string]interface{})

	return nil, jsonData
}
