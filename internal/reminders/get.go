package reminders

import (
	"time"
	"mynsb-api/internal/student"
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"mynsb-api/internal/db"
	"github.com/metakeule/fmtdate"
	"mynsb-api/internal/quickerrors"
)

func getReminders(db *sql.DB, start time.Time, end time.Time, user student.User) []Reminder {
	res, err := db.Query("SELECT * FROM reminders WHERE reminder_date_time BETWEEN $1::TIMESTAMP AND $2::TIMESTAMP AND student_id = $3 ORDER BY reminder_date_time ASC",
		start, end, user.Name)

	if err != nil {
		panic(err)
	}

	var container []Reminder

	for res.Next() {
		var reminderID int
		var headers []byte
		var studentID int
		var body string
		var tags []byte
		var reminderDateTime time.Time

		// Scan into the containers
		res.Scan(&reminderID, &studentID, &headers, &body, &tags, &reminderDateTime)

		// Start converting it into the correct types
		// Headers
		var headersContainer map[string]interface{}
		json.Unmarshal(headers, &headersContainer)

		// Tags
		var tagsContainer []string
		json.Unmarshal(tags, &tagsContainer)
		// Push into array
		container = append(container, Reminder{ReminderId: reminderID, Headers: headersContainer, Body: body, Tags: tagsContainer, ReminderDateTime: reminderDateTime})
	}

	res.Close()
	return container
}

// Get reminders handler
func GetRemindersHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	user, err := sessions.ParseSessions(r, w)

	if err != nil {
		quickerrors.NotLoggedIn(w)
		return
	}

	db.Conn("student")

	// Close that database at the end
	defer db.DB.Close()

	if ps.ByName("reqType") == "/Today" {
		reminders, _ := json.Marshal(getReminders(db.DB, time.Now().Add(time.Hour * -24), time.Now().Add(time.Hour*24), user))
		util.Error(200, "OK", string(reminders), "Response", w)
		return
	} else {
		// Get the required fields
		startTime := r.URL.Query().Get("Start_Time")
		endTime := r.URL.Query().Get("End_Time")

		if util.CompoundIsset(startTime, endTime) {

			// Start converting the dates to the correct format
			startTimeDate, err1 := fmtdate.Parse("DD-MM-YYYY", startTime)
			endTimeDate, err2 := fmtdate.Parse("DD-MM-YYYY", endTime)
			if err1 != nil || err2 != nil {
				quickerrors.MalformedRequest(w, "Dates are invalid, must follow the following format: DD-MM-YYYY hh:mm")
			}

			reminders, _ := json.Marshal(getReminders(db.DB, startTimeDate, endTimeDate, user))
			util.Error(200, "OK", string(reminders), "Response", w)

		} else {
			quickerrors.MalformedRequest(w, "Missing parameters, check the API documentation")
			return
		}
	}

}
