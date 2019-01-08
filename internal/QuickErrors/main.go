package QuickErrors

import (
	"net/http"
	"Util"
)

func InteralServerError(w http.ResponseWriter) {
	Util.SolidError(500, "Internal Server Error", "Something went wrong", "Something went wrong...", w)
}


func AlreadyLoggedIn(w http.ResponseWriter) {
	Util.SolidError(400, "Hmp??", "Already Logged In", "Something went wrong...", w)
}

func NotLoggedIn(w http.ResponseWriter) {
	Util.SolidError(400, "Hmp??", "User is not logged in or can't access this section of the API", "Something went wrong...", w)
}

func MalformedRequest(w http.ResponseWriter, error string) {
	Util.SolidError(400, "Malformed Request", error, "Invalid Request", w)
}

func NotEnoughPrivledges(w http.ResponseWriter) {
	Util.SolidError(403, "Forbidden", "User does not have sufficient privileges or is not logged in", "Invalid Request", w)
}


func OK(w http.ResponseWriter) {
	Util.SolidError(200, "OK", "Success!", "Success!", w)
}

func NotFound(w http.ResponseWriter) {
	Util.SolidError(404, "Not Found", "The file you requested could not be found on this server", "Not Found", w)
}