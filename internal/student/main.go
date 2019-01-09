package student

type User struct {
	Name        string
	Password    string
	Permissions []string
}

func IsAdmin(user User) bool {
	Permissions := user.Permissions

	for _, b := range Permissions {
		if b == "admin" {
			return true
		}
	}
	return false
}
