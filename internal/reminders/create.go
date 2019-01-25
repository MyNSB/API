package reminders

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
	"strconv"
	"time"
)


// CREATION FUNCTIONS

// insertIntoDB takes a reminder and inserts it into the database for us
func (reminder Reminder) insertIntoDB(db *sql.DB, userID string) error {

	// Convert the user ID into an integer, convert the headers and tags into json
	studentID, _ := strconv.Atoi(userID)
	tagsJSON, _ := json.Marshal(reminder.Tags)
	headersJSON, _ := json.Marshal(reminder.Headers)

	// Push into database
	_, err := db.Exec("INSERT INTO reminders (student_id, headers, body, tags, reminder_date_time)  VALUES ($1, $2, $3, $4, $5::TIMESTAMP)",
		studentID, headersJSON, reminder.Body, tagsJSON, reminder.DateTime)
	if err != nil {
		panic(err)
		return errors.New("oopsie, doopsie, doo")
	}

	return nil
}












// UTILITY FUNCTIONS

// parseReminder attains the incoming reminder within a HTTP request and returns it as a reminder object
func parseReminder(r *http.Request) (Reminder, error) {

	r.ParseForm()

	body := r.FormValue("Body")
	subject := r.FormValue("Subject")
	reminderDateTimeRAW := r.FormValue("Date_Time")
	tagsTXT := r.FormValue("Tags")

	// Determine if the request is valid
	if !util.IsSet(body, subject, reminderDateTimeRAW, tagsTXT) {
		return Reminder{}, errors.New("invalid request")
	}

	// Parse the time into a suitable format
	reminderDateTime, _ := util.ParseDateTime(reminderDateTimeRAW)


	// Parse the json into an actual structure
	var tags []string
	err := json.Unmarshal([]byte(tagsTXT), &tags)
	if err != nil {
		return Reminder{}, err
	}

	// Generate a map of the headers associated with the reminder
	headers := map[string]interface{} {
		"Content-Length": len(body),
		"Tags-Length": len(tags),
		"Created-At": time.Now().In(util.TIMEZONE).String(),
		"Subject": subject,
	}

	return Reminder{
		Headers:  headers,
		Body:     body,
		Tags:     tags,
		DateTime: reminderDateTime,
	}, nil
}











// HTTP HANDLERS

// Create reminders handler
func CreationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	allowed, user := sessions.IsUserAllowed(r, w, "user")
	if !allowed {
		quickerrors.NotLoggedIn(w)
		return
	}


	// Login into database
	db.Conn("admin")
	defer db.DB.Close()

	requestedReminder, err := parseReminder(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "Missing or invalid parameters, check the API Documentation")
		return
	}

	err = requestedReminder.insertIntoDB(db.DB, user.Name)
	if err != nil {
		quickerrors.InternalServerError(w)
		return
	}
}
