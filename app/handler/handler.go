package handler

import (
	"database/sql"
	"net/http"
)

func TestAPI(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	resposeJSON(w, http.StatusOK, "halo dari server")
}
