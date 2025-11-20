package usecase

import (
	"auth_service/internal/application/dto"
	"auth_service/internal/application/interface/repository"
	"auth_service/internal/application/interface/service"
	"auth_service/internal/domain/entity"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo   repository.UserRepository
	jwt        service.JwtService
	jwtRepo    repository.JwtRepository
}

func UserService(userRepo repository.UserRepository, jwt service.JwtService, jwtRepo repository.JwtRepository) *UserUseCase {
	serv := new(UserUseCase)
	serv.userRepo = userRepo
	serv.jwt = jwt
	serv.jwtRepo = jwtRepo
	return serv
}

func(u *UserUseCase) CreateUser(name string, email string, passwrd string, role string) dto.CreateUserOutputDto {
	userExist := u.userRepo.CheckUserByEmail(email)
	if (userExist) {
		return dto.CreateUserOutputDto{Msg : fmt.Errorf("Error: user with this email already exist").Error(), Status : "error"}
	}
	user := new(entity.User)
	user.Name = name
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwrd), bcrypt.DefaultCost)
	if err != nil {
		return dto.CreateUserOutputDto{Msg : err.Error(), Status : "error"}
	}
	user.Passwrd = string(hashed)
	user.Role = role
	user.Email = email
	roleIsNormal := user.Validate()
	if (!roleIsNormal) {
		return dto.CreateUserOutputDto{Msg : fmt.Errorf("Invalid role name").Error(), Status : "error"}
	}
	err = u.userRepo.CreateUser(user)
	
	if (err != nil) {
		return dto.CreateUserOutputDto{Msg : err.Error(), Status : "error"}
	} else {
		return dto.CreateUserOutputDto{Msg : "", Status : "ok"}
	}
}

func(u *UserUseCase) Login(email string, passwrd string) *dto.LoginOutputDto {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return &dto.LoginOutputDto{
			Token:      "",
			Msg:        fmt.Sprintf("Error in UserService: %w", err),
			Status:     "error",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Passwrd), []byte(passwrd)); err != nil {
		return &dto.LoginOutputDto{
			Token:      "",
			Msg:        fmt.Sprintf("Error in UserService: %w", err),
			Status:     "error",
		}
	}

	access, _, err := u.jwt.GenerateAccessToken(user)
	if err != nil {
		return &dto.LoginOutputDto{
			Token:      "",
			Msg:        fmt.Sprintf("Error in UserService: %w", err),
			Status:     "error",
		}
	}

	refresh, rexp, err := u.jwt.GenerateRefreshToken(user)
	if err != nil {
		return &dto.LoginOutputDto{
			Token:      "",
			Msg:        fmt.Sprintf("Error in UserService: %w", err),
			Status:     "error",
		}
	}

	if err := u.jwtRepo.SaveRefreshToken(user.ID, refresh, rexp); err != nil {
		return &dto.LoginOutputDto{
			Token:      "",
			Msg:        fmt.Sprintf("Error in UserService: %w", err),
			Status:     "error",
		} 
	}

	return &dto.LoginOutputDto{
		Token:      access,
		Msg:        "",
		Status:     "ok",
	}
}