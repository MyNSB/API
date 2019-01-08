package DB

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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
		return err
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	return nil
}
