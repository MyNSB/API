package user

import (
	"database/sql"
	"encoding/json"
	"mynsb-api/internal/util"
)

type User struct {
	ID			int
	Name        string
	Password    string
	Permissions []string
}


// Used to scan an SQL query into a user object, the assumption is that this query is from the admin table
func (user *User) ScanSQLIntoAdmin(rows *sql.Rows) {
	// Get the permissions
	var perms []byte
	rows.Scan(&user.ID, &user.Name, &user.Password, &perms)

	// Unmarshal the perms
	json.Unmarshal(perms, &user.Permissions)

}


// .IsAdmin takes a user entity and determines if the current user is an admin
func (user *User) IsAdmin() bool {
	return util.ExistsString(user.Permissions, "admin")
}