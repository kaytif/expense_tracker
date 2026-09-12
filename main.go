package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Connect to PostgreSQL.
	connectDB()

	// Send API requests to our expense handler.
	http.HandleFunc("/expenses", expenseHandler)

	// Serve the frontend when the user visits "/".
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Render provides PORT when deployed; use 8080 locally.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server running on port " + port)

	// Start the HTTP server.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("server failed: ", err)
	}
}