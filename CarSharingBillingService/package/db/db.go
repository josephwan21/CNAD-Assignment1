// db.go
package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() *sql.DB {
	if db == nil {
		var err error
		// Adjust this connection string with your MySQL credentials
		dsn := "root:rootpassword@tcp(localhost:3307)/CarSharingBillingService"
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal("Error connecting to the database:", err)
		}
	}
	return db
}
