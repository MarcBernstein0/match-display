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
	createhtml "github.com/MarcBernstein0/match-display/ulits/create-html"
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
		// w.WriteHeader(http.StatusInternalServerError)
		errorhandling.ErrorResponse(w, "Error with getting data", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(tournamentData)
	if ok, err := errorhandling.HandleError("failed to convert tournament data to json", err); ok {
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
			_, err := errorhandling.HandleError("Wrong type provided for field"+unmarshalErr.Field, err)
			log.Println(err)
			errorhandling.ErrorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
			return
		} else {
			_, err := errorhandling.HandleError("Bad request", err)
			log.Println(err)
			errorhandling.ErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	// fmt.Println(tournaments)
	matches, err := tournaments.GetMatches()
	if ok, err := errorhandling.HandleError("failed to get tournament data in controler method GetMatchData", err); ok {
		log.Println(err)
		errorhandling.ErrorResponse(w, "Error with getting data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	htmlString, err := createhtml.CreateHtml(matches)
	if ok, err := errorhandling.HandleError("failed to convert tournament data to html", err); ok {
		log.Printf("Error with getting data\n%v", err)
		// w.WriteHeader(http.StatusInternalServerError)
		errorhandling.ErrorResponse(w, "Error with creating html string", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(htmlString))
}
