package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func ReplyError(w http.ResponseWriter, r *http.Request, status int, message string) {
	http.Error(w, message, status)
}

func ReplySuccess(w http.ResponseWriter, r *http.Request, status int, message interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(message)
	if err != nil {

		w.WriteHeader(500)
		return fmt.Errorf("error marshalling JSON: %s", err)
	}
	w.WriteHeader(status)
	w.Write(dat)
	return nil
}

func CheckIfStructIsEmpty(s interface{}) bool {
	// Get the value of the struct
	val := reflect.ValueOf(s)

	// If the struct is a pointer, find the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Iterate over all fields of the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			return false // Found a non-zero field
		}
	}
	return true
}
