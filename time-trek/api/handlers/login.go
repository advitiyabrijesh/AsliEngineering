package handlers

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time-trek/auth"
	"time-trek/models"
)

func LoginHandler(c *gin.Context, db *sql.DB) {
	var loginRequest models.User

	// Bind JSON request body to the LoginRequest struct
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate login credentials
	user, err := models.GetUserByEmail(db, loginRequest.Email)
	fmt.Println(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	hashedPassword, _ := auth.HashPassword(loginRequest.Password)
	fmt.Println(hashedPassword)
	// Check if the provided password matches the stored hash
	if !auth.CheckPasswordHash(hashedPassword, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := auth.GenerateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
