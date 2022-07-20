package routes

import (
	"net/http"
	"time"

	mainlogic "github.com/MarcBernstein0/match-display/src/main-logic"
	"github.com/MarcBernstein0/match-display/src/models"
	"github.com/gin-gonic/gin"
)

func JSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

type Date struct {
	Date time.Time `form:"date" binding:"required" time_format:"2006-01-02"`
}

func MatchesGET(fetchData mainlogic.FetchData) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var date Date
		if err := c.BindQuery(&date); err != nil {
			errResponse := models.ErrorResponse{
				Message:      "did not fill out required 'date' query field",
				ErrorMessage: err.Error(),
			}
			c.JSON(http.StatusBadRequest, errResponse)
			return
		}
		// fmt.Printf("%+v\n", date)
		// get date
		// call tournaments
		tournaments, err := fetchData.FetchTournaments(date.Date.Format("2006-01-02"))
		if err != nil {
			errResponse := models.ErrorResponse{
				Message:      "failed to get tournament data",
				ErrorMessage: err.Error(),
			}
			c.JSON(http.StatusInternalServerError, errResponse)
			return
		}
		// fmt.Printf("tournaments %+v\n", tournaments)
		// call particiapnts
		participants, err := fetchData.FetchParticipants(tournaments)
		if err != nil {
			errResponse := models.ErrorResponse{
				Message:      "failed to get participant data",
				ErrorMessage: err.Error(),
			}
			c.JSON(http.StatusInternalServerError, errResponse)
			return
		}
		// fmt.Printf("list of participants %+v\n", participants)
		// call matches
		matches, err := fetchData.FetchMatches(participants)
		if err != nil {
			errResponse := models.ErrorResponse{
				Message:      "failed to get match data",
				ErrorMessage: err.Error(),
			}
			c.JSON(http.StatusInternalServerError, errResponse)
			return
		}
		// fmt.Printf("list of matches %+v\n", matches)
		// return matches
		c.JSON(http.StatusOK, matches)
	}
	return fn
}
