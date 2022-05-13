package tools

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetDotEnvVariable(key string, defaultValue string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
		return defaultValue
	}

	return os.Getenv(key)
}
