package main

import (
	"fmt"
	"time"

	"github.com/MarcBernstein0/match-display/businesslogic"
	"github.com/MarcBernstein0/match-display/ulits/errorhandling"
)

func main() {
	fmt.Println("Running program")
	currentTime := time.Now()
	tournaments, err := businesslogic.GetTournamentData(currentTime.Format("2006-01-02"))
	if ok, err := errorhandling.HandleError("failed when calling GetTournamentData", err); ok {
		panic(err)
	}
	fmt.Println(tournaments)
	matches, err := tournaments.GetMatches()
	if ok, err := errorhandling.HandleError("failed when calling GetMatches", err); ok {
		panic(err)
	}
	fmt.Println(matches)
}
