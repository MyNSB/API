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

/**
	Connection struct
	defines the basic behaviour and definition of a connection
**/
type Connection struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
}

// Database pointer
var DB *sql.DB

/**
	Func Connect:
		@param connection *Connection
		returns nil and just connects to the database
**/
func (connection *Connection) Connect() error {
	var err error
	// Connect to the database
	DB, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		connection.Host, connection.Port, connection.User, connection.Password, connection.DatabaseName))
	if err != nil {
		fmt.Printf(err.Error())
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	return nil
}

/*
	UTILITY FUNCTIONS
*/

/**
	Func getConnections:
		return db.Connection
**/
func getConnection() (Connection, error) {
	ConnectionDetails, err := filesint.DataDump("database", "/details.txt")
	if err != nil {
		return Connection{}, err
	}

	// Read the details from th file
	detailsArray := strings.Split(string(ConnectionDetails), ",")

	host := strings.Split(detailsArray[0], ":")[1]
	port := strings.Split(detailsArray[1], ":")[1]

	// Return the connection
	return Connection{
		Host: host,
		Port: port,
	}, nil
}

// Conn function for connection to database
// sensitiveLoc should look something like
func Conn(user string) error {
	// Connect to database as user
	connection, err := getConnection()
	if err != nil {
		panic(err)
	}
	// If err != nil why??
	connection.User = user

	// Attain the user password
	if pwd, err := filesint.DataDump("sensitive", fmt.Sprintf("/user pass/%s.txt", user)); err == nil {
		connection.Password = string(pwd)
	} else {
		return errors.New("could not authenticate")
	}

	connection.DatabaseName = "mynsb"

	err = connection.Connect()
	if err != nil {
		return errors.New("could not connect to database")
	}

	return nil
}
