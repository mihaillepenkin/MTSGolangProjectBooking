package entity


type User struct {
	ID        string
	Name      string
	Role      string
	Email     string
	Passwrd   string
}

func(u *User) Validate() bool {
	switch u.Role {
	case "hotelier":
		return true
	case "client":
		return true
	default:
		return false		
	}
}