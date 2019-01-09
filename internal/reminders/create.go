package reminders

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/metakeule/fmtdate"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/student"
	"mynsb-api/internal/util"
	"net/http"
	"strconv"
	"time"
)

func CreateReminder(db *sql.DB, user student.User, reminder Reminder) error {
	// Convert the student ID into an integer
	studentID, _ := strconv.Atoi(user.Name)
	jsonTXT, _ := json.Marshal(reminder.Tags)
	headersTXT, _ := json.Marshal(reminder.Headers)

	// Push into database
	_, err := db.Exec("INSERT INTO reminders (student_id, headers, body, tags, reminder_date_time)  VALUES ($1, $2, $3, $4, $5::TIMESTAMP)",
		studentID, headersTXT, reminder.Body, jsonTXT, reminder.ReminderDateTime)

	if err != nil {
		return errors.New("oopsie, doopsie, doo")
	}

	return nil
}

// Create reminders handler
func CreateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	user, _ := sessions.ParseSessions(r, w)

	r.ParseForm()

	// Login into database
	db.Conn("admin")

	// Close the connection at the end
	defer db.DB.Close()

	// Get the post vars
	body := r.PostFormValue("Body")
	subject := r.PostFormValue("Subject")
	reminderDateTime := r.PostFormValue("Reminder_Date_Time")
	tagsTXT := r.PostFormValue("Tags")

	// Get the required fields
	if util.CompoundIsset(body, reminderDateTime, tagsTXT, subject) {
		// Decode the tags
		var tags []string
		err := json.Unmarshal([]byte(tagsTXT), &tags)
		if err != nil {
			quickerrors.MalformedRequest(w, "Invalid Tags sent, must be in JSON format")
			return
		}

		// Start creating the headers
		var headers = make(map[string]interface{})
		headers["Content-Length"] = len(body)
		headers["Tags-Length"] = len(tags)
		headers["Created-At"] = time.Now().String()
		headers["Subject"] = subject

		// Parse the given date time
		reminderDateTimeVal, err := fmtdate.Parse("DD-MM-YYYY hh:mm", reminderDateTime)
		if err != nil {
			quickerrors.MalformedRequest(w, "Dates are invalid, must follow the following format: DD-MM-YYYY hh:mm")
			return
		}

		// Push everything into a reminders type
		reminder := Reminder{
			Headers:          headers,
			Body:             body,
			Tags:             tags,
			ReminderDateTime: reminderDateTimeVal,
		}

		// Push everything into the database
		err = CreateReminder(db.DB, user, reminder)
		if err != nil {
			quickerrors.MalformedRequest(w, "Unable to create event, this could be because the date time your provided is in the past")
			return
		}

		quickerrors.OK(w)
	} else {
		quickerrors.MalformedRequest(w, "Missing Parameters, check the API Documentation")
		return
	}
}
