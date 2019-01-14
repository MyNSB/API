package auth

import "mynsb-api/internal/user"
import "mynsb-api/internal/jwt"


// userToJWTData takes a user and converts their data into a format suitable for the JWT package
func userToJWTData(user user.User) jwt.JWTData {
	return jwt.JWTData{
		User: user.Name,
		Password: user.Password,
		Permissions: user.Permissions,
	}
}
