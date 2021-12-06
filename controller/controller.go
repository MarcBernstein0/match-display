package controller

import (
	"encoding/json"
	"errors"
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
	if err = errorhandling.HandleError("failed to get tournament data in controler method GetTournamentData", err); err != nil {
		log.Printf("Error with getting data\n%v", err)
		// w.WriteHeader(http.StatusInternalServerError)
		errorhandling.ErrorResponse(w, "Error with getting data", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(tournamentData)
	if err = errorhandling.HandleError("failed to convert tournament data to json", err); err != nil {
		log.Printf("Error with getting data\n%v", err)
		// w.WriteHeader(http.StatusInternalServerError)
		errorhandling.ErrorResponse(w, "Error with Marshalling data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetMatchData(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorhandling.ErrorResponse(w, "Content Type was not application/json", http.StatusUnsupportedMediaType)
		return
	}

	var tournaments businesslogic.Tournaments
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&tournaments)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			err := errorhandling.HandleError("Wrong type provided for field"+unmarshalErr.Field, err)
			log.Println(err)
			errorhandling.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
			return
		} else {
			err := errorhandling.HandleError("Bad request", err)
			log.Println(err)
			errorhandling.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// fmt.Println(tournaments)
	matches, err := tournaments.GetMatches()
	if err = errorhandling.HandleError("failed to get tournament data in controler method GetMatchData", err); err != nil {
		log.Println(err)
		errorhandling.ErrorResponse(w, "Error with getting data", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(matches)
	if err = errorhandling.HandleError("failed to convert match data to json", err); err != nil {
		log.Println(err)
		errorhandling.ErrorResponse(w, "Error with marshling data", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
