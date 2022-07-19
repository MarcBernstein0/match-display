package routes

import (
	mainlogic "github.com/MarcBernstein0/match-display/src/main-logic"
	"github.com/gin-gonic/gin"
)

func RouteSetup(fetchData mainlogic.FetchData) *gin.Engine {
	r := gin.Default()

	r.GET("/health", HealthCheck)
	// r.GET("/matches", MatchesGET(fetchData))

	return r
}
