package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"mynsb-api/internal/filesint"
	"time"
)

// NOTE ========== THE jwt PACKAGE ONLY TAKES AND DEALS IN MAPS AS IT IS MEANT TO BE INDEPENDENT OF THE OTHER PACKAGES

// Attain the directory of the sensitive data

/**
	Func GetJWT:
		@param user map[string]interface{}

		returns string which is the jwt and returns an error if something wrong happened
 **/
func GenJWT(user map[string]interface{}) (string, error) {
	// Attain student data from the student parameter
	var username = user["student"]
	var password = user["Password"]
	var permissions = user["Permissions"]

	// Generate a token
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	// Start the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["student"] = username
	claims["Password"] = password
	claims["Permissions"] = permissions
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

		returns student data based on the current jwt
 **/
func ReadJWT(token string) (map[string]interface{}, error) {
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
		return make(map[string]interface{}), errors.New("invalid jwt")
	}

	// Get claims
	claims := tokenDec.Claims.(jwt.MapClaims)
	// Push to student

	perm := claims["Permissions"].([]interface{})

	// Convert to string
	for _, b := range perm {
		permissions = append(permissions, b.(string))
	}

	// Return that shit
	var user = make(map[string]interface{})
	user["student"] = claims["student"].(string)
	user["Password"] = claims["Password"].(string)
	user["Permissions"] = permissions

	return user, nil
}
