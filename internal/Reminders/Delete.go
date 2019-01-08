package Reminders

import (
	"database/sql"
	"errors"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"Sessions"
	"QuickErrors"
	"Util"
	"strconv"
	"DB"
)

func deleteEvent(db *sql.DB, reminderId, studentID int) error {
	_, err := db.Exec("DELETE FROM reminders WHERE reminder_id = $1 AND student_id = $2", reminderId, studentID)
	if err != nil {
		return errors.New("user does not own this reminder")
	}

	return nil
}



func DeleteReminderHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	r.ParseForm()


	Util.Conn("sensitive", "database", "admin")
	defer DB.DB.Close()


	user, err := Sessions.ParseSessions(r, w)
	if err != nil || !Util.ExistsString(user.Permissions, "student") {
		QuickErrors.NotEnoughPrivledges(w)
		return
	}

	// Get the post vars
	reminderIDTXT := r.Form.Get("Reminder_ID")

	// Determine if the reminder id is really set or not
	if Util.Isset(reminderIDTXT) {

		// Begin the conversion from text to int for all the ids
		studentID, _ := strconv.Atoi(user.Name)
		reminderID, err := strconv.Atoi(reminderIDTXT)
		if err != nil {
			QuickErrors.MalformedRequest(w, "Reminder ID is not an integer")
			return
		}

		deleteEvent(DB.DB, reminderID, studentID)
		QuickErrors.OK(w)

	} else {
		QuickErrors.MalformedRequest(w, "Missing parameters, please refer to the API documentation")
		return
	}
}
