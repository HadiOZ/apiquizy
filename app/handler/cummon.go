package handler

import (
	"encoding/json"
	"net/http"
)

func resposeJSON(w http.ResponseWriter, status int, playload interface{}) {

	response, err := json.Marshal(playload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func resposeErrorJSON(w http.ResponseWriter, status int, message string) {
	resposeJSON(w, status, map[string]string{"message": message})
}
