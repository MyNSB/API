package reminders

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
	"time"
	"errors"
)



// RETRIEVAL FUNCTIONS

// getReminders returns all reminders from a specific user between two provided dates
func getReminders(db *sql.DB, start time.Time, end time.Time, studentID string) []Reminder {
	res, _ := db.Query("SELECT * FROM reminders WHERE reminder_date_time BETWEEN $1::TIMESTAMP AND $2::TIMESTAMP AND student_id = $3 ORDER BY reminder_date_time ASC",
		start, end, studentID)
	defer res.Close()


	// Response array
	var response []Reminder


	for res.Next() {
		// Core data
		var headers []byte
		var studentID int
		var tags []byte

		reminder := Reminder{}

		res.Scan(&reminder.ID, &studentID, &headers, &reminder.Body, &tags, &reminder.DateTime)

		json.Unmarshal(headers, &reminder.Headers)
		json.Unmarshal(tags, &reminder.Tags)
		response = append(response, reminder)
	}

	return response
}









// UTILITY FUNCTIONS

// parseParams takes the rqeust and parses and reads the start and end times
func parseParams(r *http.Request) (map[string]time.Time, error) {

	startTime := r.URL.Query().Get("Start_Time")
	endTime := r.URL.Query().Get("End_Time")

	start, parseErrorOne := util.ParseDateTime(startTime)
	end, parseErrorTwo := util.ParseDateTime(endTime)

	if parseErrorOne != nil || parseErrorTwo != nil {
		return nil, errors.New("could not parse datetimes")
	}


	return map[string]time.Time{
		"start": start,
		"end": end,
	}, nil
}









// HTTP HANDLERS

// RetrievalHandler takes a user's request for their reminders and returns all reminders that correspond to their request
func RetrievalHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	currUser, err := sessions.ParseSessions(r, w)
	if err != nil {
		quickerrors.NotLoggedIn(w)
		return
	}

	db.Conn("student")
	defer db.DB.Close()

	// Get all reminders for today
	if ps.ByName("reqType") == "/Today" {
		yesterday := time.Now().Add(time.Hour*-24)
		tommorrow := time.Now().Add(time.Hour*24)

		reminders, _ := json.Marshal(getReminders(db.DB, yesterday, tommorrow, currUser.Name))
		util.HTTPResponse(200, "OK", string(reminders), "Response", w)
		return
	}

	// Looks like they have a start and an end time
	params, notValid := parseParams(r)
	if  notValid != nil {
		reminders, _ := json.Marshal(getReminders(db.DB, params["start"], params["end"], currUser.Name))
		util.HTTPResponse(200, "OK", string(reminders), "Response", w)
		return
	}



	quickerrors.MalformedRequest(w, "could not determine request type")
}
