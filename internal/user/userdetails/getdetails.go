package userdetails

import (
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/user"
	"mynsb-api/internal/util"
	"net/http"
	"strconv"
)


// details struct dictates the possible response from a query returned form the user table in the database
type details struct {
	StudentID uint64
	Fname     string
	Lname     string
	Year      uint8
}



// UTILITY FUNCTIONS

// getDetails takes a user object and returns all the details stored in the DB about that user
func getDetails(db *sql.DB, user user.User) string {

	// Convert the user's ID into an integer
	studentID, _ := strconv.Atoi(user.Name)
	rows, _ := db.Query("SELECT * FROM students WHERE student_id = $1", studentID)


	details := details{}
	for rows.Next() {
		rows.Scan(&details.StudentID, &details.Fname, &details.Lname, &details.Year)
	}
	resp, _ := json.Marshal(details)

	return string(resp)
}











// HTTP HANDLERS

// RetrievalHandler handles an incoming HTTP request
func RetrievalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	db.Conn("student")
	defer db.DB.Close()


	allowed, currUser := sessions.IsUserAllowed(r, w, "user")

	// Determine if generated user is legitimately a user
	if !allowed {
		quickerrors.NotLoggedIn(w)
		return
	}


	util.Error(200, "OK", getDetails(db.DB, currUser), "Response", w)
}
