package Tests

import (
	"User"
	"testing"
	"User/JWT"
)

func TestJWT(t *testing.T) {
	// Create a test user
	user := User.User{
		Name:        "TestUser",
		Password:    "TestPWD",
		Permissions: []string{"Tester"},
	}

	// Create a jwt from this user
	jwt, err := JWT.GenJWT(user, "key.txt")
	// Check for error
	if err != nil {
		t.Errorf(err.Error())
	}

	// Decode this jwt to check that it was decoded correctly
	userDecode, err := JWT.ReadJWT(jwt, "key.txt")
	if err != nil {
		t.Errorf(err.Error())
	}

	// Check that it was decoded correctly
	if !(user.Name == userDecode.Name && user.Password == userDecode.Password && user.Permissions[0] == userDecode.Permissions[0]) {
		// Throw error
		t.Error("Token was not decoced/encoded properly")
	}

}
