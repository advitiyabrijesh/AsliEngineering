package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time-trek/api/handlers"
	"time-trek/database"
)

func main() {
	// Connect to MySQL database
	db := database.InitDB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error in closing connection with MySQL.")
		}
	}(db)

	// Connect to Redis
	redisClient := database.InitRedis()
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			fmt.Println("Error in closing connection with Redis.")
		}
	}(redisClient)

	// Initialize Gin router
	router := gin.Default()

	// Set up routes
	router.POST("/api/register", func(c *gin.Context) {
		handlers.RegisterHandler(c, db)
	})

	router.POST("/api/login", func(c *gin.Context) {
		handlers.LoginHandler(c, db)
	})

	router.GET("/api/users/:user_id/location", func(c *gin.Context) {
		handlers.GetUserLocationHandler(c, db)
	})

	router.POST("/api/users/:user_id/location", func(c *gin.Context) {
		handlers.UpdateUserLocationHandler(c, db)
	})

	// Run the server
	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
