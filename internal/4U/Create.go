package FourU

import (
	"DB"
	"Util"
	"database/sql"
	"errors"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"os"
	"time"
	"QuickErrors"
	"strconv"
)





// Http handler for four u request publications
/*
	Handler's have minimal documentation
 */
func CreateFourUHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Connect to database
	Util.Conn("sensitive", "database", "admin")
	// Close database at the end
	defer DB.DB.Close()

	// Determine if the user is allowed here and if not force them to leave
	allowed, _ := Util.UserIsAllowed(r, w, "visions", "admin")
	if !allowed {
		QuickErrors.NotEnoughPrivledges(w)
		return
	}

	// Get the incoming article
	article, err := getIncomingArticle(r)
	if err != nil {
		QuickErrors.MalformedRequest(w, "You are missing fields, please check the API documentation")
	}

	// Push the article into the database
	err = CreateFourU(article, DB.DB)
	if err != nil {
		QuickErrors.MalformedRequest(w, "4U Article/Issue already exists")
		return
	}

	QuickErrors.OK(w)
}

// HTTP handler for article creation
func CreateFourUArticleHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Connect to database
	Util.Conn("sensitive", "database", "admin")
	// Close database at the end
	defer DB.DB.Close()

	// Parse the incomming article
	article, err := getIncomingArticleSpecific(r)
	if err != nil {
		QuickErrors.MalformedRequest(w, "You are missing fields, please check the API documentation")
	}

	// Finally push the new article
	err = Create4UArtciles(Article{ArticleID: article.ArticleParentID}, article, DB.DB)
	if err != nil {
		QuickErrors.MalformedRequest(w, "Looks like that article already exists")
		return
	}

	QuickErrors.OK(w)
}



/*
	CreateFourU takes an article and pushes that into the database it also saves the incoming image for that article returns an error on a failure
	@params;
		article FourU.Article
		db *sql.DB
 */
func CreateFourU(article Article, db *sql.DB) error {
	// Check that the article does not currently exist in the database
	if count, _ := Util.CheckCount(db, "SELECT * FROM four_u WHERE article_name = $1 and link = $2 and type = $3", article.ArticleName, article.Link, article.TypePost); count > 0 {
		return errors.New("article already exists")
	}

	// Hash the issue link for diversity
	sha := Util.HashString(article.Link)
	// Save the article's image
	err := createImage(article, sha)
	if err != nil {
		return err
	}

	// Set the image url of the article
	article.ArticleImageUrl = Util.API_URL + "/api/v1/assets/4U/" + article.TypePost + "/" + sha + "/" + article.ArticleName + "/" + article.PictureHeader.Filename

	// Finally push everything into the database
	db.Exec("INSERT INTO four_u (article_name, article_desc, article_publish_date, article_image_url, link, type) VALUES($1, $2, $3::date, $4, $5, $6)",
		article.ArticleName, article.ArticleDesc, article.ArticlePublishDate, article.ArticleImageUrl, article.Link, article.TypePost)


	// Return no error
	return nil
}




/* Create 4UArtcile appends a new artcile to the 4U paper
		@params;
			article FourU.Artcile
			db *sql.db
 */

 func Create4UArtciles(article Article, uArticle FourUArticle, db *sql.DB) error {
 	// Ensure that the article actually has an id
 	if article.ArticleID == 0 {
 		return errors.New("Article is missing an ID")
	}

	// Check that the article does not exist
	if count, _ := Util.CheckCount(db, "SELECT * FROM four_u_articles WHERE four_u_article = $1 AND four_u_article_name = $2", article.ArticleID, uArticle.ArticleName); count > 0 {
		// Looks like the article already exists
		return errors.New("Article already exists")
	}


	// Insert teh article normall
	if _, err := db.Exec("INSERT INTO four_u_articles (four_u_article, page_start, four_u_article_name, four_u_article_desc) VALUES ($1, $2, $3, $4)",
		article.ArticleID, uArticle.ArticlePage, uArticle.ArticleName, uArticle.ArticleDesc); err != nil {
			return errors.New("Unable to create new article")
	}


	return nil;
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
	if _, err := os.Stat("assets/4U/" + article.TypePost + "/" + imageSaveDir); os.IsNotExist(err) {
		os.Mkdir("assets/4U/" + article.TypePost + "/" + imageSaveDir, 0777)
	}

	os.Mkdir("assets/4U/" + article.TypePost + "/" + imageSaveDir + "/" + article.ArticleName, 0777)

	// Create image and image url
	// Create a temp copy
	file, err := os.Create("assets/4U/" + article.TypePost + "/" + imageSaveDir + "/" + article.ArticleName + "/" + article.PictureHeader.Filename)
	if err != nil {
		return err
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
	if !(Util.CompoundIsset(articleName, articleDesc, issuuLink, typePost) || !(typePost == "Article" || typePost == "Issue")) {
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
		Link:     	   issuuLink,
		TypePost:	   typePost,
	}
	// Set the publish date
	article.ArticlePublishDate = time.Now()
	return article, nil
}


func getIncomingArticleSpecific(r *http.Request) (FourUArticle, error) {
	// Parse that bish
	r.ParseForm()
	// Get the article details
	articleName   := r.Form.Get("Article_Name")
	parentIssueID := r.Form.Get("Parent_ID")
	pageStart	  := r.Form.Get("Article_Page")
	articleDesc	  := r.Form.Get("Article_Desc")


	// Determine if everything is set
	if !(Util.CompoundIsset(articleDesc, articleName, parentIssueID, pageStart)) {
		return FourUArticle{}, errors.New("Invalid article")
	}

	// Convert some stuff
	parentID, _   := strconv.ParseInt(parentIssueID, 10, 64)
	pageStartN, _ := strconv.ParseInt(pageStart, 10, 64)


	// Construct the article
	article := FourUArticle{
		ArticleParentID: parentID,
		ArticlePage: pageStartN,
		ArticleName: articleName,
		ArticleDesc: articleDesc,
	}

	return article, nil
}

/*
	@ END UTIL FUNCTIONS ==================================================
 */
