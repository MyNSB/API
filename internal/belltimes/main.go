package belltimes

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"mynsb-api/internal/util"
)

var Times = make(map[string]map[string]string)

func init() {
	// Get the GOPATH
	gopath := util.GetGOPATH()
	// Set up the timetable
	bellTimesDir := filepath.FromSlash(gopath + "/mynsb-api/internal/belltimes/bellTimes.json")
	// Read the data
	data, _ := ioutil.ReadFile(bellTimesDir)

	// Load the json data into the times map
	json.Unmarshal(data, &Times)
}
