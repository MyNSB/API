package timetable

import (
	"encoding/json"
	"errors"
	"github.com/buger/jsonparser"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
	"regexp"
	"strconv"
	"path/filepath"
)






// Datastructure for holding a class
type Class struct {
	Teacher   string
	ClassRoom string
	Subject   string
	Period    string
}
type TimeTable []Class
var	timetableJSONDump map[string]interface{}



// init function reads all the timetable data into a data structure for easy access
func init() {

	// Get gopath and build a timetable directory
	gopath := util.GetGOPATH()
	timetableDir := filepath.FromSlash(gopath + "/mynsb-api/internal/timetable/daemons/Timetables.json")

	// Read the file and dump it into our data structure
	var jsonDataBuffer []byte
	jsonDataBuffer, _ = ioutil.ReadFile(timetableDir)
	json.Unmarshal(jsonDataBuffer, &timetableJSONDump)
}












// RETRIEVAL FUNCTIONS

// getSubject returns a subject from the specified students timetable, the subject is matched based off the requested
// period and day
func getSubject(studentID string, requestedPeriod string, requestedDay int) (Class, error) {

	// Read data structure and extrac the student's timetables
	studentTimetables := timetableJSONDump["student_timetables"].(map[string]interface{})
	studentTimetable := studentTimetables[studentID].([]interface{})

	// Iterate over every day in the students timetable, see timetable_sample.json to see how the timetable is structured
	for _, data := range studentTimetable {

		subject := data.(map[string]interface{})
		// Determine if the subject matches the requirements
		if subject["day"].(float64) == float64(requestedDay) && requestedPeriod == subject["period"].(string) {
			// Return the matched class
			return Class{
				Teacher:   subject["teacher"].(string),
				ClassRoom: subject["room"].(string),
				Subject:   subject["class"].(string),
				Period:    requestedPeriod,
			}, nil
		}
	}

	return Class{}, errors.New("user id or period or day does not exist")
}


// retrieveEntireTimetable gets the entire table for a student given their student id
func retrieveEntireTimetable(studentID string) (interface{}, error) {

	// Get the timetables
	studentTimetables := timetableJSONDump["student_timetables"].(map[string]interface{})

	// Get the currently logged in user's timetable through their user ID
	if _, ok := studentTimetables[studentID]; ok {
		return studentTimetables[studentID], nil
	}

	return nil, errors.New("user or period or day does not exist")
}


// getWholeDay returns a student's timetable for the entire day
func getWholeDay(requestedDay int, studentID string) ([]Class, error) {

	var timetable []Class

	allTimetables := timetableJSONDump["student_timetables"].(map[string]interface{})
	studentTimetable := allTimetables[studentID].([]interface{})

	for _, subject := range studentTimetable {
		// Extract the subject information
		subjectInfo := subject.(map[string]interface{})

		// Determine if this day matches the day requested by the user
		if subjectInfo["day"].(float64) == float64(requestedDay) {

			// Handle unsupervised periods, if there is no teacher entry that means the period is unsupervised
			teacher := "Unsupervised"
			if _, ok := subjectInfo["teacher"]; ok {
				teacher = subjectInfo["teacher"].(string)
			}

			// Create timetable
			timetable = append(timetable, Class{
				Teacher:   teacher,
				ClassRoom: subjectInfo["room"].(string),
				Subject:   subjectInfo["class"].(string),
				Period:    subjectInfo["period"].(string),
			})


			// due to the consecutive nature of the timetables we can reduce the number of iterations after the first subject that matches our day has been found
			// Check the timetablesample.json file
		} else if len(timetable) > 0 {
			break
		}

	}

	timetableResponse := sortTimetable(timetable)

	return timetableResponse, nil
}












// UTILITY FUNCTIONS

// getMaxKey returns the largest key in a map of integers
func getMaxKey(list map[int]Class) int {
	max := -9999999

	for key := range list {
		if key > max {
			max = key
		}
	}

	return max
}


// sortTimetables sorts the timetables into an order that the user will understand
func sortTimetable(timetable []Class) []Class {

	initialClass := Class{}

	// Construct a mask that dictates the "desired" arrangement for the timetable
	timetableMask := make(map[int]Class)
	for _, timetableJson := range timetable {
		if timetableJson.Period == "RC" {
			initialClass = timetableJson
			continue
		}

		p, _ := strconv.Atoi(timetableJson.Period)
		// Append this period where it should be
		timetableMask[p] = timetableJson
	}

	// Match the actual timetable with the mask
	periodNo := getMaxKey(timetableMask)
	finMask := make([]Class, periodNo)
	for x := range timetableMask {
		finMask[x-1] = timetableMask[x]
	}
	finalTimetable := finMask[:]
	finalTimetable = append(finalTimetable, initialClass)

	return finalTimetable
}


// GetStudentGrade returns the grade a student with a certain student ID is int
func GetStudentGrade(StudentID string) (string, error) {

	// Get the entire timetable for the current user
	student, err := retrieveEntireTimetable(StudentID)
	if err != nil {
		return "", err
	}
	// Marshall the result so we can perform a regex query on it
	studentTimetaleJSON, _ := json.Marshal(student)

	// Get the first subject in the students timetables
	firstSubject, _ := jsonparser.GetString(studentTimetaleJSON, "[0]", "class")
	var gradeMatcher = regexp.MustCompile(`\d+`)

	// Return first instance of the grade
	return string(gradeMatcher.Find([]byte(firstSubject))[0]), nil
}












// HTTP HANDLERS

// RetrievalHandler handles the exporting of a student's timetable
func RetrievalHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	var StudentID string
	allowed, user := sessions.IsUserAllowed(r, w, "user")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	Period := r.URL.Query().Get("Period")
	Day := r.URL.Query().Get("Day")

	StudentID = user.Name

	// Determine the request type, btw, this code is meant to be a meme:p its pre much a ternary
	reqType := map[bool]interface{}{true: 0, false: map[bool]int{true: 1, false: 2}[Period == "" && Day != ""]}[Period == "" && Day == ""] // 0: GetAll, 1: GetDay, 2: getSubject

	// GLOBAL data
	var resp interface{}
	var errGlob error


	switch reqType {
	case 1:
		// Convert the day into an integer
		day, _ := strconv.Atoi(Day)
		resp, errGlob = getWholeDay(day, StudentID)
		break
	case 2:
		day, _ := strconv.Atoi(Day)
		resp, errGlob = getSubject(StudentID, Period, day)
		break
	default:
		resp, errGlob = retrieveEntireTimetable(StudentID)
		break
	}

	if errGlob != nil {
		quickerrors.InternalServerError(w)
		return
	}


	jsonResp, _ := json.Marshal(resp)
	util.HTTPResponse(200, "OK", string(jsonResp), "Response", w)
}
