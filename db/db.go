package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/sijms/go-ora/v2"
)

var DB *sql.DB

func InitDB() {
	dbParams := map[string]string{
		"service":        os.Getenv("DB_SERVICE"),
		"username":       os.Getenv("DB_USERNAME"),
		"server":         os.Getenv("DB_SERVER"),
		"port":           os.Getenv("DB_PORT"),
		"password":       os.Getenv("DB_PASSWORD"),
		"walletLocation": os.Getenv("WALLET_LOCATION"),
	}

	connectionString := fmt.Sprintf(
		"oracle://%s:%s@%s:%s/%s?SSL=enable&SSL Verify=false&WALLET=%s",
		dbParams["username"], dbParams["password"], dbParams["server"],
		dbParams["port"], dbParams["service"], dbParams["walletLocation"],
	)

	var err error
	DB, err = sql.Open("oracle", connectionString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	// Set connection pool parameters
	DB.SetMaxOpenConns(25)                 // Set the maximum number of open connections
	DB.SetMaxIdleConns(25)                 // Set the maximum number of idle connections
	DB.SetConnMaxLifetime(5 * time.Minute) // Set the maximum amount of time a connection may be reused
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
		fmt.Println("Database connection closed.")
	}
}
