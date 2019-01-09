package belltimes

import (
	"encoding/json"
	"io/ioutil"
	"mynsb-api/internal/util"
)

var Times = make(map[string]map[string]string)

func Init() {
	// Get the GOPATH
	gopath := util.GetGOPATH()

	// Read the data
	data, _ := ioutil.ReadFile(gopath + "/src/mynsb-api/internal/belltimes/bellTimes.json")

	// Load the json data into the times map
	json.Unmarshal(data, &Times)
}
