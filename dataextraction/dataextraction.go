package dataextraction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
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

const API_URL = "https://api.challonge.com/v1"

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

type result struct {
	data []map[string]map[string]interface{}
	err  error
}

/* Generic function that calls the cahllonge api and returns the body of the response. Only builds get requests as no post calls will happen
args:
	apiPath string path of the api call
	params map[string]string all the parameters that will be passed into the
	                             request where key is the parameter and value is the parameter value
return:
	map[string]map[string]interface{} the fully built request ready to be sent
	error errors that occur when building the request
*/
func challongeApiCall(client HTTPClient, apiPath string, params map[string]string) result {
	fullAPIPath := fmt.Sprintf("%s/%s.json", API_URL, apiPath)
	req, err := http.NewRequest("GET", fullAPIPath, nil)
	if err != nil {
		return result{
			data: nil,
			err:  fmt.Errorf("failed to create request.\n%v", err),
		}
	}
	// build query
	q := req.URL.Query()
	q.Add("api_key", config.ApiKey)
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	res, err := client.Do(req)
	if err != nil {
		return result{
			data: nil,
			err:  fmt.Errorf("failed to received response from challonge api. \n%v", err)}
	}
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result{
			data: nil,
			err:  fmt.Errorf("error when reading response body.\n%v", err),
		}
	}
	var tData []map[string]map[string]interface{}
	if err = json.Unmarshal([]byte(resData), &tData); err != nil {
		return result{
			data: nil,
			err:  fmt.Errorf("failed to unmarshal json data\n%v", err),
		}
	}
	return result{
		data: tData,
		err:  nil,
	}

}

func challongeApiMultiCall(tournament int, client HTTPClient, apiPath string, params map[string]string, resultsChan chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	res := challongeApiCall(client, apiPath, params)
	resultsChan <- res
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

	// parameters to pass in
	params := map[string]string{
		"state": "in_progress",
	}

	// create request to client
	res := challongeApiCall(client, "tournaments", params)
	if res.err != nil {
		return nil, fmt.Errorf("request failed in getTournaments\n. %v", res.err)
	}

	for _, elem := range res.data {
		if tournamentID, ok := elem["tournament"]["id"].(float64); ok {

			if gameName, ok := elem["tournament"]["game_name"].(string); ok {
				tournaments[int(tournamentID)] = gameName
			} else {
				return nil, fmt.Errorf("type for game_name did not match what was expected. Expected='string' got=%T", gameName)
			}
		} else {
			return nil, fmt.Errorf("type for tournament ID did not match what was expected. Expected='float64' got=%T", tournamentID)
		}
	}

	return tournaments, nil
}

func getParticipants(tournaments map[int]string, client HTTPClient) (map[int]string, error) {
	participants := make(map[int]string)

	// cResponse := make(chan []map[string]map[string]interface{})
	// cError := make(chan error)
	// var wg sync.WaitGroup
	// for k, v := range tournaments {
	// 	wg.Add(1) // tells the waitgroup that there is no 1 pending operation
	// }

	apiPath := fmt.Sprintf("tournaments/%s/participants", "10469768")
	res := challongeApiCall(client, apiPath, nil)
	if res.err != nil {
		return nil, fmt.Errorf("request failed in getTouranments call\n. %v", res.err)
	}

	for _, elem := range res.data {
		if participantID, ok := elem["participant"]["id"].(float64); ok {

			if name, ok := elem["participant"]["name"].(string); ok {
				participants[int(participantID)] = name
			} else {
				return nil, fmt.Errorf("type for 'name' did not match what was expected. Expected='string' got=%T", name)
			}
		} else {
			return nil, fmt.Errorf("type for 'participantID' did not match what was expected. Expected='float64' got=%T", participantID)
		}
	}

	return participants, nil
}

func TestIntegration() {
	fmt.Println(getTournaments(client))
	fmt.Println(getParticipants(
		map[int]string{
			10469768: "Test2",
		},
		client,
	))
}
