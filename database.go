package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func connectDB() {
	// Load .env locally. On Redner, environment variables are provided directly.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}

	// Read the complete PostgreSQL connection URL.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Connect to PostgreSQL.
	var err error
	db, err = sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal("failed to open database: ", err)
	}

	// Verify that PostgreSQL is reachable.
	if err = db.Ping(); err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
}