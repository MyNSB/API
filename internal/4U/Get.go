package FourU

import (
	"DB"
	"Util"
	"database/sql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/metakeule/fmtdate"
	"net/http"
	"QuickErrors"
	"errors"
)




// Http handler for four u article requests
/*
	Handler's have minimal documentation
*/
func GetFourUHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// Start up database
	Util.Conn("sensitive", "database", "student")
	// Close the database at the end
	defer DB.DB.Close()

	// Determine the request type for the incoming request
	requestType, params := determineRequestType(r)


	// Perform the request in regards to parameters
	switch {

		case requestType: // Request type = all {check function documentation}
			// Perform the get all function
			articles := GetAll(DB.DB)
			// Encode result to json
			bytes, _ := json.Marshal(articles)

			// Encode result and return it to the user
			Util.Error(200, "OK", string(bytes), "Response", w)

		default: // Request type = between {check function documentation}

			// Perform it
			res, err := GetBetween(params, DB.DB)
			if err != nil {
				QuickErrors.InteralServerError(w)
				return
			}

			// encode and return the result
			bytes, _ := json.Marshal(res)
			Util.Error(200, "Ok", string(bytes), "Response", w)
	}
}





// Functions to retrieve articles ==================

/*
	GetAll returns all 4U articles currently in the database, once mynsb grows this will have to shrink to the past year but for now it can stay as the entire db
	@params;
		db *sql.DB
*/
func GetAll(db *sql.DB) []Article {
	article, _ := performRequest(db, "SELECT * FROM four_u")
	return article
}


/*
	Getbetween returns all function between specified times
	@params;
		times map[string]string
		db *sql.DB
 */
func GetBetween(times map[string]string, db *sql.DB) ([]Article, error) {
	// Convert the start and end to actual time values
	// Convert into dates
	start, err := fmtdate.Parse("DD-MM-YYYY", times["start"])
	if err != nil {
		return []Article{}, errors.New("invalid date format")
	}

	end, err := fmtdate.Parse("DD-MM-YYYY", times["end"])
	if err != nil {
		return []Article{}, errors.New("invalid date format")
	}

	return performRequest(db, "SELECT * FROM four_u WHERE article_publish_date BETWEEN $1::TIMESTAMP AND $2::TIMESTAMP", start, end)
}

// ================================





/*
	@ UTIL FUNCTIONS ==================================================
 */
 /*
 	determineRequestType determines the type of request for the incoming http.request, true for all false for between, it also returns parameters
 	@params;
 		r *http.Request
 */
func determineRequestType(r *http.Request) (bool, map[string]string) {
	// Request Type
	var typeReq bool

	startTXT := r.URL.Query().Get("Start")
	endTXT := r.URL.Query().Get("End")

	// Determine type of request based on parsed parameters
	if Util.CompoundIsset(startTXT, endTXT) {
		typeReq = false
	} else {
		typeReq = true
	}

	// Map to be returned
	toReturn := make(map[string]string)
	toReturn["start"] = startTXT
	toReturn["end"] = endTXT

	return typeReq, toReturn
}


/*
	performRequest performs a request given a query and some arguments it returns an array or articles and a possible error
	@params;
		db *sql.DB
		query string
		args ...interface{}
 */
func performRequest(db *sql.DB, query string, args ...interface{}) ([]Article, error) {
	// Array that will be returned
	var result []Article

	// Get everything
	res, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	// Close res at the end
	defer res.Close()

	// Iterate over result set
	for res.Next() {
		article := Article{}
		// Scan the rows into the article
		article.ScanFrom(res)

		// Append article to array to be returned
		result = append(result, article)
	}

	return result, nil
}
/*
	@ END UTIL FUNCTIONS ==================================================
 */