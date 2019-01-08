package events

import (
	"mime/multipart"
	"time"
	"database/sql"
)

type Event struct {
	EventName       string
	EventID         int64
	EventStart      time.Time
	EventEnd        time.Time
	EventLocation   string
	EventOrganiser  string
	EventShortDesc  string
	EventLongDesc   string
	Picture         multipart.File
	PictureHeader   *multipart.FileHeader
	EventPictureURL string
}


func (event *Event) ScanFrom(rows *sql.Rows) {
	var eventId 			int64
	var eventName 			string
	var eventStart 			time.Time
	var eventEnd 			time.Time
	var eventLocation 		string
	var eventOrganiser 		string
	var eventShortDesc 		string
	var eventLongDesc 		string
	var eventPictureUrl 	string

	rows.Scan(&eventId, &eventName, &eventStart, &eventEnd, &eventLocation, &eventOrganiser, &eventShortDesc, &eventLongDesc, &eventPictureUrl)

	event.EventID 			= eventId
	event.EventName 		= eventName
	event.EventStart 		= eventStart
	event.EventEnd 			= eventEnd
	event.EventLocation 	= eventLocation
	event.EventOrganiser 	= eventOrganiser
	event.EventShortDesc	= eventShortDesc
	event.EventLongDesc		= eventLongDesc
	event.EventPictureURL	= eventPictureUrl
}
