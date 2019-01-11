package auth

import "mynsb-api/internal/user"
import "mynsb-api/internal/jwt"

func userToJWTData(user user.User) jwt.JWTData {
	jwtData := jwt.JWTData{
		User: user.Name,
		Password: user.Password,
		Permissions: user.Permissions,
	}

	return jwtData
}
