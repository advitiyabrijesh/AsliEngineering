package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time-trek/models"
)

// GetUserLocationHandler retrieves user's current location
func GetUserLocationHandler(c *gin.Context, db *sql.DB) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	location, err := models.GetUserLocation(db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user location"})
		return
	}

	c.JSON(http.StatusOK, location)
}

// UpdateUserLocationHandler shares user's current location
func UpdateUserLocationHandler(c *gin.Context, db *sql.DB) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	var location models.UserLocation

	// Bind JSON request body to the UserLocation struct
	if err := c.BindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	location.UserID = userID

	// Store user's location in the database
	err = models.InsertUserLocation(db, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share user location"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User location shared successfully"})
}
