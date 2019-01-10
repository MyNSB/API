package quickerrors

import (
	"fmt"
	"net/http"
)

// Function to remove all that ugly code error e.t.c
func dropErr(status int, statusMessage string, body string, title string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"Status":{"Code": %d, "Status Message":"%s"},"Message": {"Title":"%s", "Body":"%s"}}`, status, statusMessage, title, body)))
}

func InternalServerError(w http.ResponseWriter) {
	dropErr(500, "Internal Server Error", "Something went wrong", "Something went wrong...", w)
}

func AlreadyLoggedIn(w http.ResponseWriter) {
	dropErr(400, "Hmp??", "Already Logged In", "Something went wrong...", w)
}

func NotLoggedIn(w http.ResponseWriter) {
	dropErr(400, "Hmp??", "user is not logged in or can't access this section of the API", "Something went wrong...", w)
}

func MalformedRequest(w http.ResponseWriter, error string) {
	dropErr(400, "Malformed Request", error, "Invalid Request", w)
}

func NotEnoughPrivileges(w http.ResponseWriter) {
	dropErr(403, "Forbidden", "user does not have sufficient privileges or is not logged in", "Invalid Request", w)
}

func OK(w http.ResponseWriter) {
	dropErr(200, "OK", "Success!", "Success!", w)
}

func NotFound(w http.ResponseWriter) {
	dropErr(404, "Not Found", "The file you requested could not be found on this server", "Not Found", w)
}
