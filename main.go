package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MarcBernstein0/match-display/controller"
	"github.com/gorilla/mux"
)

const BASE_ROUTE = "/match-display/v1"

func main() {
	// currentTime := time.Now()
	// tournaments, err := businesslogic.GetTournamentData(currentTime.Format("2006-01-02"))
	// if ok, err := errorhandling.HandleError("failed when calling GetTournamentData", err); ok {
	// 	panic(err)
	// }
	// fmt.Println(tournaments)
	// matches, err := tournaments.GetMatches()
	// if ok, err := errorhandling.HandleError("failed when calling GetMatches", err); ok {
	// 	panic(err)
	// }
	// fmt.Println(matches)
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Get all tournaments for a specific date
	r := mux.NewRouter()
	r.HandleFunc(fmt.Sprintf("%s/health-check", BASE_ROUTE), controller.HealthCheck).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("%s/tournaments", BASE_ROUTE), controller.GetTournamentData).Methods(http.MethodGet)
	r.HandleFunc(fmt.Sprintf("%s/matches", BASE_ROUTE), controller.GetMatchData).Methods(http.MethodPost)

	// Get matches
	fmt.Println("Listening on :" + port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}

}
