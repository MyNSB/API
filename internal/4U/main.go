package FourU

import (
	"mime/multipart"
	"time"
	"database/sql"
)

type Article struct {
	ArticleID		   int64
	ArticleName        string
	ArticleDesc        string
	ArticlePublishDate time.Time
	Picture            multipart.File
	PictureHeader      *multipart.FileHeader
	ArticleImageUrl    string
	Link          	   string
	TypePost		   string
}



// Function merely takes a row and scans it into an article, minimal documentation required
func (article *Article) ScanFrom(rows *sql.Rows) {
	var articleID			int64
	var articleName 		string
	var articleDesc 		string
	var articlePublishDate 	time.Time
	var articleImageUrl 	string
	var issueLink 			string
	var typeReq				string

	rows.Scan(&articleID, &articleName, &articleDesc, &articlePublishDate, &articleImageUrl, &issueLink, &typeReq)

	// Insert the values now
	article.ArticleID  			= 	articleID
	article.ArticleName 		=   articleName
	article.ArticleDesc  		=   articleDesc
	article.ArticleImageUrl 	= 	articleImageUrl
	article.ArticlePublishDate  = 	articlePublishDate
	article.Link       	 		=	issueLink
	article.TypePost		 	=	typeReq
}


type FourUArticle struct {
	ArticleID		int64
	ArticleParentID	int64
	ArticlePage		int64
	ArticleName		string
	ArticleDesc		string
}