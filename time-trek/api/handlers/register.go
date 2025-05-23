package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time-trek/auth"
	"time-trek/models"
)

func RegisterHandler(c *gin.Context, db *sql.DB) {
	var newUser models.User

	// Bind JSON request body to the User struct
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the username already exists
	if models.UserExists(db, newUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Hash the password before storing
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser.Password = hashedPassword

	// Store the user in the database
	err = models.InsertUser(db, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
