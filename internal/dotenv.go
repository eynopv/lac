package internal

import (
	"github.com/joho/godotenv"
)

func LoadDotEnv() {
	godotenv.Load(".env")
}
