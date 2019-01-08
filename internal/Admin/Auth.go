package Admin

import (
	"DB"
	"Util"
	"database/sql"
	"errors"
	"net/http"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"User"
	"Sessions"
	"QuickErrors"
)



// Http handler for admin authentication
/*
	Handler's have minimal documentation
 */
func AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Determine if a user is already logged in
	_, err := Sessions.ParseSessions(r, w)
	if err == nil {
		QuickErrors.AlreadyLoggedIn(w)
		return
	}

	// Connect to database
	Util.Conn("sensitive", "database", "student")
	defer DB.DB.Close()

	// Get the details for the incoming request
	details, err := getAuthDetails(r)
	if err != nil {
		QuickErrors.MalformedRequest(w, "Missing details")
		return
	}

	// Process these details
	err = processDetails(details, w, r)
	if err != nil {
		QuickErrors.NotEnoughPrivledges(w)
		return
	}

	// All g
	QuickErrors.OK(w)
}


/*
	Auth takes the username and password and checks if the admin really exists in the database, the function can be called from anywhere
	@params;
		Username string
		Password string
		db *sql.DB
 */
func Auth(Username, Password string, db *sql.DB) (Admin, error) {
	passwordHash := Util.HashString(Password)

	// Determine if admin exists in database
	if count, err := Util.CheckCount(db,"SELECT * FROM admins WHERE admin_name = $1 AND admin_password = $2", Username, passwordHash); err != nil || count != 1 {
		return Admin{}, errors.New("admin does not exist")
	}


	// Finally query the database for the admin's details
	rows, _ := db.Query("SELECT admin_permissions FROM admins WHERE admin_name = $1 AND admin_password = $2", Username, passwordHash)
	defer rows.Close()


	// Actual user
	// Construct the admin
	admin := Admin{
		Name:     Username,
		Password: Password,
	}


	// Scan the admin details into the admin variable
	for rows.Next() {
		admin.ScanFrom(rows)
	}

	return admin, nil
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
	processDetails takes the details map and determines if they really are an admin, it then converts that admin to a user, calculates
	the JWT representation of it generates session details for the user requesting to authenticate and then exists, phew thats a lot
	@params;
		details map[string]string
		w http.ResponseWriter
		r* http.Request
 */
func processDetails(details map[string]string, w http.ResponseWriter, r* http.Request) error {
	// Authenticate
	admin, err := Auth(details["username"], details["password"], DB.DB)
	if err != nil {
		return errors.New("invalid admin details")
	}

	// Convert the admin to a user so a jwt can be generated
	user := AdminToUser(admin)

	// Create the jwt
	jwt, err := User.GenJWT(user, "sensitive/keys/priv.txt")
	if err != nil {
		return errors.New("something went wrong")
	}

	// Generate a session from the given jwt
	err = Sessions.GenerateSession(w, r, jwt)
	if err != nil {
		return errors.New("something went wrong")
	}

	return nil
}

/*
	@ END UTIL FUNCTIONS ==================================================
 */