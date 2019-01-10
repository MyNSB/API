package auth

import (
	"database/sql"
	"errors"
	"github.com/Azure/go-ntlmssp"
	_ "github.com/SermoDigital/jose"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"mynsb-api/internal/db"
	"mynsb-api/internal/jwt"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/student"
	"mynsb-api/internal/timetable"
	"mynsb-api/internal/util"
	"net/http"
	"regexp"
	_ "regexp"
	"strconv"
	"strings"
	"unicode"
	"path/filepath"
)

// TODO: Refactor spaghetti code

// Function that parses the raw html string returned by the school and returns the school details
func getStudentDetails(rawHTML string) (string, string) {
	var re = regexp.MustCompile(`<divclass="page-header"><h1>(.*)<br><small>Homepage</small></h1></div>`)
	out := strings.TrimLeft(strings.TrimRight(re.FindString(rawHTML), `<br><small>Homepage</small></h1></div>`), `<divclass="page-header"><h1>`)
	// Match the first names, last names e.t.c...
	var firstNameS = regexp.MustCompile(`[A-Z][a-z]+`)
	var lastName = regexp.MustCompile(`[A-Z]+$`)

	return firstNameS.FindString(out), lastName.FindString(out)
}

func pushAndUpdateDB(db *sql.DB, studentID string, fName string, lName string, studentYear string) {

	// Convert student id to integer
	studentIDInt, _ := strconv.Atoi(studentID)
	studentYearInt, _ := strconv.Atoi(studentYear)

	_, err := db.Query("SELECT insert_student($1, $2, $3, $4)", studentIDInt, fName, lName, studentYearInt)
	if err != nil {
		panic(err)
	}
}

/**
	Func Authenticate:
		@param StudentID int
		@param Password string

		returns a boolean representing if the login details are valid
		(bool) 1 = Valid
		(bool) 0 = Invalid
**/
func Authenticate(StudentID string, Password string) (int, string, error) {
	client := &http.Client{
		Transport: ntlmssp.Negotiator{
			RoundTripper: &http.Transport{},
		},
	}

	req, err := http.NewRequest("GET", "http://web1.northsydbo-h.schools.nsw.edu.au", nil)
	if err != nil {
		return 0, "", errors.New("error authenticating")
	}
	req.SetBasicAuth(StudentID, Password)
	res, err := client.Do(req)

	// Strip it of whitespace
	body, _ := ioutil.ReadAll(res.Body)
	resp := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, string(body))

	if err != nil {
		return 0, "", errors.New("error authenticating")
	}

	return res.StatusCode, resp, nil
}

// Http handler for authentication
func UserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	_, err := sessions.ParseSessions(r, w)
	if err == nil {
		quickerrors.AlreadyLoggedIn(w)
		return
	}

	db.Conn("admin")
	defer db.DB.Close()


	// Retrieve auth details
	var StudentId string
	var Password string

	// Attain the details and set them
	if studentId, password, ok := r.BasicAuth(); !ok || studentId == "" || password == "" {
		quickerrors.MalformedRequest(w, "Invalid Request, missing details")
		return
	} else {
		StudentId = studentId
		Password = password
	}

	// Attain jwt from these variables
	StatusCode, rawHTML, _ := Authenticate(StudentId, Password)

	if StatusCode == 200 {
		// Generate the jwtToken
		jwtToken, err := jwt.GenJWT(userToJWTData(student.User{Name: StudentId, Password: Password, Permissions: []string{"student"}}))

		// Create the session
		sessions.GenerateSession(w, r, jwtToken)

		// DATABASE STUFF ===============
		// Get the spicy details from the details function
		fnameS, lastName := getStudentDetails(rawHTML)


		// Get the GOPATH
		gopath := util.GetGOPATH()
		// Set up the timetable
		timetableDir := filepath.FromSlash(gopath + "/mynsb-api/internal/timetable/daemons/Timetables.json")
		jsonData, err := timetable.GetJson(timetableDir)
		if err != nil {
			panic(err)
			quickerrors.InternalServerError(w)
			return
		}

		// Get the year
		year, _ := timetable.GetYear(StudentId, jsonData)
		// Parse into our little function
		pushAndUpdateDB(db.DB, StudentId, fnameS, lastName, year)
		// END DATABASE ================

		// Determine that no error occurred using compound if statements
		if err != nil {
			quickerrors.InternalServerError(w)
			return
		}

		util.SolidError(200, "OK", "Logged In as: "+StudentId, "Success!", w)
		return

	}

	util.SolidError(401, "Unauthorized", "user details provided are invalid", "Unauthorized", w)

}

// LogoutHandler function
func LogoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	err := sessions.Logout(w, r)
	if err != nil {
		quickerrors.NotLoggedIn(w)
		return
	}
	quickerrors.OK(w)
	return
}
