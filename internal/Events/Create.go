package Events

import (
	"DB"
	"User"
	"Util"
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/metakeule/fmtdate"
	"io"
	"net/http"
	"os"
	"time"
	"QuickErrors"
)











/**
	Func CreateEvent:
		@param event Event

		returns nothing and just creates an event
**/
func Create(user User.User, event Event, db *sql.DB) error {

	// Determine if the user is an admin
	if Util.IsAdmin(user) {


		// EVENT EXISTS DETERMINATION ===============================
		// Check that the event does not exist
		// Determine the count using the util function
		if count, _ := Util.CheckCount(db, "SELECT * FROM events WHERE event_name = $1 AND event_organiser = $2", event.EventName, event.EventOrganiser); count > 0 {
			return errors.New("event already exists")
		}
		// EVENT EXISTS DETERMINATION END ===============================



		// IMAGE CREATION =========================================
		// Create image for the event

		if _, err := os.Stat("assets/Events/" + event.EventOrganiser); os.IsNotExist(err) {
			os.Mkdir("assets/Events/" + event.EventOrganiser, 0777)
		}

		os.Mkdir("assets/Events/" + event.EventOrganiser + "/"  + event.EventName, 0777)

		file, err := os.Create("assets/Events/" + event.EventOrganiser + "/"  + event.EventName + "/" + event.PictureHeader.Filename)

		if err != nil {
			panic(err)
		}


		defer file.Close()
		io.Copy(file, event.Picture)
		// TODO: Change in production to actual ip
		event.EventPictureURL = Util.API_URL + "/api/v1/assets/Events/" + event.EventOrganiser + "/"  + event.EventName + "/" + event.PictureHeader.Filename
		// END IMAGE CREATION =====================================



		// Insert the event at the absolute end
		db.Exec("INSERT INTO events(event_name, event_start, event_end, event_location, event_organiser, "+
			"event_short_desc, "+ "event_long_desc, event_picture_url) "+
			"VALUES ($1, $2, $3, $4, $5 ,$6, $7, $8)", event.EventName, event.EventStart, event.EventEnd, event.EventLocation,
			event.EventOrganiser, event.EventShortDesc, event.EventLongDesc, event.EventPictureURL)

		return nil
	}

	return errors.New("user does not have sufficient privileges")
}



func validateDateTime(dateStart, format string) (bool, time.Time) {

	t, err := fmtdate.Parse(format, dateStart)
	if err != nil {
		return false, time.Time{}
	}

	return true, t

}
/*
	UTIL FUNCTIONS ============================
 */
 /* getIncomingEvent takes a request and returns the incoming event for that request
 		@params;
 			r *http.Request
 			user User.User
  */
func getIncomingEvent(r *http.Request, user User.User) (Event, error) {
	eventName := r.FormValue("Event_Name")
	eventEnd  := r.FormValue("Event_End")
	eventStart := r.FormValue("Event_Start")
	eventLocation := r.FormValue("Event_Location")
	eventOrganiser := user.Name
	eventShortDesc := r.FormValue("Event_Short_Desc")
	eventLongDesc := r.FormValue("Event_Long_Desc")


	if Util.CompoundIsset(eventName, eventEnd, eventLocation, eventOrganiser, eventLongDesc, eventStart, eventShortDesc) && Util.IsAdmin(user) {
		// Attain the image

		// Get the image uploading thing
		f, h, err := r.FormFile("Caption_Image")
		if err != nil {
			return Event{}, err
		}

		// Get the requested event
		requestedEvent := Event{
			EventName:      eventName,
			EventLocation:  eventLocation,
			EventOrganiser: eventOrganiser,
			EventShortDesc: eventShortDesc,
			EventLongDesc:  eventLongDesc,
			Picture: f,
			PictureHeader: h,
		}



		// Attain the date / time from the request ignore most of this part for maintainability
		// <+++++++++++++++++++++ DATE EXTRACTION START ++++++++++++++++++++++++++>
		err = parseDatesinto(&requestedEvent, eventStart, eventEnd)
		if err != nil {
			return Event{}, errors.New("invalid dates")
		}


		return requestedEvent, nil
	}

	return Event{}, errors.New("invalid request")

}
func parseDatesinto(requestedEvent *Event, eventStartstr string, eventEndstr string) error {
	// Event date
	pass, eventStart := validateDateTime(eventStartstr, "DD-MM-YYYY hh:mm")
	if eventStart.Unix() <= time.Now().Unix() || !pass {
		return errors.New("invalid date")
	}
	requestedEvent.EventStart = eventStart

	// Event time
	pass, eventEnd := validateDateTime(eventEndstr, "DD-MM-YYYY hh:mm")
	if !pass {
		return errors.New("invalid date")
	}
	requestedEvent.EventEnd = eventEnd


	return nil
}
/*
	END UTIL FUNCTIONS ============================
 */







// Http handler for authentication
/**
	@param event_location
	@param event_name
	@param event_organiser
	@param event_short_desc
	@param event_long_desc
	@param event_picture_url
	@param event_date_start
	@param event_date_end
	@param event_time_start
	@param event_time_end
**/
func CreateEventHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Connect to database
	Util.Conn("sensitive", "database", "admin")

	// Close the database at the end
	defer DB.DB.Close()


	// Get the user struct from an existing session and determin if they are allowed here
	allowed, user := Util.UserIsAllowed(r, w, "admin")
	if !allowed {
		QuickErrors.NotEnoughPrivledges(w)
		return
	}

	// Parse the URL
	r.ParseMultipartForm(1000000)

	event, err := getIncomingEvent(r, user)
	if err != nil {
		QuickErrors.MalformedRequest(w, "Missing parameters, user was invalid or the dates you provided were invalid")
		return
	}

	// Push the event
	err = Create(user, event, DB.DB)
	if err != nil {
		QuickErrors.MalformedRequest(w, "Event already exists")
		return
	}

	QuickErrors.OK(w)
}
