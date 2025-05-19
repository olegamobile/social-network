package config

import (
	"fmt"
	"os"
)

const DBPath = "data/app.db"
const MigrationsPath = "migrations"

var (
	Port        string
	FrontendURL string
)

func InitConfig() {
	Port = os.Getenv("PORT")
	if Port == "" {
		fmt.Println("falling back to default port")
		Port = "8080" // Default port
	}

	FrontendURL = os.Getenv("FRONTEND_URL")
	if FrontendURL == "" {
		fmt.Println("falling back to default frontend url")
		FrontendURL = "http://localhost:5173" // Default frontend URL
	}
}
