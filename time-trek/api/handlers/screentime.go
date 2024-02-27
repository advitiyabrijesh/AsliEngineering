package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetScreenTimeHandler retrieves screen time statistics for a specified duration.
func GetScreenTimeHandler(c *gin.Context, db *sql.DB) {
	// Parse query parameters
	userIDParam := c.Query("user_id")
	durationParam := c.Query("duration")

	// Validate and parse user ID parameter
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Validate and parse duration parameter
	duration, err := strconv.Atoi(durationParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration parameter. Must be a positive integer"})
		return
	}

	// Calculate start date based on the duration
	startDate := time.Now().Add(-time.Duration(duration) * time.Minute)

	// Fetch screen time data from the database for the specified duration
	screenTimeData, err := fetchScreenTimeData(userID, startDate, db)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch screen time data"})
		return
	}

	// Return the screen time statistics in the response
	c.JSON(http.StatusOK, gin.H{"screen_time_data": screenTimeData})
}

// fetchScreenTimeData fetches screen time data from the database for the specified duration
func fetchScreenTimeData(userID int, startDate time.Time, db *sql.DB) ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT timestamp FROM user_location_data WHERE user_id = ? AND timestamp >= ?", userID, startDate)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("Failed to close MySQL rows object.")
		}
	}(rows)

	screenTimeData := make([]map[string]interface{}, 0)

	duration := 0.0
	lastTime := time.Now()
	currentTime := lastTime
	for rows.Next() {
		var timestamp []uint8
		if err := rows.Scan(&timestamp); err != nil {
			return nil, err
		}
		// Parse string to time.Time
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(timestamp))

		if err != nil {
			fmt.Println("Error parsing time:", err)
			return nil, err
		}
		if lastTime == currentTime {
			lastTime = parsedTime
			continue
		}
		timeDiff := parsedTime.Sub(lastTime).Seconds()
		// Calculate the screen time duration for each record
		if timeDiff <= 60 {
			duration += timeDiff
		}
		lastTime = parsedTime
	}
	screenTimeData = append(screenTimeData, map[string]interface{}{
		"duration_in_seconds": duration,
	})

	return screenTimeData, nil
}
