package businesslogic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`
}

var (
	config Config
	client HTTPClient
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const API_URL = "https://api.challonge.com/v1/tournaments.json"

func init() {
	// Open config file
	rawConfig, err := ioutil.ReadFile("/home/marc/Projects/match-display/config.json") // will change path, problem with running test and running from main
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Fatalln(err)
	}

	client = &http.Client{}
}

/* calls challenonge api to get all running tournaments
   created recently
   args:
   	none

   returns:
	map[int]string	mapping of tournament IDs and name of the game
	error
*/
func getTournaments(client HTTPClient) (map[int]string, error) {
	// map of tournamentIDs and game names
	tournaments := make(map[int]string)

	// create request to client
	req, err := http.NewRequest("GET", API_URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request.\n%v", err)
	}

	// add request parameters
	q := req.URL.Query()
	q.Add("api_key", config.ApiKey)
	q.Add("state", "in_progress")
	req.URL.RawQuery = q.Encode()

	// call api client and handle response
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response from api failed.\n%v", err)
	}
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error when reading response body.\n%v", err)
	}
	var tData []map[string]map[string]interface{}
	if err = json.Unmarshal([]byte(resData), &tData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json data\n%v", err)
	}
	for _, elem := range tData {
		if tournamentID, ok := elem["tournament"]["id"].(float64); ok {

			if gameName, ok := elem["tournament"]["game_name"].(string); ok {
				tournaments[int(tournamentID)] = gameName
			}
		} else {
			return nil, fmt.Errorf("type for tournament ID did not match what was expected. Expected='float64' got=%T", elem["tournament"]["id"])
		}
	}

	return tournaments, nil
}

func TestIntegration() {
	fmt.Println(getTournaments(client))
}
