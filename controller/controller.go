package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/MarcBernstein0/match-display/businesslogic"
	"github.com/MarcBernstein0/match-display/ulits/errorhandling"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// specifiy status code
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// update response writer
	io.WriteString(w, `{"alive": true}`)
}

func GetTournamentData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	currentTime := time.Now().Format("2006-01-02")
	fmt.Println(currentTime)
	tournamentData, err := businesslogic.GetTournamentData(currentTime)
	if ok, err := errorhandling.HandleError("failed to get tournament data in controler method GetTournamentData", err); ok {
		log.Printf("Error with getting data\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(tournamentData)
	if ok, err := errorhandling.HandleError("failed to convert tournament data to json", err); ok {
		log.Printf("Error with getting data\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
