package fouru

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq" // This is required for sql/db
	"io"
	"mynsb-api/internal/db"
	"mynsb-api/internal/filesint"
	"mynsb-api/internal/quickerrors"
	"mynsb-api/internal/sessions"
	"mynsb-api/internal/util"
	"net/http"
	"strconv"
	"time"
)


/*	CreateIssue takes an issue and pushes that into the database it also saves the incoming image for that article returns an error on a failure
	@params;
		article FourU.Issue
		db *sql.db
*/
func CreateIssue(issue Issue, db *sql.DB) error {
	// Check that the issue does not currently exist in the database
	if count, _ := util.CheckCount(db, "SELECT * FROM four_u WHERE article_name = $1 and link = $2 and type = $3", issue.Name, issue.Link, issue.TypePost); count > 0 {
		return errors.New("issue already exists")
	}

	// Hash the issue link for diversity
	sha := util.HashString(issue.Link)
	// Save the issue's image
	err := createImage(issue, sha)
	if err != nil {
		return err
	}

	// Set the image url of the issue
	issue.ImageUrl = fmt.Sprintf("%s/api/v1/assets/4U/%s/%s/%s/%s", util.APIURL, issue.TypePost, sha, issue.Name, issue.PictureHeader.Filename)

	// Finally push everything into the database
	db.Exec("INSERT INTO four_u (article_name, article_desc, article_publish_date, article_image_url, link, type) VALUES($1, $2, $3::DATE, $4, $5, $6)",
		issue.Name, issue.Desc, issue.PublishDate, issue.ImageUrl, issue.Link, issue.TypePost)

	// Return no error
	return nil
}



/* CreateArticle appends a new article to the 4U paper
@params;
	article FourU.Issue
	db *sql.db
*/
func CreateArticle(parent Issue, article Article, db *sql.DB) error {
	// Ensure that the parent actually has an id
	if parent.ID == 0 {
		return errors.New("article is missing an ID")
	}

	// Check that the parent does not exist
	if count, _ := util.CheckCount(db, "SELECT * FROM four_u_articles WHERE four_u_article = $1 AND four_u_article_name = $2", parent.ID, article.Name); count > 0 {
		// Looks like the parent already exists
		return errors.New("article already exists")
	}

	// Insert teh parent normally
	if _, err := db.Exec("INSERT INTO four_u_articles (four_u_article, page_start, four_u_article_name, four_u_article_desc) VALUES ($1, $2, $3, $4)",
		parent.ID, article.Page, article.Name, article.Desc); err != nil {
		return errors.New("unable to create new article within issue")
	}

	return nil
}










/*
	@ UTIL FUNCTIONS ==================================================
*/
/*
	createImage takes an article and its save dir and creates the image that it requires
	@params;
		article FourU.Issue
		imageSaveDir string

*/
func createImage(issue Issue, imageSaveDir string) error {
	// Create the issue for the 4U directory
	fourUDir := fmt.Sprintf("/4U/%s/%s/%s", issue.TypePost, imageSaveDir, issue.Name)
	file, err := filesint.CreateFile("assets", fourUDir, issue.PictureHeader.Filename)

	if err != nil {
		fmt.Printf("%s", err.Error())
		return errors.New("could not create image")
	}

	// Copy the actual image into the file
	io.Copy(file, issue.Picture)

	return nil
}


/*
	getIncomingIssue attains the incoming article from the request and returns an article and/or an error
	@params;
		r *http.Request

*/
func getIncomingIssue(r *http.Request) (Issue, error) {
	// Get the issue details
	issueName := r.FormValue("Post_Name")
	issueDesc := r.FormValue("Post_Desc")
	issuuLink := r.FormValue("Link") // <--- Spelt this way on purpose

	// Issue is invalid so throw an error
	if !(util.CompoundIsset(issueName, issueDesc, issuuLink)) {
		return Issue{}, errors.New("invalid issue")
	}

	// Attain the pictures from multipart
	f, h, err := r.FormFile("Caption_Image")
	if err != nil {
		return Issue{}, errors.New("caption image does not exist")
	}

	// Create temporary issue
	issue := Issue{
		Picture:       f,
		PictureHeader: h,
		Name:          issueName,
		Desc:          issueDesc,
		Link:          issuuLink,
		TypePost:      "Issue",
	}
	// Set the publish date
	issue.PublishDate = time.Now()
	return issue, nil
}

func getIncomingArticle(r *http.Request) (Article, error) {
	// Get the article details
	articleName := r.Form.Get("Article_Name")
	parentIssueIDRaw := r.Form.Get("Parent_ID")
	pageStartRaw := r.Form.Get("Article_Page")
	articleDesc := r.Form.Get("Article_Desc")

	// Determine if everything is set
	if !(util.CompoundIsset(articleDesc, articleName, parentIssueIDRaw, pageStartRaw)) {
		return Article{}, errors.New("invalid article")
	}

	// Convert some stuff
	parentID, _ := strconv.ParseInt(parentIssueIDRaw, 10, 64)
	pageStart, _ := strconv.ParseInt(pageStartRaw, 10, 64)

	// Construct the article
	article := Article{
		ParentID: parentID,
		Page:     pageStart,
		Name:     articleName,
		Desc:     articleDesc,
	}

	return article, nil
}

/*
	@ END UTIL FUNCTIONS ==================================================
*/















// Http handler for four u request publications
/*
	Handler's have minimal documentation
*/

// CreateIssueHandler creates 4U issues
func CreateIssueHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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


	// Parse the incoming form
	r.ParseMultipartForm(1000000)

	// Get the incoming article
	issue, err := getIncomingIssue(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "You are missing fields, please check the API documentation")
		return
	}

	// Push the article into the database
	err = CreateIssue(issue, db.DB)
	if err != nil {
		panic(err)
		quickerrors.MalformedRequest(w, "4U Issue/Issue already exists")
		return
	}

	quickerrors.OK(w)
}

// CreateArticleHandler is a HTTP handler for article creation
func CreateArticleHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Connect to database
	db.Conn("admin")
	// Close database at the end
	defer db.DB.Close()

	// Parse the incoming article
	article, err := getIncomingArticle(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "You are missing fields, please check the API documentation")
	}

	// Finally push the new article
	err = CreateArticle(Issue{ID: article.ParentID}, article, db.DB)
	if err != nil {
		quickerrors.MalformedRequest(w, "Looks like that article already exists")
		return
	}

	quickerrors.OK(w)
}