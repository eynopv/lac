package internal

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Unable to load dotenv:", err)
	}
}
