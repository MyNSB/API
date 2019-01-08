package Week

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"time"
	"Util"
	"github.com/metakeule/fmtdate"
	"github.com/tidwall/gjson"
	"github.com/Azure/go-ntlmssp"
	"errors"
)

// TODO: Find a Go library that will easily allow me to parse html

func GetWeek(w http.ResponseWriter, r* http.Request, _ httprouter.Params) {

	// Determine the week type when term starts
	startWeekType, termStart := getStartWeekType()
	today := time.Now()

	// Calculate difference between two dates in terms of weeks
	diff := today.Sub(termStart)
	weeksDif := int((diff.Hours() / 24 ) / 7)

	// Determine the week type based on the weeksDiff
	// Can be a lot more efficient but tbh i am really cbbs
	if weeksDif % 2 == 1 && startWeekType == "A"{
		startWeekType = "B"
	} else if weeksDif % 2 == 1 && startWeekType == "B" {
		startWeekType = "A"
	}

	// Return our result
	Util.SolidError(200, "OK", startWeekType, "Week", w)
}





/*
	UTIL FUNCTIONS ===========================
 */

/* getStartWeek returns the type of week that the first week of term was
		@params;
			nil
 */
func getStartWeekType() (string, time.Time) {
	termDates, _ := getTermDates()

	var week string
	var termStart time.Time

	data := gjson.Parse(termDates)
	for _, name := range data.Array() {
		termStartRaw, _ := fmtdate.Parse("YYYY-MM-DD", name.Get("start_date").String())
		termEnd, _ := fmtdate.Parse("YYYY-MM-DD", name.Get("end_date").String())

		if time.Now().Before(termEnd) && time.Now().After(termStart) {
			week = name.Get("week_ab").String()
			termStart = termStartRaw
			break
		}
	}

	if week == "" {
		week = "A"
	}

	return week, termStart
}


/* getTermDates returns the term dates for that year
		@params;
			nil
 */
func getTermDates() (string, error) {

	// Set up client
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}

	req, _ := http.NewRequest("GET", "https://web3.northsydbo-h.schools.nsw.edu.au/classery/public/api/export/calendar", nil)
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
	value := gjson.Get(string(bytes), "term_dates")


	return value.String(), nil
}
/*
	END UTIL FUNCTIONS ===========================
 */
