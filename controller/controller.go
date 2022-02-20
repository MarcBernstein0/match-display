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
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// specifiy status code
	w.WriteHeader(http.StatusOK)

	// update response writer
	io.WriteString(w, `{"alive": true}`)
}

func GetMatchData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	v := r.URL.Query()

	var dateStr string
	// check if the date string is empty or not
	if len(v.Get("date")) != 0 {
		date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
		if err = errorhandling.HandleError("not a valid date", err); err != nil {
			log.Printf("Date sent from caller was not a valid date\n%v", err)
			errorhandling.ErrorResponse(w, "Did not send a valid date", http.StatusBadRequest)
			return
		}
		dateStr = date.Format("2006-01-02")
	}

	fmt.Println(dateStr)
	tournamentData, err := businesslogic.GetTournamentData(dateStr)
	if err = errorhandling.HandleError("failed to get tournament data in controler method GetTournamentData", err); err != nil {
		log.Printf("Error with getting data\n%v", err)
		errorhandling.ErrorResponse(w, "Error with getting data", http.StatusInternalServerError)
		return
	}
	matches, err := tournamentData.GetMatches()
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
