package internal

import (
	"github.com/joho/godotenv"
)

func LoadDotEnv(filePath string) {
	godotenv.Load(filePath)
}
