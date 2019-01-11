package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"mynsb-api/internal/filesint"
	"time"
)

// NOTE ========== THE jwt PACKAGE ONLY TAKES THE JWTData struct

type JWTData struct {
	User 	 	string
	Password 	string
	Permissions []string
}


// Attain the directory of the sensitive data

/**
	Func GetJWT:
		@param user map[string]interface{}

		returns string which is the jwt and returns an error if something wrong happened
 **/
func GenJWT(user JWTData) (string, error) {
	// Generate a token
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	// Start the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["User"] = user.User
	claims["Password"] = user.Password
	claims["Permissions"] = user.Permissions
	claims["Expires"] = time.Now().Add(time.Hour * 24 * 30)

	// Read key
	privKey, err := filesint.DataDump("sensitive", "/keys/priv.txt")
	if err != nil {
		return "", errors.New("error generating jwt")
	}

	// Generate the signed token
	signedToken, err := token.SignedString(privKey)
	if err != nil {
		return "", errors.New("error generating jwt")
	}

	// Throw back the signed token
	return signedToken, nil
}

/**
	Func ReadJWT:
		@param jwt string

		returns user data based on the current jwt
 **/
func ReadJWT(token string) (JWTData, error) {
	// Decode the token

	var permissions []string

	// Decode token
	tokenDec, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Read the token from the path provided
		key, err := filesint.DataDump("sensitive", "/keys/priv.txt")
		if err != nil {
			return nil, errors.New("error parsing jwt")
		}

		return []byte(key), nil

	})

	// Check for err
	if err != nil || !tokenDec.Valid {
		return JWTData{}, errors.New("invalid jwt")
	}

	// Get claims
	claims := tokenDec.Claims.(jwt.MapClaims)
	// Push to user

	perm := claims["Permissions"].([]interface{})

	// Convert to string
	for _, b := range perm {
		permissions = append(permissions, b.(string))
	}

	// Return that shit
	user := JWTData{
		User: claims["User"].(string),
		Password: claims["Password"].(string),
		Permissions: permissions,
	}

	return user, nil
}
