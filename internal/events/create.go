package events

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/metakeule/fmtdate"
	"io"
	"mynsb-api/internal/db"
	"mynsb-api/internal/filesint"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/user"
	"mynsb-api/internal/util"
	"net/http"
	"time"
)

/**
	Func CreateEvent:
		@param event Event

		returns nothing and just creates an event
**/
func Create(user user.User, event Event, db *sql.DB) error {

	// Determine if the user is an admin
	if user.IsAdmin(user) {

		// EVENT EXISTS DETERMINATION ===============================
		// Check that the event does not exist
		// Determine the count using the util function
		if count, _ := util.NumResults(db, "SELECT * FROM events WHERE event_name = $1 AND event_organiser = $2", event.Name, event.Organiser); count > 0 {
			return errors.New("event already exists")
		}
		// EVENT EXISTS DETERMINATION END ===============================

		// IMAGE CREATION =========================================
		eventPictureLoc := fmt.Sprintf("/events/%s/%s", event.Organiser, event.Name)
		file, err := filesint.CreateFile("assets", eventPictureLoc, event.PictureHeader.Filename)

		if err != nil {
			panic(err)
		}

		defer file.Close()
		io.Copy(file, event.Picture)
		event.PictureURL = fmt.Sprintf("%s/api/v1/assets%s", util.APIURL, eventPictureLoc)
		// END IMAGE CREATION =====================================

		// Insert the event at the absolute end
		db.Exec("INSERT INTO events(event_name, event_start, event_end, event_location, event_organiser, "+
			"event_short_desc, "+"event_long_desc, event_picture_url) "+
			"VALUES ($1, $2, $3, $4, $5 ,$6, $7, $8)", event.Name, event.Start, event.End, event.Location,
			event.Organiser, event.ShortDesc, event.LongDesc, event.PictureURL)

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
	user user.user
*/
func getIncomingEvent(r *http.Request, user user.User) (Event, error) {
	eventName := r.FormValue("Event_Name")
	eventEnd := r.FormValue("Event_End")
	eventStart := r.FormValue("Event_Start")
	eventLocation := r.FormValue("Event_Location")
	eventOrganiser := user.Name
	eventShortDesc := r.FormValue("Event_Short_Desc")
	eventLongDesc := r.FormValue("Event_Long_Desc")

	if util.IsSet(eventName, eventEnd, eventLocation, eventOrganiser, eventLongDesc, eventStart, eventShortDesc) {
		// Attain the image

		// Get the image uploading thing
		f, h, err := r.FormFile("Caption_Image")
		if err != nil {
			return Event{}, err
		}

		// Get the requested event
		requestedEvent := Event{
			Name:          eventName,
			Location:      eventLocation,
			Organiser:     eventOrganiser,
			ShortDesc:     eventShortDesc,
			LongDesc:      eventLongDesc,
			Picture:       f,
			PictureHeader: h,
		}

		// Attain the date / time from the request ignore most of this part for maintainability
		// <+++++++++++++++++++++ DATE EXTRACTION START ++++++++++++++++++++++++++>
		err = parseDatesInto(&requestedEvent, eventStart, eventEnd)
		if err != nil {
			return Event{}, errors.New("invalid dates")
		}

		return requestedEvent, nil
	}

	return Event{}, errors.New("invalid request")

}

func parseDatesInto(requestedEvent *Event, eventStartStr string, eventEndStr string) error {
	// Event date
	pass, eventStart := validateDateTime(eventStartStr, "DD-MM-YYYY hh:mm")
	if eventStart.Unix() <= time.Now().Unix() || !pass {
		return errors.New("invalid date")
	}
	requestedEvent.Start = eventStart

	// Event time
	pass, eventEnd := validateDateTime(eventEndStr, "DD-MM-YYYY hh:mm")
	if !pass {
		return errors.New("invalid date")
	}
	requestedEvent.End = eventEnd

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
func CreateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Connect to database
	db.Conn("admin")

	// Close the database at the end
	defer db.DB.Close()

	// Get the user struct from an existing session and determine if they are allowed here
	allowed, user := sessions.IsUserAllowed(r, w, "admin")
	if !allowed || !user.IsAdmin(user) {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// Parse the URL
	r.ParseMultipartForm(1000000)

	event, err := getIncomingEvent(r, user)
	if err != nil {
		quickerrors.MalformedRequest(w, "Missing parameters, user was invalid or the dates you provided were invalid")
		return
	}

	// Push the event
	err = Create(user, event, db.DB)
	if err != nil {
		quickerrors.MalformedRequest(w, "Event already exists")
		return
	}

	quickerrors.OK(w)
}
