package userdetails

import (
	"database/sql"
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/util"
	"mynsb-api/internal/db"
	"mynsb-api/internal/student"
)

// Don't need this but kinda ceebs fixing it so please deal with it later
type details struct {
	StudentID 	uint64
	Fname 		string
	Lname 		string
	Year  		uint8
}

func getDetails(db *sql.DB, user student.User) string {
	// Convert the student's name into an integer
	studentID, _ := strconv.Atoi(user.Name)

	rows, _ := db.Query("SELECT * FROM students WHERE student_id = $1", studentID)


	details := details{}

	// Iterate over it
	for rows.Next() {
		var StudentID uint64
		var Fname string
		var Lname string
		var Year  uint8

		rows.Scan(&StudentID, &Fname, &Lname, &Year)

		details.Year = Year
		details.StudentID = StudentID
		details.Lname = Lname
		details.Fname = Fname
	}


	// Marshall the response
	resp, _ := json.Marshal(details)

	return string(resp)
}



func GetDetailsHandler(w http.ResponseWriter, r* http.Request, _ httprouter.Params) {


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