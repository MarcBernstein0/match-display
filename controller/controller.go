package controller

import (
	"io"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// specifiy status code
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// update response writer
	io.WriteString(w, `{"alive": true}`)
}

func GetTournamentData(w http.ResponseWriter, r *http.Request) {

}
