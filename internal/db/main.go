package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"mynsb-api/internal/filesint"
	"strings"
)

// Whole thing is pre much just an abstraction

// Connection struct is a struct that defines the connection details required for DB authentication
type Connection struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
}

// Database pointer
var DB *sql.DB




// CONNECTION FUNCTIONS


// connect takes a connection configuration and connects to the psql database with it, the database connection is reflected within the db.DB object
func (connection Connection) connect() error {

	var err error

	// connect to the database
	DB, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		connection.Host, connection.Port, connection.User, connection.Password, connection.DatabaseName))
	if err != nil {
		return errors.New("could not connect to the database")
	}

	if err = DB.Ping(); err != nil {
		return errors.New("could not connect to the database")
	}

	return nil
}










// UTILITY FUNCTIONS

// getConnectionDetails returns a connection configuration based off the current database config in database/details.txt
func getConnectionDetails() (Connection, error) {

	connectionDetails, err := filesint.DataDump("database", "/details.txt")
	if err != nil {
		return Connection{}, err
	}

	// Attain the configuration through separation of delimiters
	detailsArray := strings.Split(string(connectionDetails), ",")
	host := strings.Split(detailsArray[0], ":")[1]
	port := strings.Split(detailsArray[1], ":")[1]


	return Connection{
		Host: host,
		Port: port,
	}, nil
}


// getUserPswd returns the pswd of the user inputted, this password is read off the /user pass/student.txt file
func getUserPswd(user string) string {
	if pwd, err := filesint.DataDump("sensitive", fmt.Sprintf("/user pass/%s.txt", user)); err == nil {
		return string(pwd)
	}

	return ""
}


// Conn takes a string signifying the type of user you wish to authenticate as and connects to it via the .connect() function
func Conn(user string) error {
	// connect to database as user
	connectionDetails, _ := getConnectionDetails()
	connectionDetails.User = user
	connectionDetails.Password = getUserPswd(user)
	connectionDetails.DatabaseName = "mynsb"

	err := connectionDetails.connect()
	if err != nil {
		return errors.New("could not connect to database")
	}

	return nil
}
