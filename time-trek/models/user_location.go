package models

import (
	"database/sql"
)

// UserLocation struct to represent user location data
type UserLocation struct {
	UserLocationID int     `json:"user_location_id"`
	UserID         int     `json:"user_id"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Timestamp      string  `json:"timestamp"`
}

// GetUserLocation retrieves user's current location
func GetUserLocation(db *sql.DB, userID int) (UserLocation, error) {
	var location UserLocation
	query := "SELECT * FROM user_location_data WHERE user_id = ? ORDER BY timestamp DESC LIMIT 1"
	err := db.QueryRow(query, userID).Scan(
		&location.UserLocationID,
		&location.UserID,
		&location.Latitude,
		&location.Longitude,
		&location.Timestamp,
	)
	return location, err
}

// InsertUserLocation stores user's location in the database
func InsertUserLocation(db *sql.DB, location UserLocation) error {
	query := "INSERT INTO user_location_data (user_id, latitude, longitude) VALUES (?, ?, ?)"
	_, err := db.Exec(query, location.UserID, location.Latitude, location.Longitude)
	return err
}
