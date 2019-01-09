package main

import (
	Post4U "4U/Create"
	Get4U "4U/Get"
	"Events/Create"
	"QuickErrors"
	SportsLocationsRemove "SportLocations/Update/Remove"
	"Util"
	Admin "admin/auth"
	BellTimesGet "belltimes/Get"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	User "student/auth"
)

// 404 Handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "Application/json")
	w.Write([]byte(`{"Status":{"Code": 404, "Status Message":"404 Not Found"},"Message": {"Title":"An Error Occurred", "Body":"The file you requested could not be found on this server."}}`))
}

// Index handler
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "Application/json")
	w.Write([]byte(`{"Status":{"Code": 200, "Status Message":"200 OK"},"Message": {"Title":"Hello There!", "Body":"Welcome to the MyNSB API."}}`))
}

// Http handler for timetable exports
// Should be moved somewhere else
func ExportTimetable(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	typeReq := ""

	// Overview........
	// I needed to update an existing function so instead of actually taking the time to do it I came up with a super hacked up solution

	// Get the student details
	StudentID := r.URL.Query().Get("Student_ID")
	Period := r.URL.Query().Get("Period")
	Day := r.URL.Query().Get("Day")

	if Period == "" || Day == "" {
		typeReq = "GetAll"
	}

	// Parse dat form
	r.ParseForm()

	if typeReq == "GetSubject" {
		// Shift through it and read it carefully
		day, _ := strconv.Atoi(Day)
		resp, err := RetreiveTimetable.GetSubject(StudentID, Period, day, "src/timetable/daemons/Timetables.json")
		if err != nil {
			QuickErrors.InteralServerError(w)
			return
		}

		// Return that
		jsonResp, _ := json.Marshal(resp)
		// Return response
		Util.Error(200, "OK", string(jsonResp), "Response", w)

	} else if typeReq == "GetAll" {
		// Get the student id
		StudentID := r.URL.Query().Get("Student_ID")
		// Check that it exists
		if StudentID == "" {
			Util.Error(400, "Malformed Request", "Required Parameters have not been met, please read the API docs", "Invalid Request", w)
			return
		}

		// Attain data
		Data, err := RetreiveTimetable.RetrieveAll(StudentID, "src/timetable/daemons/Timetables.json")
		if err != nil {
			QuickErrors.InteralServerError(w)
			return
		}

		// Convert to json
		jsonresp, _ := json.Marshal(Data)

		Util.Error(200, "OK", string(jsonresp), "Response", w)

	} else {
		NotFoundHandler(w, r)
	}
}

// Main function
func main() {

	// Start router
	router := httprouter.New()

	// GENERAL ===================
	// Set 404 handler
	router.NotFound = http.HandlerFunc(NotFoundHandler)
	// Regular router functions
	router.GET("/", IndexHandler)
	// END GENERAL ===============

	// TIMETABLE ==================
	// Handle timetable exports
	router.GET("/timetable/Get", ExportTimetable)
	// END TIMETABLE ==============

	// EVENTS ===================
	// Handle event creation
	router.POST("/events/Create", Create.CreateEventHandler)
	router.GET("/events/Get", GetEvent.GetEvents)
	router.GET("/events/calendar/Get", GetCalendar.GetCalendar)
	// END EVENTS =====================

	// AUTHENTICATION AND USERS AND ADMINS ======================
	// Handle authentication
	router.POST("/student/auth", User.AuthHandler)
	router.POST("/admin/auth", Admin.AuthHandler)
	router.POST("/student/logout", User.Logout)
	// END AUTHENTICATION AND USERS ==================

	// Handler for file server for assets e.t.c
	router.ServeFiles("/assets/*filepath", http.Dir("assets"))

	// BELL TIMES ====================================
	router.GET("/belltimes/Get", BellTimesGet.ServeTimeTables)
	// END BELL TIMES ================================

	// 4U STUFF =======================================
	router.GET("/4U/Get", Get4U.GetFourUHandler)
	router.POST("/4U/Upload", Post4U.Create)
	// END 4U STUFF ===================================

	// SPORTS LOCATIONS ===============================
	router.ServeFiles("/sportsLocations/Get/*filepath", http.Dir("src/SportLocations/Data"))
	router.POST("/sportsLocations/Update/Delete", SportsLocationsRemove.DeleteHandler)
	// END SPORTS LOCATIONS ===========================

	// Begin listening on port 8080
	http.ListenAndServe("0.0.0.0:8080", router)
}
