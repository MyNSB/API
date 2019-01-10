package auth

import (
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"mynsb-api/internal/admin"
	"mynsb-api/internal/db"
	"mynsb-api/internal/jwt"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
)

// Http handler for admin authentication
/*
	Handler's have minimal documentation
*/
func AdminHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Determine if a student is already logged in
	_, err := sessions.ParseSessions(r, w)
	if err == nil {
		quickerrors.AlreadyLoggedIn(w)
		return
	}

	// Connect to database
	db.Conn("student")
	defer db.DB.Close()

	// Get the details for the incoming request
	details, err := getAuthDetails(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "Missing details")
		return
	}

	// Process these details
	err = processDetails(details, w, r)
	if err != nil {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// All g
	quickerrors.OK(w)
}

/*
	auth takes the username and password and checks if the admin really exists in the database, the function can be called from anywhere
	@params;
		Username string
		Password string
		db *sql.db
*/
func Auth(Username, Password string, db *sql.DB) (admin.Admin, error) {
	passwordHash := util.HashString(Password)

	// Determine if currAdmin exists in database
	if count, err := util.CheckCount(db, "SELECT * FROM admins WHERE admin_name = $1 AND admin_password = $2", Username, passwordHash); err != nil || count != 1 {
		return admin.Admin{}, errors.New("currAdmin does not exist")
	}

	// Finally query the database for the currAdmin's details
	rows, _ := db.Query("SELECT * FROM admins WHERE admin_name = $1 AND admin_password = $2", Username, passwordHash)
	defer rows.Close()

	// Actual student
	// Construct the currAdmin
	currAdmin := admin.Admin{
		Name:     Username,
		Password: Password,
	}

	// Scan the currAdmin details into the currAdmin variable
	for rows.Next() {
		currAdmin.ScanFrom(rows)
	}

	return currAdmin, nil
}

/*
	@ UTIL FUNCTIONS ==================================================
*/
/*
	getDetails takes the incoming http request and extracts the details from it
	@params;
		r *http.Request
*/
func getAuthDetails(r *http.Request) (map[string]string, error) {
	// Retrieve auth details
	Username := ""
	Password := ""

	// Attain the details
	if username, password, ok := r.BasicAuth(); !ok || username == "" || password == "" {
		return nil, errors.New("invalid request")
	} else {
		Username = username
		Password = password
	}

	// Create a map of the username and password
	toReturn := make(map[string]string)

	// Push contents into toReturn
	toReturn["username"] = Username
	toReturn["password"] = Password

	return toReturn, nil
}

/*
	processDetails takes the details map and determines if they really are an admin, it then converts that admin to a student, calculates
	the jwt representation of it generates session details for the student requesting to authenticate and then exists, phew that's a lot
	@params;
		details map[string]string
		w http.ResponseWriter
		r* http.Request
*/
func processDetails(details map[string]string, w http.ResponseWriter, r *http.Request) error {
	// Authenticate
	currAdmin, err := Auth(details["username"], details["password"], db.DB)
	if err != nil {
		return errors.New("invalid admin details")
	}

	// Convert the admin to a student so a jwtToken can be generated
	user := admin.ToUser(currAdmin)

	// Create the jwtToken
	jwtToken, err := jwt.GenJWT(userToJWTData(user))
	if err != nil {
		return errors.New("something went wrong")
	}

	// Generate a session from the given jwtToken
	err = sessions.GenerateSession(w, r, jwtToken)
	if err != nil {
		return errors.New("something went wrong")
	}

	return nil
}

/*
	@ END UTIL FUNCTIONS ==================================================
*/
