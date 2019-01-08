package Reminders

import (
	"Sessions"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"Util"
	"User"
	"database/sql"
	"DB"
	"encoding/json"
	"time"
	"github.com/metakeule/fmtdate"
	"QuickErrors"
	"errors"
	"fmt"
)

func CreateReminder(db *sql.DB, user User.User, reminder Reminder) error {
	// Convert the student ID into an integer
	studentID, _ := strconv.Atoi(user.Name)
	jsonTXT, _ := json.Marshal(reminder.Tags)
	headersTXT, _ := json.Marshal(reminder.Headers)

	// Push into database
	_, err := db.Exec("INSERT INTO reminders (student_id, headers, body, tags, reminder_date_time)  VALUES ($1, $2, $3, $4, $5::TIMESTAMP)",
		studentID, headersTXT, reminder.Body, jsonTXT, reminder.ReminderDateTime)

	if err != nil {
		return errors.New("Oopsie, Doopsie, Doo")
	}

	return nil
}



// Create reminders handler
func CreateReminderHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	user, _ := Sessions.ParseSessions(r, w)

	r.ParseForm()


	// Login into database
	Util.Conn("sensitive", "database", "admin")

	// Close the connection at the end
	defer DB.DB.Close()


	// Get the post vars
	body := r.PostFormValue("Body")
	subject := r.PostFormValue("Subject")
	reminderDateTime := r.PostFormValue("Reminder_Date_Time")
	tagsTXT := r.PostFormValue("Tags")

	fmt.Printf("Form-Data", r.PostForm.Encode())

	// Get the required fields
	if Util.CompoundIsset(body, reminderDateTime, tagsTXT, subject) {
		// Decode the tags
		var tags []string
		err := json.Unmarshal([]byte(tagsTXT), &tags)
		if err != nil {
			QuickErrors.MalformedRequest(w, "Invalid Tags sent, must be in JSON format")
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
			QuickErrors.MalformedRequest(w, "Dates are invalid, must follow the following format: DD-MM-YYYY hh:mm")
			return
		}


		// Push everything into a reminders type
		reminder := Reminder{
			Headers: headers,
			Body: body,
			Tags: tags,
			ReminderDateTime: reminderDateTimeVal,
		}


		// Push everything into the database
		err = CreateReminder(DB.DB, user, reminder)
		if err != nil {
			QuickErrors.MalformedRequest(w, "Unable to create event, this could be because the date time your provided is in the past")
			return
		}

		QuickErrors.OK(w)
	} else {
		QuickErrors.MalformedRequest(w, "Missing Parameters, check the API Documentation")
		return
	}
}
