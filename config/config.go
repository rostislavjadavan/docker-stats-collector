package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	LogLevel            string
	DbPath              string
	WebServerAddress    string
	WebserverPort       string
	AdminPassword       string
	DockerClientVersion string
}

func LoadFromEnv() AppConfig {
	return AppConfig{
		LogLevel:            GetEnv("LOG_LEVEL", "debug"),
		DbPath:              GetEnv("DB_PATH", "./stats.db"),
		WebServerAddress:    GetEnv("WEBSERVER_ADDRESS", "127.0.0.1"),
		WebserverPort:       GetEnv("WEBSERVER_PORT", "8080"),
		AdminPassword:       GetEnv("ADMIN_PASSWORD", "admin"),
		DockerClientVersion: GetEnv("DOCKER_CLIENT_VERSION", "1.46"),
	}
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func GetEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Invalid %s parameter. Value needs to be an integer", value))
	}
	return valueAsInt
}
