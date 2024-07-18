package models

import (
	"database/sql"
	"log"

	"eoncohub.com/person_module/db"
)

func GetJudById(id int64) (*Jud, error) {
	var jud Jud
	query := `
		SELECT jud."ID_JUD", jud."NAME"
		FROM "JUD" jud 
		JOIN "LOC" loc ON jud."ID_JUD" = loc."ID_JUD"
		JOIN "ADDRESS" address ON loc."ID_LOC" = address."ID_LOC"
		WHERE address."ID_ADDRESS" = :1`

	// Prepare the statement
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		log.Fatalf("Error preparing statement: %v", err)
		return nil, err
	}
	defer stmt.Close()

	// Execute the query with the bind argument
	err = stmt.QueryRow(id).Scan(&jud.ID, &jud.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle case where no rows were returned
			return nil, nil
		}
		log.Fatalf("Error executing query: %v", err)
		return nil, err
	}

	return &jud, nil
}
