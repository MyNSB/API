package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	_4U "mynsb-api/internal/4U"
	"mynsb-api/internal/admin"
	"mynsb-api/internal/auth"
	"mynsb-api/internal/belltimes"
	"mynsb-api/internal/events"
	"mynsb-api/internal/events/calendar"
	"mynsb-api/internal/reminders"
	"mynsb-api/internal/student/userdetails"
	"mynsb-api/internal/timetable"
	"mynsb-api/internal/week"
	"net/http"
)

// NotFoundHandler deals is the response to a 404 request
func NotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/json")
	w.Write([]byte(`{"Status":{"Code": 404, "Status Message":"404 Not Found"},"Message": {"Title":"An Error Occurred", "Body":"The file you requested could not be found on this server."}}`))
}

// IndexHandler is just our plain index page
func IndexHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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
	router.NotFound = http.HandlerFunc(NotFoundHandler) /**/
	// Regular router functions
	router.GET("/api/v1", IndexHandler) /**/
	// END GENERAL ===============

	// TIMETABLE ==================
	// Handle timetable exports
	router.GET("/api/v1/timetable/Get", timetable.ExportHandler) /**/
	// END TIMETABLE ==============

	// EVENTS ===================
	// Handle event creation
	router.POST("/api/v1/events/Create", events.CreateHandler) /**/
	router.GET("/api/v1/events/Get", events.GetHandler) /**/
	router.GET("/api/v1/events/calendar/Get", calendar.GetHandler)
	// END EVENTS =====================

	// AUTHENTICATION AND USERS AND ADMINS ======================
	// Handle authentication
	router.POST("/api/v1/user/auth", auth.UserHandler) /**/
	router.POST("/api/v1/admin/auth", auth.AdminHandler) /**/
	router.POST("/api/v1/user/Logout", auth.LogoutHandler) /**/
	router.GET("/api/v1/user/GetDetails", userdetails.GetHandler) /**/
	router.GET("/api/v1/admin/GetDetails", admin.GetHandler) /**/
	// END AUTHENTICATION AND USERS ==================

	// Handler for file server for assets e.t.c
	router.ServeFiles("/api/v1/assets/*filepath", http.Dir("assets")) /**/

	// BELL TIMES ====================================
	router.GET("/api/v1/belltimes/Get", belltimes.GetHandler) /**/
	// END BELL TIMES ================================

	// 4U STUFF =======================================
	router.GET("/api/v1/4U/Get", _4U.GetIssueHandler) /**/
	router.POST("/api/v1/4U/Create/Issue", _4U.CreateIssueHandler) /**/
	router.POST("/api/v1/4U/Create/Article", _4U.CreateArticleHandler) /**/
	// END 4U STUFF ===================================

	// REMINDERS ======================================
	router.POST("/api/v1/reminders/Create", reminders.CreateHandler) /**/
	router.GET("/api/v1/reminders/Get/*reqType", reminders.GetHandler) /**/
	router.POST("/api/v1/reminders/Delete", reminders.DeleteHandler) /**/
	// END REMINDERS STUFF ============================

	// WEEK A B STUFF =================================
	router.GET("/api/v1/week/Get", week.GetHandler) /**/
	// END WEEK A B STUFF =============================

	c := cors.AllowAll().Handler(router)

	// Begin listening on port 8080
	http.ListenAndServe("0.0.0.0:8080", c)
}
