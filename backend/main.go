package main

import (
	"backend/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func setHandlers() {
	// CORS to allow developement on same address
	http.HandleFunc("/api/users", middleware.WithCORS(handlers.GetUsers))
	http.HandleFunc("/api/users/", middleware.WithCORS(handlers.HandleUserByID)) // With trailing slash
	http.HandleFunc("/api/users/search", middleware.WithCORS(handlers.SearchUsers))
	http.HandleFunc("/api/posts", middleware.WithCORS(handlers.GetPosts))
	http.HandleFunc("/api/login", middleware.WithCORS(handlers.HandleLogin))
	http.HandleFunc("/api/logout", middleware.WithCORS(handlers.HandleLogout))
	http.HandleFunc("/api/me", middleware.WithCORS(handlers.HandleMe))
	http.HandleFunc("/api/posts/create", middleware.WithCORS(handlers.HandleCreatePost))
}

func main() {
	config.InitConfig()

	err := database.NewDatabase(config.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	setHandlers()
	fmt.Printf("Backend running on port %s, allowing requests from %s\n", config.Port, config.FrontendURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
