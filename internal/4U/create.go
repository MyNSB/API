package FourU

import (
	"mynsb-api/internal/db"
	"mynsb-api/internal/util"
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"time"
	"mynsb-api/internal/quickerrors"
	"strconv"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/filesint"
	"fmt"
)

// Http handler for four u request publications
/*
	Handler's have minimal documentation
 */
func CreateFourUHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Connect to database
	db.Conn("admin")
	// Close database at the end
	defer db.DB.Close()

	// Determine if the student is allowed here and if not force them to leave
	allowed, _ := sessions.UserIsAllowed(r, w, "visions", "admin")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// Get the incoming article
	article, err := getIncomingArticle(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "You are missing fields, please check the API documentation")
	}

	// Push the article into the database
	err = CreateFourU(article, db.DB)
	if err != nil {
		quickerrors.MalformedRequest(w, "4U Article/Issue already exists")
		return
	}

	quickerrors.OK(w)
}

// HTTP handler for article creation
func CreateFourUArticleHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Connect to database
	db.Conn("admin")
	// Close database at the end
	defer db.DB.Close()

	// Parse the incoming article
	article, err := getIncomingArticleSpecific(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "You are missing fields, please check the API documentation")
	}

	// Finally push the new article
	err = Create4UArticles(Article{ArticleID: article.ArticleParentID}, article, db.DB)
	if err != nil {
		quickerrors.MalformedRequest(w, "Looks like that article already exists")
		return
	}

	quickerrors.OK(w)
}

/*
	CreateFourU takes an article and pushes that into the database it also saves the incoming image for that article returns an error on a failure
	@params;
		article FourU.Article
		db *sql.db
 */
func CreateFourU(article Article, db *sql.DB) error {
	// Check that the article does not currently exist in the database
	if count, _ := util.CheckCount(db, "SELECT * FROM four_u WHERE article_name = $1 and link = $2 and type = $3", article.ArticleName, article.Link, article.TypePost); count > 0 {
		return errors.New("article already exists")
	}

	// Hash the issue link for diversity
	sha := util.HashString(article.Link)
	// Save the article's image
	err := createImage(article, sha)
	if err != nil {
		return err
	}

	// Set the image url of the article
	article.ArticleImageUrl = fmt.Sprintf("%s/api/v1/assets/4U/%s/%s/%s/%s", util.APIURL, article.TypePost, sha, article.ArticleName, article.PictureHeader.Filename)

	// Finally push everything into the database
	db.Exec("INSERT INTO four_u (article_name, article_desc, article_publish_date, article_image_url, link, type) VALUES($1, $2, $3::DATE, $4, $5, $6)",
		article.ArticleName, article.ArticleDesc, article.ArticlePublishDate, article.ArticleImageUrl, article.Link, article.TypePost)

	// Return no error
	return nil
}

/* Create 4UArticle appends a new article to the 4U paper
		@params;
			article FourU.Article
			db *sql.db
 */

func Create4UArticles(article Article, uArticle FourUArticle, db *sql.DB) error {
	// Ensure that the article actually has an id
	if article.ArticleID == 0 {
		return errors.New("article is missing an ID")
	}

	// Check that the article does not exist
	if count, _ := util.CheckCount(db, "SELECT * FROM four_u_articles WHERE four_u_article = $1 AND four_u_article_name = $2", article.ArticleID, uArticle.ArticleName); count > 0 {
		// Looks like the article already exists
		return errors.New("article already exists")
	}

	// Insert teh article normally
	if _, err := db.Exec("INSERT INTO four_u_articles (four_u_article, page_start, four_u_article_name, four_u_article_desc) VALUES ($1, $2, $3, $4)",
		article.ArticleID, uArticle.ArticlePage, uArticle.ArticleName, uArticle.ArticleDesc); err != nil {
		return errors.New("unable to create new article")
	}

	return nil
}

/*
	@ UTIL FUNCTIONS ==================================================
 */
/*
	createImage takes an article and its save dir and creates the image that it requires
	@params;
		article FourU.Article
		imageSaveDir string

 */
func createImage(article Article, imageSaveDir string) error {
	// Create the article for the 4U directory
	fourUDir := fmt.Sprintf("/4U/%s/%s/%s/%s", article.TypePost, imageSaveDir, article.ArticleName, article.PictureHeader.Filename)
	file, err := filesint.CreateFile("assets", fourUDir)

	if err != nil {
		return errors.New("could not create image")
	}

	// Copy the actual image into the file
	io.Copy(file, article.Picture)

	return nil
}

/*
 	getIncomingArticle attains the incoming article from the request and returns an article and/or an error
 	@params;
 		r *http.Request

  */
func getIncomingArticle(r *http.Request) (Article, error) {
	// Parse incoming form
	r.ParseForm()
	// Get the article details
	articleName := r.Form.Get("Post_Name")
	articleDesc := r.Form.Get("Post_Desc")
	issuuLink := r.Form.Get("Link")
	typePost := r.Form.Get("Post_Type")

	// Article is invalid so throw an error
	if !(util.CompoundIsset(articleName, articleDesc, issuuLink, typePost) || !(typePost == "Article" || typePost == "Issue")) {
		return Article{}, errors.New("invalid article")
	}

	// Attain the pictures from multipart
	f, h, err := r.FormFile("Caption_Image")
	if err != nil {
		return Article{}, errors.New("caption image does not exist")
	}

	// Create temporary article
	article := Article{
		Picture:       f,
		PictureHeader: h,
		ArticleName:   articleName,
		ArticleDesc:   articleDesc,
		Link:          issuuLink,
		TypePost:      typePost,
	}
	// Set the publish date
	article.ArticlePublishDate = time.Now()
	return article, nil
}

func getIncomingArticleSpecific(r *http.Request) (FourUArticle, error) {
	// Parse that boy
	r.ParseForm()
	// Get the article details
	articleName := r.Form.Get("Article_Name")
	parentIssueID := r.Form.Get("Parent_ID")
	pageStart := r.Form.Get("Article_Page")
	articleDesc := r.Form.Get("Article_Desc")

	// Determine if everything is set
	if !(util.CompoundIsset(articleDesc, articleName, parentIssueID, pageStart)) {
		return FourUArticle{}, errors.New("invalid article")
	}

	// Convert some stuff
	parentID, _ := strconv.ParseInt(parentIssueID, 10, 64)
	pageStartN, _ := strconv.ParseInt(pageStart, 10, 64)

	// Construct the article
	article := FourUArticle{
		ArticleParentID: parentID,
		ArticlePage:     pageStartN,
		ArticleName:     articleName,
		ArticleDesc:     articleDesc,
	}

	return article, nil
}

/*
	@ END UTIL FUNCTIONS ==================================================
 */
