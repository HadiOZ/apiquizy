package handler

import (
	"apiquizyfull/app/model"
	"database/sql"
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

func TestAPI(w http.ResponseWriter, r *http.Request, db sql.DB) {
	resposeJSON(w, http.StatusOK, "halo dari server")
}

func HandlerTest(w http.ResponseWriter, r *http.Request, db sql.DB) {
	res := model.SelectAllCalon(&db)
	resposeJSON(w, http.StatusOK, res)
}
