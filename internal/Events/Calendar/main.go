package Calendar

import (
	"Util"
	"errors"
	"fmt"
	"github.com/Azure/go-ntlmssp"
	"github.com/julienschmidt/httprouter"
	"github.com/metakeule/fmtdate"
	"io/ioutil"
	"net/http"
	"QuickErrors"
	"Sessions"
	"github.com/tidwall/gjson"
)



// Function to get all data from the
func GetAll() (string, int, error) {
	content, _ := sendRequest("https://web3.northsydbo-h.schools.nsw.edu.au/classery/public/api/export/calendar")
	// Return it
	return string(content), 0, nil
}




func GetBetween(startDate, endDate string) (string, int, error) {
	// Create times
	// Create times

	// Convert the two strings to dates
	start, err := fmtdate.Parse("DD-MM-YYYY", startDate)
	if err != nil {
		return "", 0, errors.New("invalid Date format, must fit format: DD-MM-YYYY")
	}
	end, err := fmtdate.Parse("DD-MM-YYYY", endDate)
	if err != nil {
		return "", 0, errors.New("invalid Date format, must fit format: DD-MM-YYYY")
	}


	// This is the format that the school's API accepts
	apiFormat := "YYYY-MM-DD"

	// Send the http request
	url := fmt.Sprintf("http://web3.northsydbo-h.schools.nsw.edu.au/classery/public/api/export/calendar?start=%s&end=%s", fmtdate.Format(apiFormat, start), fmtdate.Format(apiFormat,end))
	bytes, err := sendRequest(url)


	return string(bytes), 200, nil

}




// Function send a reuqest with the details provided to us by the school
func sendRequest(url string) (string, error) {

	// Set up client
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}

	req, _ := http.NewRequest("GET", url, nil)
	// Set up the basic auth headers
	req.SetBasicAuth("skedular", "chickenfarm")
	req.Header.Set("X-AUTH", "!te5D?DI<c0#t=2nZir0_eC4.(`i1>p/xEj[Qk_v10dF|G~*{zvwcwTw+`MS&o)M")

	// Perform request
	res, err := client.Do(req)
	if err != nil {
		return "", errors.New("something went wrong when trying to retrieve calendar")
	}

	// Attain the resutls
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)


	// Parse data as json
	value := gjson.Get(string(bytes), "events")


	return value.String(), nil
}





// Http handler for calendar thingy ma bop
func GetCalendar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Get the currently logged in user

	user, err := Sessions.ParseSessions(r, w)
	if err != nil || !Util.ExistsString(user.Permissions, "student") {
		QuickErrors.NotEnoughPrivledges(w)
		return
	}


	// Get the get variables from the query
	dateStart := r.URL.Query().Get("Start")
	dateEnd := r.URL.Query().Get("End")
	fmt.Printf("Start: %s, End: %s", dateStart, dateEnd)
	// If these variables do not exist just spew everything
	if Util.CompoundIsset(dateStart, dateEnd) {
		resp, status, err := GetBetween(dateStart, dateEnd)

		if err != nil && status == 0 {
			QuickErrors.InteralServerError(w)
			return
		} else if status != 200 {
			Util.Error(status, "Something went horribly wrong", "This could be because a failed attempt to authenticate with the school servers, please try again later", "Something went horribly wrong", w)
			return
		} else {
			Util.Error(200, "OK", resp, "Calendar", w)
			return
		}
	} else {
		// Else spew everything
		// Get everything first
		resp, status, err := GetAll()
		if err != nil {
			QuickErrors.InteralServerError(w)
		} else {
			Util.Error(status, "Here: ", resp, "Here: ", w)
		}
	}
}
