package user

import (
	"database/sql"
	"encoding/json"
)

type User struct {
	ID			int
	Name        string
	Password    string
	Permissions []string
}

// USed only for admin scanning
func (user *User) AdminScanFrom(rows *sql.Rows) {
	// Get the permissions
	var perms []byte
	rows.Scan(&user.ID, &user.Name, &user.Password, &perms)

	// Unmarshal the perms
	json.Unmarshal(perms, &user.Permissions)

}
