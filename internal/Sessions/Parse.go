package Sessions

import (
	"errors"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
	"User"
	"time"
)



// Attain the cookie store password from the sensitive info
var key, _ = ioutil.ReadFile("sensitive/keys/priv.txt")



// USe it to set up the cookie store
var store = sessions.NewCookieStore(key)




// Function for parsing sessions and generating a user from those sessions
func ParseSessions(r *http.Request, w http.ResponseWriter) (User.User, error) {
	// Attain session
	session, _ := store.Get(r, "user-data")

	// Get the expiry date
	if !(session.Values["token"] == nil) {
		// Attain values and parse token
		// Parse JWT
		user, err := User.ReadJWT(session.Values["token"].(string), "sensitive/keys/priv.txt")
		if err != nil {
			return User.User{}, err
		}

		// Renew the session
		session.Options.MaxAge = int(time.Duration(time.Hour*24*30).Seconds())
		// Save the sessions
		session.Save(r, w)

		return user, nil
	}

	// Tell the user that auth is required again
	return User.User{}, errors.New("session is invalid, please authenticate again")
}




func GenerateSession(w http.ResponseWriter, r *http.Request, token string) error {
	// Create the session
	session, err := store.New(r, "user-data")
	if err != nil {
		return err
	}
	// Expire 1 month from now
	session.Values["token"] = token
	session.Options.Path = "/"

	session.Save(r, w)

	return nil
}




func Logout(w http.ResponseWriter, r *http.Request) error {

	// Attain session
	sess, _ := store.Get(r, "user-data")

	if sess.Values["token"] == nil {
		return errors.New("hmph")
	}

	// Set the max age
	sess.Options.MaxAge = -1
	sess.Save(r, w)

	return nil
}
