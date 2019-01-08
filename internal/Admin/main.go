package Admin

import (
	"User"
	"database/sql"
	"encoding/json"
)

type Admin struct {
	Name        string
	Password    string
	Permissions []string
}

func AdminToUser(admin Admin) User.User {
	User := User.User{
		Name:        admin.Name,
		Password:    admin.Password,
		Permissions: admin.Permissions,
	}

	return User
}


func (admin *Admin) ScanFrom(rows *sql.Rows) {
	// Get the permissions
	var adminPermissions []byte
	rows.Scan(&adminPermissions)

	// Unmarshal the perms
	var perms []string
	json.Unmarshal([]byte(adminPermissions), &perms)

	admin.Permissions = perms
}