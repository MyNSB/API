package main

import (
	_4U "mynsb-api/internal/4U"
	"mynsb-api/internal/admin"
	"mynsb-api/internal/belltimes"
	"mynsb-api/internal/events"
	"mynsb-api/internal/reminders"
	"mynsb-api/internal/events/calendar"
	"mynsb-api/internal/student/userdetails"
	"mynsb-api/internal/auth"
	"mynsb-api/internal/week"
	"mynsb-api/internal/timetable"
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
	router.GET("/mynsb-api/v1", IndexHandler)
	// END GENERAL ===============

	// TIMETABLE ==================
	// Handle timetable exports
	router.GET("/mynsb-api/v1/timetable/Get", timetable.ExportTimetable)
	// END TIMETABLE ==============

	// EVENTS ===================
	// Handle event creation
	router.POST("/mynsb-api/v1/events/Create", events.CreateEventHandler)
	router.GET("/mynsb-api/v1/events/Get", events.GetEvents)
	router.GET("/mynsb-api/v1/events/calendar/Get", calendar.GetCalendar)
	// END EVENTS =====================

	// AUTHENTICATION AND USERS AND ADMINS ======================
	// Handle authentication
	router.POST("/mynsb-api/v1/student/auth", auth.UserAuthHandler)
	router.POST("/mynsb-api/v1/admin/auth", auth.AdminAuthHandler)
	router.POST("/mynsb-api/v1/student/Logout", auth.Logout)
	router.GET("/mynsb-api/v1/student/GetDetails", userdetails.GetDetailsHandler)
	router.GET("/mynsb-api/v1/admin/GetDetails", admin.GetDetailsHandler)
	// END AUTHENTICATION AND USERS ==================

	// Handler for file server for assets e.t.c
	router.ServeFiles("/mynsb-api/v1/assets/*filepath", http.Dir("assets"))

	// BELL TIMES ====================================
	router.GET("/mynsb-api/v1/belltimes/Get", belltimes.ServeBellTimes)
	// END BELL TIMES ================================

	// 4U STUFF =======================================
	router.GET("/mynsb-api/v1/4U/Get", _4U.GetFourUHandler)
	router.POST("/mynsb-api/v1/4U/Create", _4U.CreateFourUHandler)
	router.POST("/mynsb-api/v1/4U/Create/Article", _4U.CreateFourUArticleHandler)
	// END 4U STUFF ===================================

	// REMINDERS ======================================
	router.POST("/mynsb-api/v1/reminders/Create", reminders.CreateReminderHandler)
	router.GET("/mynsb-api/v1/reminders/Get/*reqType", reminders.GetRemindersHandler)
	router.POST("/mynsb-api/v1/reminders/Delete", reminders.DeleteReminderHandler)
	// END REMINDERS STUFF ============================

	// WEEK A B STUFF =================================
	router.GET("/mynsb-api/v1/week/Get", week.GetWeek)
	// END WEEK A B STUFF =============================


	c := cors.AllowAll().Handler(router)


	// Begin listening on port 8080
	http.ListenAndServe("0.0.0.0:8080", c)
}
