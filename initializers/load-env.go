package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("VERCEL") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}
