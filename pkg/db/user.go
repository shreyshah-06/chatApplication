package db

import (
	"database/sql"
	"errors"
	"fmt"
	"gochatapp/utils"
)

// User model
type User struct {
	Username string
	Password string
}

// RegisterNewUser registers a new user in the database
func RegisterNewUser(db *sql.DB, username, password string) error {
	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err := db.Exec(query, username, password)
	if err != nil {
		return fmt.Errorf("error registering user: %w", err)
	}
	return nil
}

// IsUserExist checks if a user exists in the database
func IsUserExist(db *sql.DB, username string) (bool) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)"
	err := db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// IsUserAuthentic checks if the username and password match for authentication
func IsUserAuthentic(db *sql.DB, username, password string) error {
	var storedPassword string
	query := "SELECT password FROM users WHERE username = $1"
	err := db.QueryRow(query, username).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		// If the user doesn't exist, return an error
		return errors.New("invalid username or password")
	}
	if err != nil {
		// Handle other potential database errors
		return err
	}

	// Compare the hashed password stored in the database with the provided password
	if !utils.CheckPasswordHash(password, storedPassword) {
		// If the password doesn't match, return an error
		return errors.New("invalid username or password")
	}

	// If everything is valid, return nil (successful authentication)
	return nil
}
