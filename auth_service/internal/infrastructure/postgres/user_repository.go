package postgres

import (
	"auth_service/internal/domain/entity"
	"database/sql"
	"fmt"
	"log"
)

type UserRepoPG struct {
	DB *sql.DB
}

func (rep *UserRepoPG) GetUserById(ID int) (*entity.User, error) {
	tx, err := rep.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	user := &entity.User{}
	row := tx.QueryRow(`SELECT id, name, role, email, password FROM users WHERE id = $1`, ID)
	if err := row.Scan(&user.ID, &user.Name, &user.Role, &user.Email, &user.Passwrd); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return user, nil
}

func (rep *UserRepoPG) GetUserByEmail(email string) (*entity.User, error) {
	tx, err := rep.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	user := &entity.User{}
	row := tx.QueryRow(`SELECT id, name, role, email, password FROM users WHERE email = $1`, email)
	if err := row.Scan(&user.ID, &user.Name, &user.Role, &user.Email, &user.Passwrd); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return user, nil
}
func (rep *UserRepoPG) CheckUserByEmail(email string) bool {
	_, err := rep.GetUserByEmail(email)
	log.Printf("error checking user by email %s: %v", email, err)
	return err == nil
}

func (rep *UserRepoPG) CreateUser(user *entity.User) error {
	tx, err := rep.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO users (name, role, email, password) VALUES ($1, $2, $3, $4)`,
		user.Name, user.Role, user.Email, user.Passwrd)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil

}
