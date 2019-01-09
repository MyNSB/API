package fouru

import (
	"database/sql"
	"mime/multipart"
	"time"
)

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

// Parses SQL query result into current Issue object
func (issue *Issue) ScanFrom(rows *sql.Rows) {
	rows.Scan(
		&issue.ID,
		&issue.Name,
		&issue.Desc,
		&issue.PublishDate,
		&issue.ImageUrl,
		&issue.Link,
		&issue.TypePost)
}

type Article struct {
	ID       int64
	ParentID int64
	Page     int64
	Name     string
	Desc     string
}
