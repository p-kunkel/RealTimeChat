package config

import (
	"github.com/joho/godotenv"
)

func LoadEnvToCache() error {
	return godotenv.Load()
}
