package service


type UserService interface {
	CreateUser(name string, email string, passwrd string, role string) error
	Login(email string, passwrd string) error 
}