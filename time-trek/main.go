package main

import (
	"github.com/gin-gonic/gin"
	"time-trek/api/handlers"
	"time-trek/database"
)

func main() {
	// Connect to MySQL database
	db := database.InitDB()
	defer db.Close()

	// Connect to Redis
	redisClient := database.InitRedis()
	defer redisClient.Close()

	// Initialize Gin router
	router := gin.Default()

	// Set up routes
	router.POST("/api/register", func(c *gin.Context) {
		handlers.RegisterHandler(c, db)
	})

	// Run the server
	router.Run(":8080")
}
