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


// INSERTION FUNCTIONS

// issue.insertIntoDB takes an issue of the 4U Paper, checks whether if it already exists within the database and inserts it into the DB if it doesn't
func (issue Issue) insertIntoDB(db *sql.DB) error {
	// Determine if the issue has already been entered into the DB
	numInstances, _ := util.NumResults(db, "SELECT * FROM four_u WHERE article_name = $1 and link = $2", issue.Name, issue.Link);
	if numInstances != 0 {
		return errors.New("issue already exists")
	}

	// Copy the image attached with the issue into a folder stored onto the server
	// Save the issue's image
	imageLocation, err := saveImage(issue)
	if err != nil {
		return err
	}

	// Set the Issue's image URL to be the location of the saved image on our server
	issue.ImageUrl = fmt.Sprintf("%s/api/v1/assets/%s", util.APIURL, imageLocation)

	// Push everything into the DB
	db.Exec("INSERT INTO four_u (article_name, article_desc, article_publish_date, article_image_url, link, type) VALUES($1, $2, $3::DATE, $4, $5, $6)",
		issue.Name, issue.Desc, issue.PublishDate, issue.ImageUrl, issue.Link)

	return nil
}


// article.insertIntoDB inserts a "4U Article" into the DB and like issue.insertIntoDB it also checks weather is already exists
func (article Article) insertIntoDB(parent Issue, db *sql.DB) error {
	// Determine if the article has already been entered into the DB
	numInstances, _ := util.NumResults(db, "SELECT * FROM four_u_articles WHERE four_u_article = $1 AND four_u_article_name = $2", parent.ID, article.Name)
	if numInstances != 0 {
		return errors.New("article already exists")
	}

	// Insert into the DB, if there is a failure that generally implies that the parent was invalid
	if _, err := db.Exec("INSERT INTO four_u_articles (four_u_article, page_start, four_u_article_name, four_u_article_desc) VALUES ($1, $2, $3, $4)",
		parent.ID, article.Page, article.Name, article.Desc); err != nil {
		return errors.New("unable to insert article into DB, most likely because the parent provided did not exist")
	}

	return nil
}












// UTILITY FUNCTIONS

// saveImage takes an issue of the 4U Paper and saves the image associated with it to disk
func saveImage(issue Issue) (string, error) {
	// We take the string's hash to be the directory we will be using to save the issue
	// The reason why we are hashing the link is as they will generally be unique from issue to issue and that reduces the number of possible hash collisions
	imageSaveDir := util.HashString(issue.Link)

	// Create the directory that will be used to save the image
	fourUDir := fmt.Sprintf("/4U/%s/%s/%s", issue.TypePost, imageSaveDir, issue.Name)
	file, err := filesint.CreateFile("assets", fourUDir, issue.PictureHeader.Filename)
	if err != nil {
		return "", errors.New("could not create image")
	}

	// Copy the actual image into the file object
	io.Copy(file, issue.Picture)

	return fourUDir, nil
}


// getIncomingIssue parses the issue being sent by the user via HTTP into an actual issue object
func getIncomingIssue(r *http.Request) (Issue, error) {
	r.ParseMultipartForm(1000000)

	issueName := r.FormValue("Post_Name")
	issueDesc := r.FormValue("Post_Desc")
	issuuLink := r.FormValue("Link") // <--- Spelt this way on purpose

	if !(util.IsSet(issueName, issueDesc, issuuLink)) {
		return Issue{}, errors.New("not all parameters have been provided")
	}

	// Read the attached image using the multipart package
	f, h, err := r.FormFile("Caption_Image")
	if err != nil {
		return Issue{}, errors.New("caption image does not exist")
	}

	// Return the parsed issue
	return Issue{
		Picture:       f,
		PictureHeader: h,
		Name:          issueName,
		Desc:          issueDesc,
		Link:          issuuLink,
		PublishDate:   time.Now().In(util.TIMEZONE),
		TypePost:      "Issue",
	}, nil
}


// getIncomingArticle parses the article being sent by the user via HTTP into an actual article object
func getIncomingArticle(r *http.Request) (Article, error) {
	r.ParseMultipartForm(1000000)

	articleName := r.Form.Get("Name")
	parentIssueIDRaw := r.Form.Get("Parent_ID")
	pageStartRaw := r.Form.Get("Page")
	articleDesc := r.Form.Get("Desc")

	if !(util.IsSet(articleDesc, articleName, parentIssueIDRaw, pageStartRaw)) {
		return Article{}, errors.New("not all parameters have been provided")
	}

	// Convert the text provided by the user into integers
	parentID, _ := strconv.ParseInt(parentIssueIDRaw, 10, 64)
	pageStart, _ := strconv.ParseInt(pageStartRaw, 10, 64)

	// Return the parsed article
	return Article{
		ParentID: parentID,
		Page:     pageStart,
		Name:     articleName,
		Desc:     articleDesc,
	}, nil
}












// HTTP HANDLERS

// IssueCreationHandler creates 4U issues
func IssueCreationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// connect to database
	db.Conn("admin")
	defer db.DB.Close()

	// Determine if the user is allowed here and if not force them to leave
	allowed, _ := sessions.IsUserAllowed(r, w, "visions", "admin")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// Get the incoming issue
	issue, err := getIncomingIssue(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "You are missing fields, please check the API documentation")
		return
	}

	// Push the issue into the database
	err = issue.insertIntoDB(db.DB)
	if err != nil {
		panic(err)
		quickerrors.MalformedRequest(w, "4U Issue/Issue already exists")
		return
	}

	quickerrors.OK(w)
}


// ArticleCreationHandler is a HTTP handler for article creation
func ArticleCreationHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// connect to database as an admin
	db.Conn("admin")
	defer db.DB.Close()

	// Determine if the user is allowed here and if not force them to leave
	allowed, _ := sessions.IsUserAllowed(r, w, "visions", "admin")
	if !allowed {
		quickerrors.NotEnoughPrivileges(w)
		return
	}

	// Attain the incoming article
	article, err := getIncomingArticle(r)
	if err != nil {
		quickerrors.MalformedRequest(w, "You are missing fields, please check the API documentation")
		return
	}

	// Push the article into the database
	parentIssue := Issue{ID: article.ParentID}
	err = article.insertIntoDB(parentIssue, db.DB)
	if err != nil {
		quickerrors.MalformedRequest(w, "Looks like that article already exists in our DB")
		return
	}

	// All Clear :)
	quickerrors.OK(w)
}