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

// Handler for serving timetables
/*
	http handlers require minimal documentation
*/
func GetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Determine if the currently logged in student is allowed here
	allowed, _ := sessions.UserIsAllowed(r, w, "student")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	params, err := getParams(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "Assembly is not boolean")
		return
	}

	// Return the error
	util.Error(200, "OK", getTimes(params["term"].(string), params["day"].(string), params["assembly"].(bool)), "Response", w)
}

/*
	@ UTIL FUNCTIONS ===========================================
*/
/*
	getTimes returns the times given specific parameters
	@params;
		term string
		day string
		assembly bool
*/
func getTimes(term string, day string, assembly bool) string {
	// Load up the hash map
	var json []byte

	// Determine what to return
	timetable := Times

	// Determine if the thursday bell times should be changed
	if term == "2" || term == "3" {
		// Convert to non crawford shield times
		timetable["Thursday"]["Lunch"] = "12:38pm - 1:17pm"
	}

	// Determine if they want assembly
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

/*
	getParams returns the parameters of the incoming request
	@params;
		r *http.Request
*/
func getParams(r *http.Request) (map[string]interface{}, error) {
	term := r.URL.Query().Get("Term")
	day := r.URL.Query().Get("Day")
	assembly := r.URL.Query().Get("Assembly")

	assemblyBool := false

	// Convert to bool
	if util.Isset(assembly) {
		var err error
		assemblyBool, err = strconv.ParseBool(assembly)
		if err != nil {
			return nil, errors.New("assembly is not boolean")
		}
	}

	// Construct a map of results
	toReturn := make(map[string]interface{})
	toReturn["term"] = term
	toReturn["day"] = day
	toReturn["assembly"] = assemblyBool

	return toReturn, nil
}

/*
	@ END UTIL FUNCTIONS ===========================================
*/
