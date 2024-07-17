package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/godror/godror"
	"github.com/joho/godotenv"
)

var db *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbParams := map[string]string{
		"service":        os.Getenv("DB_SERVICE"),
		"username":       os.Getenv("DB_USERNAME"),
		"server":         os.Getenv("DB_SERVER"),
		"port":           os.Getenv("DB_PORT"),
		"password":       os.Getenv("DB_PASSWORD"),
		"walletLocation": os.Getenv("WALLET_LOCATION"),
	}

	// Construct the connection string with wallet location
	db, err = sql.Open("godror", fmt.Sprintf(`user="%s" password="%s"
		connectString="tcps://%s:%s/%s?wallet_location=%s"
		   `, dbParams["username"], dbParams["password"], dbParams["server"], dbParams["port"], dbParams["service"], dbParams["walletLocation"]))

	if err != nil {
		panic(fmt.Errorf("error in sql.Open: %w", err))
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Println("Can't close connection: ", err)
		}
	}()

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}

	fmt.Println("Successfully connected to the database!")

}

func CloseDB() {

	err := db.Close()
	if err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
	fmt.Println("Database connection closed.")
}

func someAdditionalActions() {
	// Example of additional database actions
	rows, err := db.Query("SELECT 2+3 FROM dual")
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		fmt.Println(rows.Columns())
	}

	err = rows.Err()
	if err != nil {
		log.Fatalf("Error processing rows: %v", err)
	}
}
