package Events

import (
	"DB"
	"Util"
	"database/sql"
	_ "database/sql"
	json2 "encoding/json"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"QuickErrors"
	"github.com/metakeule/fmtdate"
	"errors"
)



// HTTP handler for attaining all events
func GetEvents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Connect to the database
	Util.Conn("sensitive", "database", "student")
	// Close the request at the end
	defer DB.DB.Close()

	// Determine the request type and get the parameters
	requestType, params := determineRequestType(r)
	var jsonResp []byte

	// Determine which function is being request
	switch requestType{
	case "GetAll":
		// Perform that function
		json, _ := GetAll(DB.DB)
		// Json encode the result set
		jsonResp, _ = json2.Marshal(json)

	case "Get":
		// Get the event id
		eventId, _ := strconv.Atoi(params["eventID"])
		// Perform the function
		json, _ := Get(Event{EventID: int64(eventId)}, DB.DB)
		// Encode the response
		jsonResp, _ = json2.Marshal(json)
	case "Range":
		json, _ := GetBetween(DB.DB, params["start"], params["end"])
		// Encode the response
		jsonResp, _ = json2.Marshal(json)
	default:
		QuickErrors.NotFound(w)
		return
	}

	Util.Error(200, "OK", string(jsonResp), "Result: ", w)
}





// =======================================================
// Function to attain an event with given details
/*
	Get just returns an event given a particular event id
	@params;
		event Events.Event
		db *sql.DB
 */
func Get(event Event, db *sql.DB) (Event, error) {
	// Extract necessary details
	eventID := event.EventID
	eventArr, err := performRequest(db, "SELECT * FROM events WHERE event_id = $1", eventID)
	return eventArr[0], err
}


/*
	GetBetween returns all events between two given date strings
	@params;
		db *sql.DB
		start string
		end string
 */
func GetBetween(db *sql.DB, start string, end string)  ([]Event, error){
	// Begin converting our strings to actual dates
	Start, err := fmtdate.Parse("DD-MM-YYYY", start)
	if err != nil {
		return []Event{}, errors.New("invalid date")
	}

	End, err := fmtdate.Parse("DD-MM-YYYY", end)
	if err != nil {
		return []Event{}, errors.New("invalid date")
	}
	// Perform the actual request
	return performRequest(db, "SELECT * FROM events WHERE event_start BETWEEN $1::TIMESTAMP AND $2::TIMESTAMP", Start, End)
}


/*
	GetAll returns all currently active events
	@params;
		db *sql.DB
 */
func GetAll(db *sql.DB) ([]Event, error) {
	return performRequest(db, "SELECT * FROM events")
}
// ==========================================================





/*
	@ UTIL FUNCTIONS START ====================================
 */
/*
	performrequest takes a question and some arguments and returns an array of events corresponding to that query
	@params;
		db *sql.DB
		query string
		args ...interface{}
 */
 func performRequest(db *sql.DB, query string, args ...interface{}) ([]Event, error) {
	// Perform 	query
	rows, err := db.Query(query, args...)
	// Throw error if exits
	if err != nil {
		return nil, err
	}

	// Read from row
	// close it at the very end
	defer rows.Close()

	// Declare our events array
	var finEvent []Event

	// Read it
	for rows.Next() {
		event := Event{}
		// Scan into the event
		event.ScanFrom(rows)
		// Push into the events array
		finEvent = append(finEvent, event)
	}

	// Return that
	return finEvent, nil
}


/*
	determineRequestType determines the request type of the incoming request and returns all parameters related to that request
	@params;
		r *http.Request
 */
func determineRequestType(r *http.Request) (string, map[string]string) {
	typeReq := ""
	if r.URL.Query().Get("Event_ID") == "" {
		typeReq = "GetAll"
	} else if Util.CompoundIsset(r.URL.Query().Get("Start"), r.URL.Query().Get("End")) {
		typeReq = "Range"
	} else{
		typeReq = "Get"
	}

	// Prepare a map to return
	toReturn := make(map[string]string)
	toReturn["eventID"] = r.URL.Query().Get("Event_ID")
	toReturn["start"]   = r.URL.Query().Get("Start")
	toReturn["end"]     = r.URL.Query().Get("End")

	return typeReq, toReturn
}
/*
	@ UTIL FUNCTIONS END ====================================
 */