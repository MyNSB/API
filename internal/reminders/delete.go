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


// DELETION FUNCTIONS

// deleteReminder takes a reminderID and deletes the reminder associated with that ID, the studentID is added for security
// in case someone tries to delete a reminder that is not theirs
func deleteReminder(db *sql.DB, reminderIdRAW, studentIDRAW string) error {

	// Convert our ids to integers
	studentID, _ := strconv.Atoi(studentIDRAW)
	reminderID, err := strconv.Atoi(reminderIdRAW)
	if err != nil {
		errors.New("reminderID is not an integer")
	}

	_, err = db.Exec("DELETE FROM reminders WHERE reminder_id = $1 AND student_id = $2", reminderID, studentID)
	if err != nil {
		return errors.New("user does not own this reminder")
	}

	return nil
}










// HTTP HANDLERS

// DeletionHandler is a HTTP handler that handles the deletion of a requested user event
func DeletionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	db.Conn("admin")
	defer db.DB.Close()

	allowed, user := sessions.IsUserAllowed(r, w, "user")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}


	r.ParseForm()
	reminderID := r.Form.Get("Reminder_ID")
	if util.NonNull(reminderID) {
		quickerrors.MalformedRequest(w, "Missing parameters, please refer to the API documentation")
		return
	}



	err := deleteReminder(db.DB, reminderID, user.Name)
	if err != nil {
		quickerrors.MalformedRequest(w, "User does not own requested reminder or reminder ID is not an integer")
		return
	}
	quickerrors.OK(w)

}
