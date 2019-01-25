package belltimes

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"mynsb-api/internal/util"
)

var timetableData = make(map[string]map[string]string)

// init function loads the data into the timetableData data structure
func init() {

	// Attain the location of the bellTimes data
	gopath := util.GetGOPATH()
	bellTimesDir := filepath.FromSlash(gopath + "/src/mynsb-api/internal/belltimes/bellTimes.json")

	// Read the file and unmarshal it into the data structure
	data, _ := ioutil.ReadFile(bellTimesDir)
	json.Unmarshal(data, &timetableData)
}
