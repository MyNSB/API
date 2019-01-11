package admin

import (
	json2 "encoding/json"
	"github.com/julienschmidt/httprouter"
	"mynsb-api/internal/db"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/user"
	"mynsb-api/internal/util"
	"net/http"
)

type details struct {
	AdminName   string
	Permissions string
}

// Retrieves the details for an incoming user
/*
	http handlers need minimal documentation
*/
func GetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Determine if the user is allowed here
	allowed, user := sessions.UserIsAllowed(r, w, "admin")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// Connect to Database
	db.Conn("user")
	defer db.DB.Close()

	// Get the details
	util.Error(200, "OK", getDetails(user), "Response", w)
	return
}

/*
	@ UTIL FUNCTIONS ==================================================
*/
/*
	getDetails takes a user and returns the details for that user
	@params;
		user user.user
*/
func getDetails(user user.User) string {
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
