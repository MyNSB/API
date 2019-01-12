package calendar

import (
	"errors"
	"fmt"
	"github.com/Azure/go-ntlmssp"
	"github.com/julienschmidt/httprouter"
	"github.com/metakeule/fmtdate"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
	"time"
)



// RETRIEVAL FUNCTIONS

// getAll returns all events currently stored on the school's calendar
func getAll() (string, int) {
	content, statusCode := sendRequest("https://web3.northsydbo-h.schools.nsw.edu.au/classery/public/api/export/calendar")
	return string(content), statusCode
}

// getBetween returns all events betwen a start an an end date
func getBetween(startDate, endDate string) (string, int) {

	eventStart, parseErrorOne := parseDate(startDate)
	eventEnd, parseErrorTwo := parseDate(endDate)
	if parseErrorOne != nil || parseErrorTwo != nil {
		return "", 400
	}

	// This is the only date format the school accepts
	schoolDateFormat := "YYYY-MM-DD"

	// Create a request URL
	requestURL := fmt.Sprintf("http://web3.northsydbo-h.schools.nsw.edu.au/classery/public/api/export/calendar?start=%s&end=%s", fmtdate.Format(schoolDateFormat, eventStart), fmtdate.Format(schoolDateFormat, eventEnd))

	bytes, statusCode := sendRequest(requestURL)
	return string(bytes), statusCode
}












// UTILITY FUNCTIONS

// sendRequest takes a url and sends the HTTP request for us, it utilises an NTLM client
func sendRequest(url string) (string, int) {

	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("skedular", "chickenfarm")
	req.Header.Set("X-AUTH", "!te5D?DI<c0#t=2nZir0_eC4.(`i1>p/xEj[Qk_v10dF|G~*{zvwcwTw+`MS&o)M")
	res, err := client.Do(req)
	if err != nil {
		return "", 500
	}

	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)

	// Parse and read the request
	value := gjson.Get(string(bytes), "events")
	return value.String(), 200
}


// parseDate takes a single date and parses it into the date format that is univsersal within the API
func parseDate(date string) (time.Time, error) {
	return fmtdate.Parse("DD-MM-YYYY", date)
}












// HTTP HANDLERS

// CalendarRetrievalHandler takes a users request for a calendar and returns the corresponding event
func CalendarRetrievalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	allowed, _ := sessions.IsUserAllowed(r, w, "user")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
	}

	dateStart := r.URL.Query().Get("Start")
	dateEnd := r.URL.Query().Get("End")

	// Determine what type of request they want
	// getBetween
	if util.IsSet(dateStart, dateEnd) {
		resp, statusCode := getBetween(dateStart, dateEnd)

		switch statusCode {
		case 400:
			quickerrors.MalformedRequest(w, "Invalid parameters")
			return
			break
		case 500:
			quickerrors.InternalServerError(w)
			return
			break
		}

		util.Error(200, "OK", resp, "calendar", w)
		return

	}

	// getAll
	resp, statusCode := getAll()
	if statusCode != 200 {
		util.Error(statusCode, "Here: ", resp, "Here: ", w)
	}
}
