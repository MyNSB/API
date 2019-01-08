package main

import (
	_4U "4U"
	"Admin"
	"BellTimes"
	"Reminders"
	"Events"
	"Timetable"
	"User"
	"Week"
	"Events/Calendar"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"github.com/rs/cors"
)




// 404 Handler
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte(`{"Status":{"Code": 404, "Status Message":"404 Not Found"},"Message": {"Title":"An Error Occured", "Body":"The file you requested could not be found on this server."}}`))
}




// Index handler
func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte(`{"Status":{"Code": 200, "Status Message":"200 OK"},"Message": {"Title":"Hello There!", "Body":"Welcome to the MyNSB API."}}`))
}


// Main function
func main() {

	// Start router
	router := httprouter.New()


	// GENERAL ===================
	// Set 404 handler
	router.NotFound = http.HandlerFunc(NotFoundHandler)
	// Regular router functions
	router.GET("/api/v1", IndexHandler)
	// END GENERAL ===============

	// TIMETABLE ==================
	// Handle timetable exports
	router.GET("/api/v1/timetable/Get", Timetable.ExportTimetable)
	// END TIMETABLE ==============

	// EVENTS ===================
	// Handle event creation
	router.POST("/api/v1/events/Create", Events.CreateEventHandler)
	router.GET("/api/v1/events/Get", Events.GetEvents)
	router.GET("/api/v1/events/Calendar/Get", Calendar.GetCalendar)
	// END EVENTS =====================

	// AUTHENTICATION AND USERS AND ADMINS ======================
	// Handle authentication
	router.POST("/api/v1/user/Auth", User.AuthHandler)
	router.POST("/api/v1/admin/Auth", Admin.AuthHandler)
	router.POST("/api/v1/user/Logout", User.Logout)
	router.GET("/api/v1/user/GetDetails", User.GetDeatilsHandler)
	router.GET("/api/v1/admin/GetDetails", Admin.GetDetailsHandler)
	// END AUTHENTICATION AND USERS ==================

	// Handler for file server for assets e.t.c
	router.ServeFiles("/api/v1/assets/*filepath", http.Dir("assets"))

	// BELL TIMES ====================================
	router.GET("/api/v1/belltimes/Get", BellTimes.ServeBellTimes)
	// END BELL TIMES ================================

	// 4U STUFF =======================================
	router.GET("/api/v1/4U/Get", _4U.GetFourUHandler)
	router.POST("/api/v1/4U/Create", _4U.CreateFourUHandler)
	router.POST("/api/v1/4U/Create/Article", _4U.CreateFourUArticleHandler)
	// END 4U STUFF ===================================

	// REMINDERS ======================================
	router.POST("/api/v1/reminders/Create", Reminders.CreateReminderHandler)
	router.GET("/api/v1/reminders/Get/*reqType", Reminders.GetRemindersHandler)
	router.POST("/api/v1/reminders/Delete", Reminders.DeleteReminderHandler)
	// END REMINDERS STUFF ============================

	// WEEK A B STUFF =================================
	router.GET("/api/v1/week/Get", Week.GetWeek)
	// END WEEK A B STUFF =============================


	c := cors.AllowAll().Handler(router)


	// Begin listening on port 8080
	http.ListenAndServe("0.0.0.0:8080", c)
}
