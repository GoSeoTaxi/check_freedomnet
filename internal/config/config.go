package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Servers    []string // Список серверов через ;
	MaxRetries int      // Максимальное количество попыток
	Port       string   // Порт для запуска сервиса
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system environment variables")
	}

	return &Config{
		Servers:    splitServers(os.Getenv("SERVERS")),
		MaxRetries: getEnvAsInt("MAX_RETRIES", 10),
		Port:       getEnvAsString("PORT", "8080"),
	}
}

func splitServers(servers string) []string {
	return strings.Split(servers, ";")
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultVal
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultVal
	}
	return value
}

func getEnvAsString(key string, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
