package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	row := us.DB.QueryRow(`
	  INSERT INTO users (email, password_hash)
	  VALUES ($1, $2) RETURNING id;`, email, passwordHash)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	// Normalize data
	email = strings.ToLower(email)

	// Create the user instance
	user := User{
		Email: email,
	}

	// Query to the DB for the id and password_hash of the user
	row := us.DB.QueryRow(`
	SELECT id, password_hash
	FROM users
	WHERE email = $1;`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	// Authenticate the password_hash with provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return &user, nil
}

func (us *UserService) UpdateUser(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = us.DB.Exec(`
	  UPDATE users
	  SET password_hash=$2
	  WHERE id=$1`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}
