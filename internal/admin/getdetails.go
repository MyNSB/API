package admin

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/util"
	"mynsb-api/internal/quickerrors"
	json2 "encoding/json"
	"mynsb-api/internal/db"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/student"
)

type details struct {
	AdminName   string
	Permissions string
}

// Retrieves the details for an incoming student
/*
	http handlers need minimal documentation
 */
func GetDetailsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Determine if the student is allowed here
	allowed, user := sessions.UserIsAllowed(r, w, "admin")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// Connect to Database
	db.Conn("student")
	defer db.DB.Close()

	// Get the details
	util.Error(200, "OK", getDetails(user), "Response", w)
	return
}

/*
	@ UTIL FUNCTIONS ==================================================
 */
/*
	getDetails takes a student and returns the details for that student
	@params;
		student student.student
 */
func getDetails(user student.User) string {
	rows, _ := db.DB.Query("SELECT admin_name, admin_permissions FROM admins WHERE admin_name = $1", user.Name)
	defer rows.Close()

	// userDetails structure
	details := details{}

	// Scan into the rows
	for rows.Next() {
		rows.Scan(&details.AdminName, &details.Permissions)
	}

	json, _ := json2.Marshal(details)

	return string(json)
}

/*
	@ END UTIL FUNCTIONS ==================================================
 */
