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

/* Generic function that builds and http request. Only builds get requests as no post calls will happen
args:
	apiPath string path of the api call
	params map[string]string all the parameters that will be passed into the
	                             request where key is the parameter and value is the parameter value
return:
	*http.Request the fully built request ready to be sent
	error errors that occur when building the request
*/
func httpQueryBuilder(apiPath string, params map[string]string) (*http.Request, error) {
	fullAPIPath := fmt.Sprintf("%s/%s.json", API_URL, apiPath)
	req, err := http.NewRequest("GET", fullAPIPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request.\n%v", err)
	}
	// build query
	q := req.URL.Query()
	q.Add("api_key", config.ApiKey)
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
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
	req, err := httpQueryBuilder("tournaments", params)
	if err != nil {
		return nil, fmt.Errorf("req failed in getTouranments call\n. %v", err)
	}

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

func getParticipants(tournaments map[int]string, client HTTPClient) (map[int]string, error) {
	apiPath := fmt.Sprintf("tournaments/%s/participants", "10469768")
	req, err := httpQueryBuilder(apiPath, nil)
	if err != nil {
		return nil, fmt.Errorf("req failed in getTouranments call\n. %v", err)
	}
	fmt.Println(req.URL.String())
	return nil, nil
}

func TestIntegration() {
	fmt.Println(getTournaments(client))
}
