package businesslogic

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
			err:  fmt.Errorf("failed to received response from challonge api.\n%v", err)}
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

func challongeApiMultiCall(client HTTPClient, apiPath string, params map[string]string, resultsChan chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	res := challongeApiCall(client, apiPath, params)
	resultsChan <- res
}
