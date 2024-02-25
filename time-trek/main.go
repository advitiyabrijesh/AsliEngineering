package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"time"
	_ "time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	redisClient *redis.Client
	mysqlDB     *sql.DB
)

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	TimeEpoch int64   `json:"timeEpoch"`
	UserID    int64   `json:"userId"`
}

func init() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password
		DB:       0,  // Default DB
	})

	// Initialize MySQL database connection
	var err error
	mysqlDB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/time-trek")
	if err != nil {
		log.Fatal(err)
	}

	err = mysqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Set Gin to production mode for better performance
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.POST("/storeLocation", storeLocationHandler)

	fmt.Println("Server is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func storeLocationHandler(c *gin.Context) {
	var locationData LocationData

	// Decode JSON payload
	if err := c.BindJSON(&locationData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Store data in Redis
	if err := storeInRedis(locationData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing data in Redis"})
		return
	}

	// Store data in MySQL
	if err := storeInMySQL(locationData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing data in MySQL"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Data stored successfully"})
}

func storeInRedis(data LocationData) error {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Store data in Redis
	err = redisClient.Set(
		fmt.Sprintf("locationData:%d", data.UserID),
		jsonData,
		time.Duration(1)*time.Minute,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func storeInMySQL(data LocationData) error {
	// Check if the table exists
	if _, err := mysqlDB.Exec("CREATE TABLE IF NOT EXISTS location_data (id INT AUTO_INCREMENT PRIMARY KEY, latitude DOUBLE, longitude DOUBLE, time_epoch INT, user_id INT)"); err != nil {
		return err
	}

	// Prepare SQL statement
	stmt, err := mysqlDB.Prepare("INSERT INTO location_data(latitude, longitude, time_epoch, user_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute SQL statement
	_, err = stmt.Exec(data.Latitude, data.Longitude, data.TimeEpoch, data.UserID)
	if err != nil {
		return err
	}

	return nil
}
