package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

type GeoPoint struct {
	Latitude  float64
	Longitude float64
}

func main() {
	// Create a Redis client
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Your Redis server address
		Password: "",               // No password
		DB:       0,                // Default DB
	})

	locations := []GeoPoint{
		{17.479631, 78.360558}, // Hafeezpet, Hyderabad
		{17.482456, 78.357940}, // Hafeezpet, Hyderabad
		{27.580289, 77.701320}, // GandhiPark, Vrindavan
		{28.509658, 77.406836}, // Paramount, Noida
		{17.458106, 78.372505}, // Google office, Hyderabad
		{17.239925, 78.430061}, // Hyderabad Airport
	}

	// Key to store geospatial data in Redis
	redisKey := "locations"
	for _, location := range locations {
		latitude := location.Latitude
		longitude := location.Longitude
		fmt.Println(location)
		// Add the geospatial location to Redis
		addLocation(client, redisKey, latitude, longitude)
		// Get the geohash for the location
		getGeohash(client, redisKey, latitude, longitude)
	}

}

// Function to add a geospatial location to Redis
func addLocation(client *redis.Client, key string, latitude, longitude float64) {
	err := client.GeoAdd(key, &redis.GeoLocation{
		Name:      "location",
		Latitude:  latitude,
		Longitude: longitude,
	}).Err()

	if err != nil {
		fmt.Println("Error adding location:", err)
	} else {
		fmt.Println("Location added")
	}
}

// Function to get the geohash for a given location from Redis
func getGeohash(client *redis.Client, key string, latitude, longitude float64) {
	result, err := client.GeoHash(key, "location").Result()

	if err != nil {
		fmt.Println("Error getting geohash:", err)
	} else {
		fmt.Println("Geohash:", result)
	}
}
