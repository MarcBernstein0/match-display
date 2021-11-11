package businesslogic

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`
}

var config Config

func init() {
	// Open config file
	rawConfig, err := ioutil.ReadFile("../config.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Fatalln(err)
	}

	// read jsonFile as byte
}

/* calls challenonge api to get all running tournaments
   created recently
   args:
   	none

   returns:
	map[int]string	mapping of tournament IDs and name of the game
	error
*/
func getTournaments() (map[int]string, error) {

	return nil, nil
}
