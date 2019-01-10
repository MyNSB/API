package auth

import "mynsb-api/internal/student"
import "mynsb-api/internal/jwt"

func userToJWTData(user student.User) jwt.JWTData {
	jwtData := jwt.JWTData{
		User: user.Name,
		Password: user.Password,
		Permissions: user.Permissions,
	}

	return jwtData
}
