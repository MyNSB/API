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
	"mynsb-api/internal/user/userdetails"
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
	router.GET("/api/v1/timetable/get", timetable.RetrievalHandler) /**/
	// END TIMETABLE ==============

	// EVENTS ===================
	// Handle event creation
	router.POST("/api/v1/events/create", events.EventCreationHandler) /**/
	router.GET("/api/v1/events/get", events.EventRetrievalHandler)    /**/
	router.GET("/api/v1/events/calendar/get", calendar.CalendarRetrievalHandler)
	// END EVENTS =====================

	// AUTHENTICATION AND USERS AND ADMINS ======================
	// Handle authentication
	router.POST("/api/v1/user/auth", auth.UserAuthenticationHandler)     /**/
	router.POST("/api/v1/admin/auth", auth.AdminAuthenticationHandler)   /**/
	router.POST("/api/v1/user/logout", auth.LogoutRequestHandler)        /**/
	router.GET("/api/v1/user/getDetails", userdetails.RetrievalHandler)  /**/
	router.GET("/api/v1/admin/getDetails", admin.DetailRetrievalHandler) /**/
	// END AUTHENTICATION AND USERS ==================

	// Handler for file server for assets e.t.c
	router.ServeFiles("/api/v1/assets/*filepath", http.Dir("assets")) /**/

	// BELL TIMES ====================================
	router.GET("/api/v1/belltimes/get", belltimes.RetrievalHandler) /**/
	// END BELL TIMES ================================

	// 4U STUFF =======================================
	router.GET("/api/v1/4U/get", _4U.IssueRetrievalHandler)              /**/
	router.POST("/api/v1/4U/create/issue", _4U.IssueCreationHandler)     /**/
	router.POST("/api/v1/4U/create/article", _4U.ArticleCreationHandler) /**/
	// END 4U STUFF ===================================

	// REMINDERS ======================================
	router.POST("/api/v1/reminders/create", reminders.CreationHandler)       /**/
	router.GET("/api/v1/reminders/get/*reqType", reminders.RetrievalHandler) /**/
	router.POST("/api/v1/reminders/delete", reminders.DeletionHandler)       /**/
	// END REMINDERS STUFF ============================

	// WEEK A B STUFF =================================
	router.GET("/api/v1/week/get", week.GetHandler) /**/
	// END WEEK A B STUFF =============================

	c := cors.AllowAll().Handler(router)


	http.ListenAndServe(":8080", c)
}
