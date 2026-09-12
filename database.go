package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"

	// PostgreSQL driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

// db will hold our connection to PostgreSQL
var db *sql.DB

func connectDB() {

	var err error
	// Load variables from the local .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file: ", err)
	}

	// Read database credentials from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Build the postgreSQL connection string.
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user, password, host, port, dbname)
	

	// Open the database connection
	db, err = sql.Open("pgx", connectionString)

	if err != nil {
		log.Fatal("failed to open database: ", err)
	}

	// Check that PostreSQ; is actually reachable
	if err = db.Ping(); err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
}
