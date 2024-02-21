package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func ParseMultipartFieldsToJson(
	r *http.Request,
) ([]byte, error) {
	var jsonStructure map[string]interface{} = make(map[string]interface{})

	for key, values := range r.PostForm {
		if len(values) > 1 {
			// this is for all the fields that have one key with multiple values
			jsonStructure[key] = values
		} else {
			var err error
			var number int64
			// Value[0] it's creating arrays of input values
			// in case there are two inputs with the same name
			// it will create an array with the same key and more values
			// check if int
			if number, err = strconv.ParseInt(values[0], 10, 64); err == nil {
				jsonStructure[key] = number
			}
			// check if checkbox or string
			if err != nil {
				if values[0] == "on" {
					jsonStructure[key] = true
				} else {
					jsonStructure[key] = values[0]
				}
			}
		}
	}
	// json marshal
	jsonString, err := json.Marshal(jsonStructure)
	if err != nil {
		return nil, fmt.Errorf("error marshaling the json structure %v", err)
	}

	return jsonString, err
}
