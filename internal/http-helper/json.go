package httphelper

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statuscode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	return nil
}

func WriteError(w http.ResponseWriter, status int, message string) {
	if err := WriteJSON(w, status, map[string]string{
		"error": message,
	}); err != nil {

		log.Println(err)
	}
}
