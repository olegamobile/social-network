package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	// Import the godotenv library
)

const DBPath = "data/db/app.db"
const MigrationsPath = "migrations"

var (
	Port        string
	FrontendURL string
)

func InitConfig() {
	err := loadEnvFile("./config/.env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

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

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Ignore comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		// Split by first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // or return error if preferred
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		value = strings.Trim(value, `"'`)

		os.Setenv(key, value)
	}

	return scanner.Err()
}
