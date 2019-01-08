package auth

import "mynsb-api/internal/student"

func userToMap(user student.User) map[string]interface{} {
	details := make(map[string]interface{})
	details["student"] = user.Name
	details["Password"] = user.Password
	details["Permissions"] = user.Permissions

	return details
}