package config

import (
	"os"
)

const DBPath = "data/app.db"

var (
	Port        string
	FrontendURL string
)

func InitConfig() {
	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080" // Default port
	}

	FrontendURL = os.Getenv("FRONTEND_URL")
	if FrontendURL == "" {
		FrontendURL = "http://localhost:5173" // Default frontend URL
	}
}
