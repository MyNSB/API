package Timetable

import (
	"QuickErrors"
	"Util"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"github.com/buger/jsonparser"
	"github.com/julienschmidt/httprouter"
)


/*
	Data Holders =============
 */
// Class holds a certain class during any day and any given period
type Class struct {
	Teacher string
	ClassRoom string
	Subject string
	Period string
}

type TimeTable []Class
// Would represent something like: {1 = monday week A}
/*
	END Data Holders =============
 */





/*
	CORE FUNCTIONS =============================
 */
/* GetSubject returns the subject that the user has requested through a given period and a given day
 		@params;
			StudentID string
			Period string,
			Day int,
			filepath string
 */
func GetSubject(StudentID string, Period string, Day int, filepath string) (Class, error) {

	// Create a data holder
	var data map[string]interface{}
	var err error

	if data, err = getJson(filepath); err != nil {
		return Class{}, errors.New("could not parse json")
	}

	// Read from the data container
	studentTimetables := data["student_timetables"].(map[string]interface{})
	// Convert this into a array of interfaces
	// Student timetables look like: "441567081": {
	// 												{},
	// 												{},
	// 												{}
	// 											  }
	studentTimetable := studentTimetables[StudentID].([]interface{})

	// Iterate through timetable array of all the students
	for _, data := range studentTimetable {
		details := data.(map[string]interface{})
		if details["day"].(float64) == float64(Day) && Period == details["period"].(string) {
			// Parse the details to turn it into a class
			return Class{
				Teacher: details["teacher"].(string),
				ClassRoom: details["room"].(string),
				Subject: details["class"].(string),
				Period: Period,
			}, nil
		}
	}


	// Return error if it does not exist
	return Class{}, errors.New("student id or period or day does not exist")

}


/* getYear returns the year a specific user is in, it is only really used during authentication when the user details are stored in the database
		@params;
			StudentID string
			filepath string
 */
func GetYear(StudentID string, filepath string) (string, error) {
	// Get the timetable for the correct student
	student, err := RetrieveAll(StudentID, filepath)
	if err != nil {
		return "", err
	}

	// Get the year off the first class using   r e g e x
	rawJson, _ := json.Marshal(student)

	firstSubject, _ := jsonparser.GetString(rawJson, "[0]", "class")
	// r e g e x  that boi
	var numberRegex = regexp.MustCompile(`\d+`)

	return string(numberRegex.Find([]byte(firstSubject))[0]), nil
}


/* RetrieveAll returns the timetable for a particular student
		@params;
			StudentID string
			filepath string

 */
func RetrieveAll(StudentID string, filepath string) (interface{}, error) {

	var data map[string]interface{}
	var err error

	// Retrieve the timetables
	if data, err = getJson(filepath); err != nil {
		return nil, errors.New("could not read timetable dump")
	}

	// Get the timetables
	studentTimetables := data["student_timetables"].(map[string]interface{})

	// Get the currently logged in student's timetable through their student ID
	if _, ok := studentTimetables[StudentID]; ok {
		return studentTimetables[StudentID], nil
	}

	return nil, errors.New("student or period or day does not exist")
}


/* GetWholeDay returns the timetable for a student on a given day
		@params;
			day int
			filepath string
 */
func GetWholeDay(day int, studentID string, filepath string) ([]Class, error) {

	// Timetable type that we start with
	timetables := []Class{}

	var data map[string]interface{}
	var err error

	// Retrieve the timetables
	if data, err = getJson(filepath); err != nil {
		return TimeTable{}, errors.New("could not get timetables")
	}


	studentTimetables := data["student_timetables"].(map[string]interface{})
	studentsTimetable := studentTimetables[studentID].([]interface{})

	// Iterate through timetable array of all the students
	for _, data := range studentsTimetable {
		details := data.(map[string]interface{})
		if details["day"].(float64) == float64(day) {
			// Handle unsupervised periods
			teacher := "Unsupervised"
			if _, ok := details["teacher"]; ok {
				teacher = details["teacher"].(string)
			}

			// Create timetable
			timetables = append(timetables, Class{
				Teacher: teacher,
				ClassRoom: details["room"].(string),
				Subject: details["class"].(string),
				Period: details["period"].(string),
			})
		}
	}

	json_l := Class{}
	// This now needs to be sorted
	timetableMask := make(map[int]Class)
	for _, json := range timetables {
		if json.Period == "RC" {
			json_l = json
			continue
		}

		p, _ := strconv.Atoi(json.Period)
		// Append this period where it should be
		timetableMask[p] = json
	}
	// Sort this mask
	periodNo := getMaxKey(timetableMask)
	finMask := make([]Class, periodNo)
	for x, _ := range timetableMask {
		finMask[x-1] = timetableMask[x]
	}
	ln := finMask[:]
	ln = append(ln, json_l)





	return ln, nil
}
/*
	END CORE FUNCTIONS =============================
 */





/*
	UTIL FUNCTIONS =========================
 */
/* getMaxKey returns the maximum key in a map
		@params;
			map[int]interface{}
 */
 func getMaxKey(list map[int]Class) int {
 	max := -9999999

 	for key, _ := range list {
 		if key > max {
 			max = key
		}
	}

	return max
 }

/* getJson retrieves the json dump as a map of strings and interfaces
		@params;
			filepath string
 */
func getJson(filepath string) (map[string]interface{}, error) {
	// Get the jsonpath
	jsonPath := filepath
	// Holders for stuff like json data and errors
	var jsonData []byte
	var err error

	// Read everything from the timetable export
	if jsonData, err = ioutil.ReadFile(jsonPath); err != nil {
		return nil, errors.New("could not open timetable dump")
	}

	// Create a data holder
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, errors.New("could not convert to JSON")
	}

	// Otherwise return the json data as a map
	return data, nil
}
 /*
 	END UTIL FUNCTIONS =====================
 */

// Http handler for timetable exports
// Should be moved somewhere else
func ExportTimetable(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	typeReq := ""

	var StudentID string

	allowed, user := Util.UserIsAllowed(r, w, "student")
	if !allowed {
		QuickErrors.NotEnoughPrivledges(w)
		return
	}
	// Set the student id variable from the user
	StudentID = user.Name

	// Overview........
	// Get the user details
	Period := r.URL.Query().Get("Period")
	Day := r.URL.Query().Get("Day")

	// Determine the type of request being sent
	if Period == "" && Day == "" {
		typeReq = "GetAll"
	} else if Period == "" && Day != "" {
		typeReq = "GetDay"
	}

	// Perform a request given the data we are given
	if typeReq == "GetSubject" {
		// Shift through it and read it carefully
		day, _ := strconv.Atoi(Day)
		resp, err := GetSubject(StudentID, Period, day, "src/Timetable/Daemons/Timetables.json")
		if err != nil {
			QuickErrors.InteralServerError(w)
			return
		}

		// Return that
		jsonResp, _ := json.Marshal(resp)
		// Return response
		Util.Error(200, "OK", string(jsonResp), "Response", w)

	} else if typeReq == "GetAll" {
		// Attain data
		Data, err := RetrieveAll(StudentID, "src/Timetable/Daemons/Timetables.json")
		if err != nil {
			QuickErrors.InteralServerError(w)
			return
		}

		// Convert to json
		jsonresp, _ := json.Marshal(Data)

		Util.Error(200, "OK", string(jsonresp), "Response", w)

	} else if typeReq == "GetDay" {
		// Convert the day
		day, _ := strconv.Atoi(Day)

		// Attain data
		Data, err := GetWholeDay(day, StudentID, "src/Timetable/Daemons/Timetables.json")
		if err != nil {
			QuickErrors.InteralServerError(w)
			return
		}

		// Convert to json
		jsonresp, _ := json.Marshal(Data)

		Util.Error(200, "OK", string(jsonresp), "Response", w)
	}
}
