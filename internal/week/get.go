package week

import (
	"errors"
	"github.com/Azure/go-ntlmssp"
	"github.com/julienschmidt/httprouter"
	"github.com/metakeule/fmtdate"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"mynsb-api/internal/util"
	"net/http"
	"time"
)




// RETRIEVAL FUNCTIONS

// getCurrTermInfo determines what type of week the current SCHOOL TERM starts on
func getCurrTermInfo() (string, time.Time) {

	termDates, _ := getTermDates()
	var week string
	var termStart time.Time
	termData := gjson.Parse(termDates)


	for _, name := range termData.Array() {
		termStartRaw, _	:= parseSchoolDate(name.Get("start_date").String())
		termEnd, _ 		:= parseSchoolDate(name.Get("end_date").String())

		if time.Now().In(util.TIMEZONE).Before(termEnd) && time.Now().In(util.TIMEZONE).After(termStart) {
			week = name.Get("week_ab").String()
			termStart = termStartRaw
			break
		}
	}

	// Just some cleaning up because the json response is kinda dodgy
	if week == "" {
		week = "A"
	}

	return week, termStart
}












// UTILITY FUNCTIONS

// getTermDates returns the dates for this current term based off the school's calendar
func getTermDates() (string, error) {

	// Set up NTLM client
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}


	req, _ := http.NewRequest("GET", "https://web3.northsydbo-h.schools.nsw.edu.au/classery/public/api/export/calendar", nil)
	req.SetBasicAuth("skedular", "chickenfarm")
	req.Header.Set("X-AUTH", "!te5D?DI<c0#t=2nZir0_eC4.(`i1>p/xEj[Qk_v10dF|G~*{zvwcwTw+`MS&o)M")


	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return "", errors.New("something went wrong when trying to retrieve calendar")
	}


	// Read data and extract term-dates via json
	bytes, _ := ioutil.ReadAll(res.Body)
	value := gjson.Get(string(bytes), "term_dates")


	// The result is a json extraction of the term dates
	return value.String(), nil
}


// parseSchoolDate parses a string and turns it into a date, note, this date is specifically in the school's date format
func parseSchoolDate(time string) (time.Time, error) {
	return fmtdate.Parse("YYYY-MM-DD", time)
}











// HTTP HANDLERS

// GetHandler takes a simple HTTP request for the current week and returns the current week
func GetHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {

	// Determine what week the term started on
	startWeekType, termStartDate := getCurrTermInfo()
	today := time.Now().In(util.TIMEZONE)

	// Calculate difference between two dates in terms of weeks
	diff := today.Sub(termStartDate)
	weeksDif := int((diff.Hours() / 24) / 7)

	// Determine the week type based on the weeksDiff
	if weeksDif % 2 == 1 {
		if startWeekType == "A" {
			startWeekType = "B"
		} else  if startWeekType == "B" {
			startWeekType = "A"
		}
	}


	util.HTTPResponseArr(200, "OK", startWeekType, "week", w)
}