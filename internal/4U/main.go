package FourU

import (
	"mime/multipart"
	"time"
	"database/sql"
)

type Article struct {
	ArticleID          int64
	ArticleName        string
	ArticleDesc        string
	ArticlePublishDate time.Time
	Picture            multipart.File
	PictureHeader      *multipart.FileHeader
	ArticleImageUrl    string
	Link               string
	TypePost           string
}

// Parses an SQL query result into the current Article object
func (article *Article) ScanFrom(rows *sql.Rows) {
	rows.Scan(&article.ArticleID,
	          &article.ArticleName,
	          &article.ArticleDesc,
	          &article.ArticlePublishDate,
	          &article.ArticleImageUrl,
	          &article.Link,
	          &article.TypePost)
}

type FourUArticle struct {
	ArticleID       int64
	ArticleParentID int64
	ArticlePage     int64
	ArticleName     string
	ArticleDesc     string
}
