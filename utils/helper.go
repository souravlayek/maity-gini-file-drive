package utils

import (
	"fmt"

	"github.com/joho/godotenv"
)

// LoadENV loads the environment variables from .env file
func LoadENV() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}
