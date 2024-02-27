// api/handlers/nearby_locations.go

package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Location represents a geographical location
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Type      string  `json:"type"`
	Rating    float64 `json:"rating"`
}

// GetNearbyLocationsHandler retrieves top places to explore nearby within a specified radius.
func GetNearbyLocationsHandler(c *gin.Context, db *sql.DB) {
	// Parse query parameters
	userIDParam := c.Query("user_id")
	radiusParam := c.Query("radius")

	// Validate and parse user ID parameter
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Validate and parse radius parameter
	radius, err := strconv.ParseFloat(radiusParam, 64)
	if err != nil || radius <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-positive radius"})
		return
	}

	// Fetch user's current location from the database
	userLocation, err := fetchUserLocation(userID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user location"})
		return
	}

	// Fetch nearby locations from the database within the specified radius
	nearbyLocations, err := fetchNearbyLocations(userLocation, radius, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch nearby locations"})
		return
	}

	// Return the nearby locations in the response
	c.JSON(http.StatusOK, gin.H{"nearby_locations": nearbyLocations})
}

// fetchUserLocation fetches the user's current location from the database
func fetchUserLocation(userID int, db *sql.DB) (*Location, error) {
	var userLocation Location
	err := db.QueryRow(
		"SELECT latitude, longitude FROM user_location_data WHERE user_id = ? ORDER BY timestamp DESC LIMIT 1",
		userID,
	).Scan(&userLocation.Latitude, &userLocation.Longitude)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user location not found")
		}
		return nil, err
	}

	return &userLocation, nil
}

// fetchNearbyLocations fetches nearby locations from the database within the specified radius
func fetchNearbyLocations(userLocation *Location, radius float64, db *sql.DB) ([]Location, error) {
	// Use the Haversine formula to calculate distances
	query := "SELECT latitude, longitude, type, rating FROM top_locations " +
		"WHERE ST_Distance_Sphere(point(latitude, longitude), point(?, ?)) <= ?"
	rows, err := db.Query(query, userLocation.Latitude, userLocation.Longitude, radius)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Error in closing rows object.")
		}
	}(rows)

	var nearbyLocations []Location

	for rows.Next() {
		var location Location
		err := rows.Scan(&location.Latitude, &location.Longitude, &location.Type, &location.Rating)
		if err != nil {
			return nil, err
		}
		nearbyLocations = append(nearbyLocations, location)
	}

	return nearbyLocations, nil
}
