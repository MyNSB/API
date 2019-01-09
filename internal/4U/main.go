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

// Function merely takes a row and scans it into an article, minimal documentation required
func (article *Issue) ScanFrom(rows *sql.Rows) {
	rows.Scan(
		article.ID,
		&article.Name,
		&article.Desc,
		&article.PublishDate,
		&article.ImageUrl,
		&article.Link,
		&article.TypePost)
}

type Article struct {
	ID       int64
	ParentID int64
	Page     int64
	Name     string
	Desc     string
}
