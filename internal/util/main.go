package util

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"go/build"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/metakeule/fmtdate"
)

const APIURL = "https://mynsb.nsbvisions.com"
var TIMEZONE *time.Location


func init() {
	TIMEZONE, _ = time.LoadLocation("Australia/Sydney")
}



// ExistsString takes an array and an entry string and determines if that entry string resides in that array
func ExistsString(array []string, entry string) bool {
	for _, b := range array {
		if b == entry {
			return true
		}
	}
	return false
}


// ParseDate takes a date string, parses it and returns a time object
func ParseDate(date string) (time.Time, error){
	return fmtdate.Parse("DD-MM-YYYY", date)
}

// ParseDateTime takes a datetime string, parses it and returns a time object
func ParseDateTime(datetime string) (time.Time, error) {
	return fmtdate.Parse("DD-MM-YYYY hh:mm", datetime)
}


// GetGOPATH returns the system's first GOPATH variable
func GetGOPATH() string {
	// NOTE... This is a hack, as soon as you can find a better option... use it PLEASE!!!
	gopath := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	if len(gopath) == 0 {
		gopath = append(gopath, build.Default.GOPATH)
	} else if gopath[0] == "" {
		gopath[0] = build.Default.GOPATH
	}

	return gopath[0]
}


// NonNull determines if a string is empty, honestly just a code styling thing
func NonNull(thing string) bool {
	return thing != ""
}


// IsSet behaves much like php's isset and takes a list of variables and determins if they are null
func IsSet(vars ...string) bool {
	for _, varVal := range vars {
		if !NonNull(varVal) {
			return false
		}
	}
	return true
}


// HashString takes an input string and hashes it with the sha256 algorithm, unfortunately not the one you would see in your hashtables ;) too big :p
func HashString(toHash string) string {
	// Create a hasher
	h := sha256.New()
	// Write our data to it for hashing, in this case it is the string
	h.Write([]byte(toHash))
	// Create the final hash
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	return sha256Hash
}


// IsSubset takes two arrays of strings and determines if the second array is a subset of the first
func IsSubset(first, second []string) bool {
	for _, val := range first {
		if !ExistsString(second, val) {
			return false
		}
	}

	return true
}



// HTTPResponse function for returning array based responses
func HTTPResponse(status int, statusMessage string, body string, title string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"Status":{"Code": %d, "Status Message":"%s"},"Message": {"Title":"%s", "Body":[%s]}}`, status, statusMessage, title, body)))
}



// HTTPResponse function for returning object based responses
func HTTPResponseArr(status int, statusMessage string, body string, title string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"Status":{"Code": %d, "Status Message":"%s"},"Message": {"Title":"%s", "Body":"%s"}}`, status, statusMessage, title, body)))
}


// NumResults takes a db and a query and determines how many results were returned from that query
func NumResults(db *sql.DB, query string, args ...interface{}) (int, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	counter := 0
	for rows.Next() {
		counter += 1
	}
	rows.Next()
	return counter, nil
}
