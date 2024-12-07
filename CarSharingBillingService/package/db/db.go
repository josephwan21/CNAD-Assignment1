// db.go
package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

func InitDB() error {
	var err error
	dbConn, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/carsharingbillingservice?parseTime=true")
	if err != nil {
		return err
	}

	// Verify the connection
	if err := dbConn.Ping(); err != nil {
		return err
	}
	return nil
}

// GetDBConn retrieves the database connection
func GetDBConn() *sql.DB {
	return dbConn
}

// CloseDB closes the database connection pool (if needed)
func CloseDB() {
	if dbConn != nil {
		err := dbConn.Close()
		if err != nil {
			log.Println("Error closing the database connection:", err)
		}
	}
}
