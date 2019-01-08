package User

import (
	"database/sql"
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"Sessions"
	"QuickErrors"
	"Util"
	"DB"
)

// Don't need this but kinda ceebs fixing it so please deal with it later
type details struct {
	StudentID 	uint64
	Fname 		string
	Lname 		string
	Year  		uint8
}

func getDetails(db *sql.DB, user User) string {
	// Convert the user's name into an integer
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



func GetDeatilsHandler(w http.ResponseWriter, r* http.Request, _ httprouter.Params) {


	user, err := Sessions.ParseSessions(r, w)

	Util.Conn("sensitive", "database", "student")

	defer DB.DB.Close()

	// Determine if the sessions is an actual user
	if err != nil || !Util.ExistsString(user.Permissions, "student") {
		QuickErrors.NotLoggedIn(w)
		return
	}


	Util.Error(200, "OK", getDetails(DB.DB, user), "Response", w)
	return
}