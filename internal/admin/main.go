package admin

import (
	"mynsb-api/internal/student"
	"database/sql"
	"encoding/json"
)

type Admin struct {
	Name        string
	Password    string
	Permissions []string
}

func AdminToUser(admin Admin) student.User {
	User := student.User{
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