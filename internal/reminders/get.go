package reminders

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/metakeule/fmtdate"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/student"
	"mynsb-api/internal/util"
	"net/http"
	"time"
)

func getReminders(db *sql.DB, start time.Time, end time.Time, user student.User) []Reminder {
	res, err := db.Query("SELECT * FROM reminders WHERE reminder_date_time BETWEEN $1::TIMESTAMP AND $2::TIMESTAMP AND student_id = $3 ORDER BY reminder_date_time ASC",
		start, end, user.Name)

	if err != nil {
		panic(err)
	}

	var container []Reminder

	for res.Next() {
		var headers []byte
		var studentID int
		var tags []byte

		var reminder Reminder

		// Scan into the containers
		res.Scan(&reminder.ReminderId, &studentID, headers, &reminder.Body, &tags, &reminder.ReminderDateTime)

		// Start converting it into the correct types
		// Headers
		json.Unmarshal(headers, &reminder.Headers)
		// Tags
		json.Unmarshal(tags, &reminder.Tags)
		// Push into array
		container = append(container, reminder)
	}

	res.Close()
	return container
}

// Get reminders handler
func GetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	user, err := sessions.ParseSessions(r, w)

	if err != nil {
		quickerrors.NotLoggedIn(w)
		return
	}

	db.Conn("student")

	// Close that database at the end
	defer db.DB.Close()

	if ps.ByName("reqType") == "/Today" {
		reminders, _ := json.Marshal(getReminders(db.DB, time.Now().Add(time.Hour*-24), time.Now().Add(time.Hour*24), user))
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
