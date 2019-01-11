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

// USe it to set up the cookie store
var store = sessions.NewCookieStore(key)

// Function for parsing sessions and generating a user from those sessions
func ParseSessions(r *http.Request, w http.ResponseWriter) (user.User, error) {
	// Attain session
	session, _ := store.Get(r, "user-data")

	// Get the expiry date
	if !(session.Values["token"] == nil) {
		// Attain values and parse token
		// Parse jwt
		currUser, err := jwt.ReadJWT(session.Values["token"].(string))
		if err != nil {
			return user.User{}, err
		}

		// Renew the session
		session.Options.MaxAge = int(time.Duration(time.Hour * 24 * 30).Seconds())
		// Save the sessions
		session.Save(r, w)

		return jwtDataToUser(currUser), nil
	}

	// Tell the user that auth is required again
	return user.User{}, errors.New("session is invalid, please authenticate again")
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

func UserIsAllowed(r *http.Request, w http.ResponseWriter, requirements ...string) (bool, user.User) {
	currUser, err := ParseSessions(r, w)
	if err != nil {
		return false, user.User{}
	}
	return util.IsSubset(requirements, currUser.Permissions), currUser
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

func jwtDataToUser(details jwt.JWTData) user.User {
	return user.User{
		Name:        details.User,
		Password:    details.Password,
		Permissions: details.Permissions,
	}
}
