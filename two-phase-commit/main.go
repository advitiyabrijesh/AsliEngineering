package main

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	dsn := "root:@tcp(127.0.0.1:3306)/asli_engineering?charset=utf8mb4&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDatabase() *gorm.DB {
	return db
}

type Order struct {
	gorm.Model
	IsAvailable bool `gorm:"column:is_available" json:"is_available"`
}

type DeliveryAgent struct {
	gorm.Model
	IsAvailable bool   `gorm:"column:is_available" json:"is_available"`
	OrderId     string `gorm:"column:order_id" json:"order_id"`
}

func Init() {
	Connect()
	db = GetDatabase()
	err := db.AutoMigrate(&Order{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&DeliveryAgent{})
	if err != nil {
		panic(err)
	}
}

func ReserveAgent() (*DeliveryAgent, error) {
	txn := db.Begin()
	if txn.Error != nil {
		// Handle error. For example:
		log.Fatalf("Failed to start transaction: %v", txn.Error)
	}
	var deliveryAgent DeliveryAgent
	row := db.Where("is_reserved = ? AND order_id", false, nil).First(&deliveryAgent)
	if row.Error != nil || errors.Is(row.Error, gorm.ErrRecordNotFound) {
		txn.Rollback()
		return nil, errors.New("no delivery agent found")
	}

	result := txn.Model(&DeliveryAgent{}).Where("id = ?", deliveryAgent.ID).Update("is_reserved", true)
	if result.Error != nil {
		txn.Rollback()
		return nil, result.Error
	}

	commit := txn.Commit()
	if commit.Error != nil {
		return nil, commit.Error
	}
	return &deliveryAgent, nil

}

func BookAgent(orderId string) (*DeliveryAgent, error) {
	txn := db.Begin()
	if txn.Error != nil {
		// Handle error. For example:
		log.Fatalf("Failed to start transaction: %v", txn.Error)
	}
	var deliveryAgent DeliveryAgent
	row := db.Where("is_reserved = ? AND order_id", true, nil).First(&deliveryAgent)
	if row.Error != nil || errors.Is(row.Error, gorm.ErrRecordNotFound) {
		txn.Rollback()
		return nil, errors.New("no delivery agent found")
	}
	result := txn.Model(&DeliveryAgent{}).Where("id = ?", deliveryAgent.ID).Updates(map[string]interface{}{
		"is_reserved": false,
		"order_id":    orderId,
	})
	if result.Error != nil {
		txn.Rollback()
		return nil, result.Error
	}

	commit := txn.Commit()
	if commit.Error != nil {
		return nil, commit.Error
	}
	return &deliveryAgent, nil
}

func main() {
	Init()
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/delivery/agent/reserve", func(c *gin.Context) {

	})

	router.POST("/delivery/agent/book", func(c *gin.Context) {

	})

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
