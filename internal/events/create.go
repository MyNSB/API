package events

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"mynsb-api/internal/db"
	"mynsb-api/internal/filesint"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
)



// EVENT CREATION FUNCTIONS


// create takes an event and a database and inserts the event into the database
func (event Event) create(db *sql.DB) error {

	// Determine if the event exists in our DB or not
	numEvents, _ := util.NumResults(db, "SELECT * FROM events WHERE event_name = $1 AND event_organiser = $2", event.Name, event.Organiser)
	if numEvents != 0 {
		return errors.New("event already exists")
	}

	// create the image associated with the requested event
	eventPictureLoc := fmt.Sprintf("/events/%s/%s", event.Organiser, event.Name)
	file, err := filesint.CreateFile("assets", eventPictureLoc, event.PictureHeader.Filename)
	defer file.Close()
	if err != nil {
		return err
	}

	// Copy the picture into the file pointer
	io.Copy(file, event.Picture)
	event.PictureURL = fmt.Sprintf("%s/api/v1/assets%s", util.APIURL, eventPictureLoc)

	// Insert the event into the database
	db.Exec("INSERT INTO events(event_name, event_start, event_end, event_location, event_organiser, "+
		"event_short_desc, "+"event_long_desc, event_picture_url) "+
			"VALUES ($1, $2, $3, $4, $5 ,$6, $7, $8)", event.Name, event.Start, event.End, event.Location,
			event.Organiser, event.ShortDesc, event.LongDesc, event.PictureURL)

	return nil

}












// UTILITY FUNCTIONS

// parseIncomingEvent takes a http request and parses the incoming event that is being sent by the user
func parseIncomingEvent(r *http.Request, organiser string) (Event, error) {

	r.ParseMultipartForm(1000000)

	eventName := r.FormValue("Name")
	eventEndRAW := r.FormValue("End")
	eventStartRAW := r.FormValue("Start")
	eventLocation := r.FormValue("Location")
	eventOrganiser := organiser
	eventShortDesc := r.FormValue("Short_Desc")
	eventLongDesc := r.FormValue("Long_Desc")
	// Check that the variables are actually set
	if !(util.IsSet(eventName, eventEndRAW, eventLocation, eventOrganiser, eventLongDesc, eventStartRAW, eventShortDesc)) {
		return Event{}, errors.New("user is mising paramaters")
	}



	// Attain the attached image with the request
	f, h, err := r.FormFile("Caption_Image")
	if err != nil {
		return Event{}, err
	}


	// Parse the event datetimes
	eventStart, parseErrorOne := util.ParseDateTime(eventStartRAW)
	eventEnd, parseErrorTwo := util.ParseDateTime(eventEndRAW)
	if parseErrorTwo != nil || parseErrorOne != nil {
		return Event{}, errors.New("could not parse date")
	}


	return Event{
		Name:          eventName,
		Location:      eventLocation,
		Organiser:     eventOrganiser,
		ShortDesc:     eventShortDesc,
		LongDesc:      eventLongDesc,
		Picture:       f,
		PictureHeader: h,
		Start:		   eventStart,
		End:		   eventEnd,
	}, nil
}












// HTTP HANDLERS

// EventCreationHandler is a http handler that handles creation requests from users
func EventCreationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	db.Conn("admin")
	defer db.DB.Close()

	// Check user privledges
	allowed, currUser := sessions.IsUserAllowed(r, w, "admin")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}


	event, err := parseIncomingEvent(r, currUser.Name)
	if err != nil {
		quickerrors.MalformedRequest(w, "Missing parameters, user was invalid or the dates you provided were invalid")
		return
	}

	err = event.create(db.DB)
	if err != nil {
		quickerrors.MalformedRequest(w, "Event already exists")
		return
	}

	quickerrors.OK(w)
}
