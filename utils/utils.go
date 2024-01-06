package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
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
