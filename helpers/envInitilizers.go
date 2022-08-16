package helpers

import (
	"log"

	"github.com/joho/godotenv"
)

func EnvInitilizer() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to load environment variables")
	}
}
