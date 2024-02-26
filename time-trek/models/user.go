package models

import (
	"database/sql"
	"fmt"
)

// User struct to represent user data
type User struct {
	UserId    int     `json:"user_id"`
	CreatedAt []uint8 `json:"created_at"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
}

// UserExists checks if a user with the given username already exists
func UserExists(db *sql.DB, email string) bool {

	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		fmt.Println("Error checking user existence:", err)
		return false
	}

	return count > 0
}

// InsertUser inserts a new user into the database
func InsertUser(db *sql.DB, user User) error {
	query := "INSERT INTO users (email, password, first_name, last_name) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, user.Email, user.Password, user.FirstName, user.LastName)
	return err
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(db *sql.DB, email string) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(
		&user.UserId,
		&user.CreatedAt,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
	)
	return user, err
}
