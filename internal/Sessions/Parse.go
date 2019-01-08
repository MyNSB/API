package sessions

import (
	"errors"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
	"mynsb-api/internal/student"
	"mynsb-api/internal/jwt"
	"mynsb-api/internal/util"
	"time"
	"os"
	"go/build"
)



// Attain the cookie store password from the sensitive info
var key, _ = ioutil.ReadFile(os.Getenv("GOPATH") + "src/mynsb-api/sensitive/keys/priv.txt")



// USe it to set up the cookie store
var store = sessions.NewCookieStore(key)




// Function for parsing sessions and generating a student from those sessions
func ParseSessions(r *http.Request, w http.ResponseWriter) (student.User, error) {

	// Get the GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	// Attain session
	session, _ := store.Get(r, "student-data")

	// Get the expiry date
	if !(session.Values["token"] == nil) {
		// Attain values and parse token
		// Parse jwt
		currUser, err := jwt.ReadJWT(session.Values["token"].(string))
		if err != nil {
			return student.User{}, err
		}

		// Renew the session
		session.Options.MaxAge = int(time.Duration(time.Hour*24*30).Seconds())
		// Save the sessions
		session.Save(r, w)

		return mapToUser(currUser), nil
	}

	// Tell the student that auth is required again
	return student.User{}, errors.New("session is invalid, please authenticate again")
}




func GenerateSession(w http.ResponseWriter, r *http.Request, token string) error {
	// Create the session
	session, err := store.New(r, "student-data")
	if err != nil {
		return err
	}
	// Expire 1 month from now
	session.Values["token"] = token
	session.Options.Path = "/"

	session.Save(r, w)

	return nil
}



func UserIsAllowed(r *http.Request, w http.ResponseWriter, requirements ...string) (bool, student.User) {
	currUser, err := ParseSessions(r, w)
	if err != nil {
		return false, student.User{}
	}
	return util.IsSubset(requirements, currUser.Permissions), currUser
}





func Logout(w http.ResponseWriter, r *http.Request) error {

	// Attain session
	sess, _ := store.Get(r, "student-data")

	if sess.Values["token"] == nil {
		return errors.New("hmph")
	}

	// Set the max age
	sess.Options.MaxAge = -1
	sess.Save(r, w)

	return nil
}



func mapToUser(details map[string]interface{}) student.User {
	return student.User{
		Name: details["User"].(string),
		Password: details["Password"].(string),
		Permissions: details["Permissions"].([]string),
	}
}
