package belltimes

import (
	json2 "encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
	"strconv"
)


// RETRIEVAL FUNCTIONS

// getBellTimes takes a set of parameters and returns the belltimes in accordance to those params
func getBellTimes(term string, day string, assembly bool) string {

	json := []byte{}
	// Pull the data from the file
	timetable := timetableData


	if term == "2" || term == "3" {
		// Convert to non crawford shield timetableData
		timetable["Thursday"]["Lunch"] = "12:38pm - 1:17pm"
	}
	if !assembly {
		// Switch monday with friday
		val := timetable["Friday"]
		timetable["Monday"] = val
	}

	// Determine what to return
	if day == "" {
		json, _ = json2.Marshal(timetable)
	} else {
		table := timetable[day]
		json, _ = json2.Marshal(table)
	}

	return string(json)
}












// UTILITY FUNCTIONS

// getParams takes a http request and returns the user parameters from it
func getParams(r *http.Request) (map[string]interface{}, error) {

	term := r.URL.Query().Get("Term")
	day := r.URL.Query().Get("Day")
	assemblyRaw := r.URL.Query().Get("Assembly")

	assembly := false

	// Convert assembly to a boolean
	if util.NonNull(assemblyRaw) {
		var err error
		assembly, err = strconv.ParseBool(assemblyRaw)
		if err != nil {
			return nil, errors.New("assembly is not boolean")
		}
	}


	return map[string]interface{}{
		"term": term,
		"day": day,
		"assembly": assembly,
	}, nil
}












// HTTP HANDLERS

// GetHandler is a http handler that handles requests for the belltimes
func GetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	allowed, _ := sessions.IsUserAllowed(r, w, "user")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	params, err := getParams(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "Assembly is not boolean")
		return
	}

	belltimes := getBellTimes(
		params["term"].(string),
		params["day"].(string),
		params["assembly"].(bool))

	// Return the error
	util.Error(200, "OK", belltimes, "Response", w)
}