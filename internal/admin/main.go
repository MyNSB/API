package admin

import (
	"database/sql"
	"encoding/json"
	"mynsb-api/internal/student"
)

type Admin struct {
	ID			int
	Name        string
	Password    string
	Permissions []string
}

func ToUser(admin Admin) student.User {
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
	rows.Scan(&admin.ID, &admin.Name, &admin.Password, &adminPermissions)

	// Unmarshal the perms
	json.Unmarshal(adminPermissions, &admin.Permissions)

}
