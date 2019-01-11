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
	"database/sql"
)


// UTILITY FUNCTIONS

// getDetails takes a user and returns their all the information regarding that user, the assumption here is that the admin is a user
func getDetails(db *sql.DB, user user.User) string {

	// Query for the admin based off their name
	result, _ := db.Query("SELECT admin_name, admin_permissions FROM admins WHERE admin_name = $1", user.Name)
	defer result.Close()

	// Push it into a user
	details := user.User{}
	for result.Next() {
		result.Scan(&user.Name, &user.Permissions)
	}

	response, _ := json2.Marshal(details)
	return string(response)
}












// HTTP HANDLERS

// DetailRetrievalHandler retrieves the details of the admin that is currently logged in, if the requesting user is not an admin then they are blocked
func DetailRetrievalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Determine if the user is allowed here
	allowed, user := sessions.UserIsAllowed(r, w, "admin")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// Connect to Database
	db.Conn("user")
	defer db.DB.Close()

	util.Error(200, "OK", getDetails(db.DB, user), "Response", w)
	return
}
