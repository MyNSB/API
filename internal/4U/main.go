package fouru

import (
	"database/sql"
	"mime/multipart"
	"time"
)



// Data structures for holding the information regarding the 4U Paper
type Issue struct {
	ID            int64
	Name          string
	Desc          string
	PublishDate   time.Time
	Picture       multipart.File
	PictureHeader *multipart.FileHeader
	ImageUrl      string
	Link          string
	TypePost      string
}
// Articles are held within an issue however sometimes we might have a stand alone article
type Article struct {
	ID       int64
	ParentID int64
	Page     int64
	Name     string
	Desc     string
}




// ReplaceWith reads and sql.Rows object and pushes the information into an Issue object
// The assumption here is that the sql.Rows object is a result from a query involving the 4U Table
func (issue *Issue) ReplaceWith(rows *sql.Rows) {
	rows.Scan(
		&issue.ID,
		&issue.Name,
		&issue.Desc,
		&issue.PublishDate,
		&issue.ImageUrl,
		&issue.Link)
}