package User

import (
	"Sessions"
	"Util"
	"errors"
	"github.com/Azure/go-ntlmssp"
	_ "github.com/SermoDigital/jose"
	"github.com/julienschmidt/httprouter"
	"net/http"
	_ "regexp"
	"database/sql"
	"io/ioutil"
	"strings"
	"unicode"
	"regexp"
	// Required to get the year of the student
	"Timetable"
	"DB"
	"QuickErrors"
	"strconv"
)


// TODO: Refactor spaghetti code

// Function that parses the raw html string returned by the school and returns the school details
func getStudentDetails(rawHTML string) (string, string){
	var re = regexp.MustCompile(`<divclass="page-header"><h1>(.*)<br><small>Homepage</small></h1></div>`)
	out := strings.TrimLeft(strings.TrimRight(re.FindString(rawHTML),`<br><small>Homepage</small></h1></div>`),`<divclass="page-header"><h1>`)
	// Match the first names, last names e.t.c...
	var firstNameS = regexp.MustCompile(`[A-Z][a-z]+`)
	var lastName = regexp.MustCompile(`[A-Z]+$`)

	return firstNameS.FindString(out), lastName.FindString(out)
}



func pushAndUpdateDB(db *sql.DB, studentID string, fname string, lname string, studentYear string) {

	// Convert student id to integer
	studentIDInt, _ := strconv.Atoi(studentID)
	studentYearInt, _ := strconv.Atoi(studentYear)

	_, err := db.Query("SELECT insert_student($1, $2, $3, $4)", studentIDInt, fname, lname, studentYearInt)
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
func AuthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	_, err := Sessions.ParseSessions(r, w)
	if err == nil {
		QuickErrors.AlreadyLoggedIn(w)
		return
	}


	Util.Conn("sensitive", "database", "admin")
	defer DB.DB.Close()

	// Retrieve auth details
	var StudentId string
	var Password string

	// Attain the details and set them
	if studentId, password, ok := r.BasicAuth(); !ok || studentId == "" || password == "" {
		QuickErrors.MalformedRequest(w, "Invalid Request, missing details")
		return
	} else {
		StudentId = studentId
		Password = password
	}

	// Attain JWT from these variables
	StatusCode, rawHTML, _ := Authenticate(StudentId, Password)

	if StatusCode == 200 {
		// Generate the JWT
		jwt, err := GenJWT(User{Name: StudentId, Password: Password, Permissions: []string{"student"}}, "sensitive/keys/priv.txt")

		// Create the session
		Sessions.GenerateSession(w, r ,jwt)

		// DATABASE STUFF ===============
		// Get the spicy details from the details function
		fnameS, lastName := getStudentDetails(rawHTML)
		// Get the year
		year, _ := Timetable.GetYear(StudentId, "src/Timetable/Daemons/Timetables.json")
		// Parse into our little function
		pushAndUpdateDB(DB.DB, StudentId, fnameS, lastName, year)
		// END DATABASE ================


		// Determine that no error occurred using compound if statements
		if err != nil {
			QuickErrors.InteralServerError(w)
			return
		}

		Util.SolidError(200, "OK", "Logged In as: " + StudentId, "Success!", w)
		return

	}

	Util.SolidError(401, "Unauthorized", "Details provided are invalid", "Unauthorized", w)

}



// Logout function
func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {


	err := Sessions.Logout(w, r)
	if err != nil {
		QuickErrors.NotLoggedIn(w)
		return
	}
	QuickErrors.OK(w)
	return
}