package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"mynsb-api/internal/filesint"
	"time"
)

// NOTE ========== THE jwt PACKAGE ONLY TAKES THE JWTData struct


// JWTData struct for dealing with JWTData
type JWTData struct {
	User 	 	string
	Password 	string
	Permissions []string
}




// UTILITY FUNCTIONS

// GenJWT generates a JWT token based off the userData provided
func GenJWT(user JWTData) (string, error) {

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	oneMonth := time.Hour * 24 * 30

	claims := token.Claims.(jwt.MapClaims)
	claims["User"] = user.User
	claims["Password"] = user.Password
	claims["Permissions"] = user.Permissions
	claims["Expires"] = time.Now().Add(oneMonth)

	// Get the private key
	privateKey, err := filesint.DataDump("sensitive", "/keys/priv.txt")
	if err != nil {
		return "", errors.New("error generating jwt")
	}

	// Generate the signed token
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.New("error generating jwt")
	}

	return signedToken, nil
}


// ReadJWT takes a JWT token and decodes it into a JWTData object
func ReadJWT(token string) (JWTData, error) {

	var permissions []string

	// Decode token
	tokenData, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Get the private key to decode the token
		privateKey, err := filesint.DataDump("sensitive", "/keys/priv.txt")
		if err != nil {
			return nil, errors.New("error parsing jwt")
		}

		return []byte(privateKey), nil
	})
	if err != nil || !tokenData.Valid {
		return JWTData{}, errors.New("invalid jwt")
	}

	// Get claims
	claims := tokenData.Claims.(jwt.MapClaims)

	// Convert the permissions to a string array
	perms := claims["Permissions"].([]interface{})
	for _, b := range perms {
		permissions = append(permissions, b.(string))
	}


	return JWTData{
		User: claims["User"].(string),
		Password: claims["Password"].(string),
		Permissions: permissions,
	}, nil
}
