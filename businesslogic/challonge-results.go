package businesslogic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/MarcBernstein0/match-display/ulits/errorhandling"
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
	// _, f, _, _ := runtime.Caller(0)
	// rawConfig, err := ioutil.ReadFile(filepath.Dir(f) + "/../config.json") // will change path, problem with running test and running from main
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// err = json.Unmarshal(rawConfig, &config)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	config.Username = os.Getenv("USER_NAME")
	config.ApiKey = os.Getenv("API_KEY")
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
	if err = errorhandling.HandleError("failed to create request.", err); err != nil {
		return result{
			data: nil,
			err:  err,
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
	// TODO: check if result back http code
	if err = errorhandling.HandleError("failed to received response from challonge api.", err); err != nil {
		return result{
			data: nil,
			err:  err,
		}
	}
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	if err = errorhandling.HandleError("error when reading response body.", err); err != nil {
		return result{
			data: nil,
			err:  err,
		}
	}
	var tData []map[string]map[string]interface{}
	err = json.Unmarshal([]byte(resData), &tData)
	if err := errorhandling.HandleError("failed to unmarshal json data", err); err != nil {
		return result{
			data: nil,
			err:  err,
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
