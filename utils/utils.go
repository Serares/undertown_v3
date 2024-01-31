package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
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

func CreateDisplayPrice(price int64) string {
	return strconv.Itoa(int(price)) + " €"
}

func CreateDisplayCreatedAt(createdAt time.Time) string {
	passedTime := int(time.Since(createdAt).Hours() / 24)
	if passedTime < 0 {
		return "Adaugat recent"
	}

	if passedTime > 10 {
		return ""
	}

	return fmt.Sprintf("Adaugat cu %d zile in urma", passedTime)
}

// ❗ The GET_PROPERTY backend should search for the propertyId query string
func CreatePropertyPath(title, humanReadableId string) string {
	replacedString := strings.ReplaceAll(title, " ", string('%'))
	return strings.ToLower(replacedString) + strings.Join([]string{"?propertyId=", humanReadableId}, "")
}
