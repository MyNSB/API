package reminders

import (
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
	"strconv"
)

func deleteEvent(db *sql.DB, reminderId, studentID int) error {
	_, err := db.Exec("DELETE FROM reminders WHERE reminder_id = $1 AND student_id = $2", reminderId, studentID)
	if err != nil {
		return errors.New("user does not own this reminder")
	}

	return nil
}

func DeleteHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	r.ParseForm()

	db.Conn("admin")
	defer db.DB.Close()

	user, err := sessions.ParseSessions(r, w)
	if err != nil || !util.ExistsString(user.Permissions, "user") {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// Get the post vars
	reminderIDTXT := r.Form.Get("Reminder_ID")

	// Determine if the reminder id is really set or not
	if util.NonNull(reminderIDTXT) {

		// Begin the conversion from text to int for all the ids
		studentID, _ := strconv.Atoi(user.Name)
		reminderID, err := strconv.Atoi(reminderIDTXT)
		if err != nil {
			quickerrors.MalformedRequest(w, "Reminder ID is not an integer")
			return
		}

		deleteEvent(db.DB, reminderID, studentID)
		quickerrors.OK(w)

	} else {
		quickerrors.MalformedRequest(w, "Missing parameters, please refer to the API documentation")
		return
	}
}
