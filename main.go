package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	// Load database credentials from environment variables
	dbUser, exist := os.LookupEnv("DB_USER")
	if !exist {
		log.Fatalf("Failed to retrieve DB_USER variable")
	}
	dbPassword, exist := os.LookupEnv("DB_PASSWORD")
	if !exist {
		log.Fatalf("Failed to retrieve DB_PASSWORD variable")
	}
	dbServer, exist := os.LookupEnv("DB_SERVER")
	if !exist {
		log.Fatalf("Failed to retrieve DB_SERVER variable")
	}
	dbName, exist := os.LookupEnv("DB_NAME")
	if !exist {
		log.Fatalf("Failed to retrieve DB_NAME variable")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	version := "1.0.1"

	// Construct connection string
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(dbUser, dbPassword),
		Host:     dbServer,
		RawQuery: "database=" + dbName,
	}
	// Connect to the database
	db, err := sql.Open("sqlserver", u.String())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create table if it doesn't exist
	_, err = db.Exec("IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'Requests') BEGIN CREATE TABLE Requests (ID INT IDENTITY(1,1) PRIMARY KEY, CreateDate DATETIME) END")
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Handle HTTP requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Insert the current datetime into the database
		_, err := db.Exec("INSERT INTO Requests (CreateDate) VALUES (GETDATE())")
		if err != nil {
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}

		// Retrieve the number of rows
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM Requests").Scan(&count)
		if err != nil {
			http.Error(w, "Failed to count rows", http.StatusInternalServerError)
			return
		}

		// Measure database response time
		dbResponseTime := time.Since(start)

		// Get hostname
		hostname, err := os.Hostname()
		if err != nil {
			http.Error(w, "Failed to get hostname", http.StatusInternalServerError)
			return
		}

		// Respond with details
		response := fmt.Sprintf("Hello, world!\nVersion: %s\nNumber of requests: %d\nHostname: %s\nDB Response Time: %v\n", version, count, hostname, dbResponseTime)
		log.Printf("Success\t%s\n", r.RemoteAddr)
		fmt.Fprint(w, response)
	})

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
