package businesslogic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	Username string `json:"username"`
	ApiKey   string `json:"api_key"`

	Server string `json:"server"`
	Port   int    `json:"port"`
}

var config Config

func init() {
	// Open config file
	rawConfig, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(rawConfig, &config)
	if err != nil {
		log.Fatalln(err)
	}

	// read jsonFile as byte
}
func GetMatches() {
	fmt.Println(config)
}
