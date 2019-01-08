package User

import (
	"io/ioutil"
	"time"
	"errors"
	"github.com/dgrijalva/jwt-go"
)

/**
	Func GetJWT:
		@param user User
		@param keypath string

		returns string which is the JWT and returns an error if something wrong happened
 **/
func GenJWT(user User, privkeyPath string) (string, error) {
	// Attain user data from the user paramater
	var username = user.Name
	var password = user.Password
	var permissions = user.Permissions


	// Generate a token
	token := jwt.New(jwt.GetSigningMethod("HS256"))

	// Start the claims
	claims := token.Claims.(jwt.MapClaims)
	claims["User"] = username
	claims["Password"] = password
	claims["Permissions"] = permissions
	claims["Expires"] = time.Now().Add(time.Hour * 24 * 30)

	// Read key
	privKey, err := ioutil.ReadFile(privkeyPath)
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
		@param JWT string

		returns User data based on the current jwt
 **/
func ReadJWT(token, keypath string) (User, error) {
	// Decode the token

	// Generate temp details
	var username string
	var password string
	var permissions []string

	// Decode token
	tokenDec, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Read the token from the path provided
		key, err := ioutil.ReadFile(keypath)
		if err != nil {
			return nil, errors.New("error parsing jwt")
		}

		return []byte(key), nil

	})

	// Check for err
	if err != nil || !tokenDec.Valid {
		return User{}, errors.New("invalid jwt")
	}


	// Get claims
	claims := tokenDec.Claims.(jwt.MapClaims)
	// Push to user
	username = claims["User"].(string)
	password = claims["Password"].(string)
	perm := claims["Permissions"].([]interface{})

	// Convert to string
	for _, b := range perm {
		permissions = append(permissions, b.(string))
	}

	// Return that shit
	var user = User{
		Name:        username,
		Password:    password,
		Permissions: permissions,
	}

	return user, nil
}