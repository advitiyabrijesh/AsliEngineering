package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func WithoutConnectionPool() {
	startTime := time.Now()

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/asli_engineering")
	if err != nil {
		fmt.Println(err)
	}

	db.SetMaxOpenConns(100)

	execute(db, startTime)
}

func execute(db *sql.DB, startTime time.Time) {
	rows, err := db.Query("SELECT * FROM KV_STORE")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var key_name, value_name string
		var expiry_at string
		err = rows.Scan(&key_name, &value_name, &expiry_at)
		if err != nil {
			fmt.Println(err)
		}
		// Parse the string into time.Time
		expiryTime, err := time.Parse("2006-01-02 15:04:05", expiry_at) // Adjust the layout string as needed
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}
		fmt.Printf("Key: %s, Value: %s, Expiry: %s\n", key_name, value_name, expiryTime)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}

	endTime := time.Now()
	execitonTime := endTime.Sub(startTime)

	fmt.Printf("Execution time: %v\n", execitonTime)
}

type ConnectionPool struct {
	pool           *sql.DB
	maxConnections int
	numConnections int
	mutex          *sync.Mutex
}

func NewConnectionPool(dsn string, maxConnections int) (*ConnectionPool, error) {
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	pool.SetMaxOpenConns(maxConnections)

	p := &ConnectionPool{
		pool:           pool,
		maxConnections: maxConnections,
		numConnections: 0,
		mutex:          &sync.Mutex{},
	}
	return p, nil
}

func (p *ConnectionPool) GetConnection() (*sql.DB, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.numConnections == p.maxConnections {
		return nil, fmt.Errorf("connection pool is full")
	}

	p.numConnections++

	return p.pool, nil
}

func (p *ConnectionPool) ReleaseConnection(conn *sql.DB) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.numConnections--
}

func WithConnectionPool() {

	startTime := time.Now()
	pool, err := NewConnectionPool("root:@tcp(localhost:3306)/asli_engineering", 100)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.pool.Close()

	conn, err := pool.GetConnection()
	if err != nil {
		log.Fatal(err)
	}
	execute(conn, startTime)
}

func main() {
	WithConnectionPool()
}
