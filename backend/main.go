package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db          *sql.DB
	port        string
	frontendURL string
)

func runInitSQL(db *sql.DB, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Split statements on semicolon
	//queries := strings.SplitSeq(string(sqlBytes), ";")
	queries := strings.Split(string(sqlBytes), ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error executing query: %q\n%v", query, err)
		}
	}
	return nil
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database/app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = runInitSQL(db, "./database/init.sql")
	if err != nil {
		log.Fatalf("Failed to run init.sql: %v", err)
	}

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	frontendURL = os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173" // Default frontend URL
	}

	// CORS to allow developement on same address
	http.HandleFunc("/api/users", withCORS(getUsers))
	http.HandleFunc("/api/users/", withCORS(handleUserByID)) // With trailing slash
	http.HandleFunc("/api/posts", withCORS(getPosts))
	http.HandleFunc("/api/login", withCORS(handleLogin))
	http.HandleFunc("/api/logout", withCORS(handleLogout))
	http.HandleFunc("/api/me", withCORS(handleMe))
	http.HandleFunc("/api/posts/create", withCORS(handleCreatePost))

	fmt.Printf("Backend running on port %s, allowing requests from %s\n", port, frontendURL)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
