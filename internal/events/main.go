package events

import (
	"database/sql"
	"mime/multipart"
	"time"
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
	rows.Scan(
		&event.EventID,
		&event.EventName,
		&event.EventStart,
		&event.EventEnd,
		&event.EventLocation,
		&event.EventOrganiser,
		&event.EventShortDesc,
		&event.EventLongDesc,
		&event.EventPictureURL)
}
