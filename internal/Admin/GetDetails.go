package Admin

import (
	"User"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"Util"
	"QuickErrors"
	json2 "encoding/json"
	"DB"
)

type details struct {
	AdminName 		string
	Persmissions 	string
}


// Retrieves the details for an incoming user
/*
	http handlers need minimal documentation
 */
func GetDetailsHandler (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Determine if the user is allowed here
	allowed, user := Util.UserIsAllowed(r, w, "admin")
	if !allowed {
		QuickErrors.NotEnoughPrivledges(w)
		return
	}

	// Connect to Database
	Util.Conn("sensitive", "database", "student")
	defer DB.DB.Close()

	// Get the details
	Util.Error(200, "OK", getDetails(user), "Response", w)
	return
}




/*
	@ UTIL FUNCTIONS ==================================================
 */
 /*
 	getDetails takes a user and returns the details for that user
 	@params;
 		user User.User
  */
func getDetails(user User.User) string {
	rows, _ := DB.DB.Query("SELECT admin_name, admin_permissions FROM admins WHERE admin_name = $1", user.Name)
	defer rows.Close()

	// Details structure
	details := details{}

	// Scan into the rows
	for rows.Next() {
		rows.Scan(&details.AdminName, &details.Persmissions)
	}

	json, _ := json2.Marshal(details)

	return string(json)
}
/*
	@ END UTIL FUNCTIONS ==================================================
 */