package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
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

	// Store the user in the database
	err := models.InsertUser(db, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
