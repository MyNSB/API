package events

import (
	"database/sql"
	_ "database/sql"
	json2 "encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/metakeule/fmtdate"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/util"
	"net/http"
	"strconv"
	time2 "time"
)



// RETRIEVAL FUNCTIONS

// getWithID takes an eventID and a db and queries the database for an event with the matching event ID
func getWithID(eventID int, db *sql.DB) (Event, error) {
	eventArr, err := performRequest(db, "SELECT * FROM events WHERE event_id = $1", eventID)
	return eventArr[0], err
}


// getBetween returns all events between two times
func getBetween(db *sql.DB, startStr string, endStr string) ([]Event, error) {

	start, parseErrorOne := util.ParseDate(startStr)
	end, parseErrorTwo := util.ParseDate(endStr)

	if parseErrorOne != nil || parseErrorTwo != nil {
		return []Event{}, errors.New("invalid date")
	}

	return performRequest(db, "SELECT * FROM events WHERE event_start BETWEEN $1::TIMESTAMP AND $2::TIMESTAMP", start, end)
}


// getAll returns all events currently in the database
func getAll(db *sql.DB) ([]Event, error) {
	return performRequest(db, "SELECT * FROM events")
}












// UTILITY FUNCTIONS

// performRequest takes a simple request and a string of parameters, executes the query and returns the result
func performRequest(db *sql.DB, query string, args ...interface{}) ([]Event, error) {

	rows, err := db.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var allEvents []Event

	for rows.Next() {
		// Parse data into an event object and append it to the allEvents structure
		event := Event{}
		event.ScanFrom(rows)
		allEvents = append(allEvents, event)
	}

	return allEvents, nil
}


// parseRequest takes an http request and determines what type of request they are performing
func parseRequest(r *http.Request) (string, map[string]string) {

	eventID := r.URL.Query().Get("Event_ID")
	start := r.URL.Query().Get("Start")
	end := r.URL.Query().Get("End")


	// Determine what type of request the user wants
	requestType := ""
	switch {
	case util.IsSet(start, end):
		requestType = "Range"
		break
	case util.IsSet(eventID):
		requestType = "getWithID"
		break
	default:
		requestType = "getAll"
	}


	return requestType, map[string]string{
		"eventID": eventID,
		"start": start,
		"end": end,
	}
}












// HTTP HANDLERS

// EventRetrievalHandler returns all events according to a http request
func EventRetrievalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	db.Conn("student")
	defer db.DB.Close()

	requestType, params := parseRequest(r)
	var response []byte

	// Determine which function is being requested
	switch requestType {
	case "getAll":
		json, _ := getAll(db.DB)
		response, _ = json2.Marshal(json)
		break

	case "getWithID":
		// convert the eventID to an integer
		eventId, _ := strconv.Atoi(params["eventID"])
		json, _ := getWithID(eventId, db.DB)
		response, _ = json2.Marshal(json)
		break

	case "Range":
		json, _ := getBetween(db.DB, params["start"], params["end"])
		response, _ = json2.Marshal(json)
		break

	default:
		quickerrors.NotFound(w)
		return
	}

	util.HTTPResponse(200, "OK", string(response), "Result: ", w)
}