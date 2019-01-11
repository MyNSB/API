package auth

import (
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"mynsb-api/internal/db"
	"mynsb-api/internal/jwt"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
	"mynsb-api/internal/user"
)


// AUTHENTICATION FUNCTION

// authenticateAdmin takes an admin's details and checks weather they exist in the DB or not :)
func authenticateAdmin(Username, Password string, db *sql.DB) (user.User, error) {

	// Hash the password
	passwordHash := util.HashString(Password)

	// Determine if current Admin exists in database
	numAdmins, _ := util.NumResults(db, "SELECT * FROM admins WHERE admin_name = $1 AND admin_password = $2", Username, passwordHash)
	if numAdmins != 1 {
		return user.User{}, errors.New("currAdmin does not exist")
	}

	// Finally, get the admin's details
	// NOTE: This section may be optimised by amalgamating the request for the number of admins with the query for the admin's details but that produces messier and unreadable code
	// Also the change in performance is not worth it
	rows, _ := db.Query("SELECT * FROM admins WHERE admin_name = $1 AND admin_password = $2", Username, passwordHash)
	defer rows.Close()

	// Attain the result
	currAdmin := user.User{
		Name:     Username,
		Password: Password,
	}
	for rows.Next() {
		currAdmin.ScanSQLIntoAdmin(rows)
	}

	return currAdmin, nil
}












// UTILITY FUNCTIONS

// parseRequest takes the incoming request and attains the sent parameters from it
func parseRequest(r *http.Request) (map[string]string, error) {
	// Retrieve auth details
	username := ""
	password := ""
	ok		 := true

	// Attain the parameters
	if username, password, ok = r.BasicAuth(); !ok || !util.IsSet(username, password) {
		return nil, errors.New("invalid request")
	}

	return  map[string]string{
		"username": username,
		"password": password,
	}, nil
}


// processDetails handles the core of the authentication including JWT and session generation
func processDetails(details map[string]string, w http.ResponseWriter, r *http.Request) error {

	currAdmin, err := authenticateAdmin(details["username"], details["password"], db.DB)
	if err != nil {
		return errors.New("invalid admin details")
	}

	// Create the jwtToken and then set up a session
	jwtToken, _ := jwt.GenJWT(userToJWTData(currAdmin))
	sessions.GenerateSession(w, r, jwtToken)


	return nil
}












// HTTP HANDLERS


// AdminAuthenticationHandler is a HTTP handler that handles an authentication request from an admin
func AdminAuthenticationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	_, notLoggedIn := sessions.ParseSessions(r, w)
	if notLoggedIn == nil {
		quickerrors.AlreadyLoggedIn(w)
		return
	}

	// connect to database
	db.Conn("user")
	defer db.DB.Close()

	// Get the details for the incoming request
	details, err := parseRequest(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "Missing details")
		return
	}

	err = processDetails(details, w, r)
	if err != nil {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	quickerrors.OK(w)
}