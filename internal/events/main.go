package events

import (
	"database/sql"
	"mime/multipart"
	"time"
)

type Event struct {
	ID            int64
	Name          string
	Start         time.Time
	End           time.Time
	Location      string
	Organiser     string
	ShortDesc     string
	LongDesc      string
	Picture       multipart.File
	PictureHeader *multipart.FileHeader
	PictureURL    string
}

func (event *Event) ScanFrom(rows *sql.Rows) {
	rows.Scan(
		&event.ID,
		&event.Name,
		&event.Start,
		&event.End,
		&event.Location,
		&event.Organiser,
		&event.ShortDesc,
		&event.LongDesc,
		&event.PictureURL)
}
