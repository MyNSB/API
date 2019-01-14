package sessions

import (
	"errors"
	"github.com/gorilla/sessions"
	"mynsb-api/internal/filesint"
	"mynsb-api/internal/jwt"
	"mynsb-api/internal/user"
	"mynsb-api/internal/util"
	"net/http"
	"time"
)


// Attain the cookie store password from the sensitive info
var key, _ = filesint.DataDump("sensitive", "/keys/priv.txt")
// Use it to set up the cookie store
var store = sessions.NewCookieStore(key)



// CORE FUNCTIONS

// ParseSessions takes a http request, extracts the session data from it, parses it and returns a user object
func ParseSessions(r *http.Request, w http.ResponseWriter) (user.User, error) {

	// Attain session
	session, err := store.Get(r, "user-data")
	if err != nil {
		return user.User{}, err
	}

	oneMonth := time.Hour * 24 * 30

	if !(session.Values["token"] == nil) {

		currUser, err := jwt.ReadJWT(session.Values["token"].(string))
		if err != nil {
			return user.User{}, err
		}

		session.Options.MaxAge = int(time.Duration(oneMonth).Seconds())
		session.Save(r, w)

		return jwtDataToUser(currUser), nil
	}

	return user.User{}, errors.New("session is invalid, please authenticate again")
}


// GenerateSession generates a session based off a JWT-token
func GenerateSession(w http.ResponseWriter, r *http.Request, token string) error {

	// Create the session
	session, err := store.New(r, "user-data")
	if err != nil {
		return err
	}

	session.Values["token"] = token
	session.Options.Path = "/"
	session.Save(r, w)

	return nil
}


// IsUserAllowed determines if the currently logged in user meets a set of requirements
func IsUserAllowed(r *http.Request, w http.ResponseWriter, requirements ...string) (bool, user.User) {

	currUser, err := ParseSessions(r, w)
	if err != nil {
		return false, user.User{}
	}

	return util.IsSubset(requirements, currUser.Permissions), currUser
}


// Logout allows for a user to logout and destroy their session
func Logout(w http.ResponseWriter, r *http.Request) error {

	sess, _ := store.Get(r, "user-data")
	dead := -1

	if sess.Values["token"] == nil {
		return errors.New("session is invalid")
	}

	sess.Options.MaxAge = dead
	sess.Save(r, w)

	return nil
}










// jwtDataToUser converts a JWT data object into a user, this is done in an attempt to isolate the JWT module from any other project packages
func jwtDataToUser(details jwt.JWTData) user.User {
	return user.User{
		Name:        details.User,
		Password:    details.Password,
		Permissions: details.Permissions,
	}
}
