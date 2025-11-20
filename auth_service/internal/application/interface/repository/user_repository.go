package repository

import "auth_service/internal/domain/entity"


type UserRepository interface {
	GetUserById(ID int) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	CheckUserByEmail(email string) bool
	CreateUser(user *entity.User) error
}