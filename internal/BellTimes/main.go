package BellTimes

import (
	"io/ioutil"
	"encoding/json"
)

var Times = map[string]map[string]string{}

func Init() {
	// Read the data
	data, _ := ioutil.ReadFile("src/BellTimes/bellTimes.json")


	// Load the json data into the times map
	json.Unmarshal(data, &Times)
}
