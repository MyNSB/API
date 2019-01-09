package util

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"go/build"
	"net/http"
	"os"
)

const APIURL = "http://127.0.0.1"

/**
	Func Search:
		@param array AnyArray[]
		@param value

		returns Boolean
		True: Item exists
		False: Item doesn't exist
**/

func ExistsString(array []string, entry string) bool {
	for _, b := range array {
		if b == entry {
			return true
		}
	}
	return false
}

func GetGOPATH() string {
	// Get the GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	return gopath
}

/**
	Func isset:
		@param *string

		returns boolean
		True: Var exists
		False: Var doesn't exist
**/

func Isset(thing string) bool {
	return thing != ""
}

// Abstraction of isset
func CompoundIsset(vars ...string) bool {
	for _, varVal := range vars {
		if !Isset(varVal) {
			return false
		}
	}
	return true
}

func HashString(toHash string) string {
	// Create a hasher
	h := sha256.New()
	// Write our data to it for hashing, in this case it is the string
	h.Write([]byte(toHash))
	// Create the final hash
	sha256Hash := hex.EncodeToString(h.Sum(nil))

	return sha256Hash
}

func IsSubset(first, second []string) bool {
	for _, val := range first {
		if !ExistsString(second, val) {
			return false
		}
	}

	return true
}

// Function to remove all that ugly code error e.t.c
func Error(status int, statusMessage string, body string, title string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"Status":{"Code": %d, "Status Message":"%s"},"Message": {"Title":"%s", "Body":[%s]}}`, status, statusMessage, title, body)))
}

// Function to remove all that ugly code error e.t.c
func SolidError(status int, statusMessage string, body string, title string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{"Status":{"Code": %d, "Status Message":"%s"},"Message": {"Title":"%s", "Body":"%s"}}`, status, statusMessage, title, body)))
}

// Function to encrypt error messages for fixing later

// Function to return the number of returned rows it takes an actual query coz go is fucking stupid and will only let you iterate over the fucking set one fucking time!!!!
func CheckCount(db *sql.DB, query string, args ...interface{}) (int, error) {
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
