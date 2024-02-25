package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// InitDB initializes the MySQL database connection
func InitDB() *sql.DB {
	connectionString := "root:@tcp(localhost:3306)/time-trek"
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil
	}

	// Ping the database to check the connection
	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging database:", err)
		return nil
	}

	fmt.Println("Connected to the database")

	db = conn
	return db
}
