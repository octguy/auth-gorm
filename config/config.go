package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port      string
	JWTSecret string
	GinMode   string
}

func Load() *Config {
	godotenv.Load() // Load .env file if it exists

	return &Config{
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		GinMode:   os.Getenv("GIN_MODE"),
	}
}
