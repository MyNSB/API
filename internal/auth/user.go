package auth

import (
	"database/sql"
	"github.com/Azure/go-ntlmssp"
	_ "github.com/SermoDigital/jose"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"mynsb-api/internal/db"
	"mynsb-api/internal/jwt"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/user"
	"mynsb-api/internal/timetable"
	"mynsb-api/internal/util"
	"net/http"
	"regexp"
	_ "regexp"
	"strconv"
	"strings"
	"unicode"
)



// AUTHENTICATION FUNCTIONS

// authenticateStudent takes a student's provided details and checks weather those details exist within the school's system
func authenticateStudent(StudentID string, Password string) (int, string) {

	// Setup NTLM client
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}

	// Send the request
	req, _ := http.NewRequest("GET", "http://web1.northsydbo-h.schools.nsw.edu.au", nil)
	req.SetBasicAuth(StudentID, Password)
	res, err := client.Do(req)
	if err != nil {
		return 0, ""
	}

	body, _ := ioutil.ReadAll(res.Body)
	return res.StatusCode, sanitise(body)
}











// UTILITY FUNCTIONS

// sanitise strips the whitespace from the input string
func sanitise(input []byte) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, string(input))
}


// getStudentDetails takes a raw HTML dump of the school's classery and extract's the current student's details from it
func getStudentName(rawHTML string) (string, string) {
	// Regex to match the student details from the classert
	var re = regexp.MustCompile(`<divclass="page-header"><h1>(.*)<br><small>Homepage</small></h1></div>`)

	// Search for the user data using the regex and remove all the unnecessary information
	userData := strings.TrimLeft(strings.TrimRight(re.FindString(rawHTML), `<br><small>Homepage</small></h1></div>`), `<divclass="page-header"><h1>`)

	// Regex for extracting further information
	var firstNameMatch = regexp.MustCompile(`[A-Z][a-z]+`)
	var lastNameMatch  = regexp.MustCompile(`[A-Z]+$`)

	return firstNameMatch.FindString(userData), lastNameMatch.FindString(userData)
}


// insertStudentIntoDB inserts a student into our database based off the provided information
func insertStudentIntoDB(db *sql.DB, studentID string, fName string, lName string, studentYear string) {

	// Convert userID and the student year into integers
	studentIDInt, _ := strconv.Atoi(studentID)
	studentYearInt, _ := strconv.Atoi(studentYear)

	_, err := db.Query("SELECT insert_student($1, $2, $3, $4)", studentIDInt, fName, lName, studentYearInt)
	if err != nil {
		// yikes
		panic(err)
	}
}












// HTTP HANDLERS

// UserAuthenticationHandler handles a user's request to log into the app
func UserAuthenticationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Determine if the user is already logged in
	_, notLoggedIn := sessions.ParseSessions(r, w)
	if notLoggedIn == nil {
		quickerrors.AlreadyLoggedIn(w)
		return
	}

	// connect to DB
	db.Conn("admin")
	defer db.DB.Close()


	studentID := ""
	password  := ""
	ok		  := true

	// Attain the incoming student's details
	if studentID, password, ok = r.BasicAuth(); !ok || !util.IsSet(studentID, password) {
		quickerrors.MalformedRequest(w, "Invalid Request, missing details")
		return
	}

	statusCode, authResponse := authenticateStudent(studentID, password)

	if statusCode == 200 {
		// Generate the jwtToken for the user and generate a session for the user too
		jwtToken, err := jwt.GenJWT(
			userToJWTData(
				user.User{
					Name: studentID,
					Password: password,
					Permissions: []string{"user"}}))
		sessions.GenerateSession(w, r, jwtToken)


		firstName, lastName := getStudentName(authResponse)
		studentGrade, err := timetable.GetStudentGrade(studentID)
		insertStudentIntoDB(db.DB, studentID, firstName, lastName, studentGrade)

		if err != nil {
			quickerrors.InternalServerError(w)
			return
		}

		util.HTTPResponseArr(200, "OK", "Logged In as: "+studentID, "Success!", w)
		return
	}

	util.HTTPResponseArr(401, "Unauthorized", "user details provided are invalid", "Unauthorized", w)
	return
}


// LogoutRequestHandler handles a logout request from a student or an admin
func LogoutRequestHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	err := sessions.Logout(w, r)
	if err != nil {
		quickerrors.NotLoggedIn(w)
		return
	}
	quickerrors.OK(w)
	return
}
