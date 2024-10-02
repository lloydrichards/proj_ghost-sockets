package database

import (
	"database/sql"
	"log"
)

func migrate(db *sql.DB) error {
	log.Println("Running migrations...")
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS clients (
    		id VARCHAR PRIMARY KEY,
    		username VARCHAR NOT NULL,
    		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    		last_seen TIMESTAMP
		)
	`)
	if err != nil {
		log.Println("Error running migrations:", err)
		return err
	}
	log.Println("Migrations ran successfully")
	return nil
}
