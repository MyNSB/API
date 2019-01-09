package userdetails

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/student"
	"mynsb-api/internal/util"
	"net/http"
	"strconv"
)

// Don't need this but kinda ceebs fixing it so please deal with it later
type details struct {
	StudentID uint64
	Fname     string
	Lname     string
	Year      uint8
}

func getDetails(db *sql.DB, user student.User) string {
	// Convert the student's name into an integer
	studentID, _ := strconv.Atoi(user.Name)

	rows, _ := db.Query("SELECT * FROM students WHERE student_id = $1", studentID)

	details := details{}

	// Iterate over it
	for rows.Next() {
		rows.Scan(&details.StudentID, &details.Fname, &details.Lname, &details.Year)
	}

	// Marshall the response
	resp, _ := json.Marshal(details)

	return string(resp)
}

func GetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	user, err := sessions.ParseSessions(r, w)

	db.Conn("student")

	defer db.DB.Close()

	// Determine if the sessions is an actual student
	if err != nil || !util.ExistsString(user.Permissions, "student") {
		quickerrors.NotLoggedIn(w)
		return
	}

	util.Error(200, "OK", getDetails(db.DB, user), "Response", w)
	return
}
