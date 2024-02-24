package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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
	mysqlDB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/asli_engineering")
	if err != nil {
		log.Fatal(err)
	}

	err = mysqlDB.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/storeLocation", storeLocationHandler).Methods("POST")

	http.Handle("/", r)

	fmt.Println("Server is running on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func storeLocationHandler(w http.ResponseWriter, r *http.Request) {
	var locationData LocationData

	// Decode JSON payload
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&locationData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Store data in Redis
	err = storeInRedis(locationData)
	if err != nil {
		http.Error(w, "Error storing data in Redis", http.StatusInternalServerError)
		return
	}

	// Store data in MySQL
	err = storeInMySQL(locationData)
	if err != nil {
		http.Error(w, "Error storing data in MySQL", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data stored successfully"))
}

func storeInRedis(data LocationData) error {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Store data in Redis
	err = redisClient.Set(redisClient.Context(), fmt.Sprintf("locationData:%d", data.UserID), jsonData, 0).Err()
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
