package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	BotAdmin         int64
	BotToken         string

	ServiceName string
	LoggerLevel string
	HTTPPort    string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error !", err.Error())
	}
	cfg := Config{}

	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	cfg.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", "5432"))
	cfg.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "your user"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "your password"))
	cfg.PostgresDB = cast.ToString(getOrReturnDefault("POSTGRES_DB", "your database"))
	cfg.BotAdmin = cast.ToInt64(getOrReturnDefault("BOT_ADMIN", "your telegram id"))
	cfg.BotToken = cast.ToString(getOrReturnDefault("BOT_TOKEN", "bot token"))
	cfg.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "store"))
	cfg.LoggerLevel = cast.ToString(getOrReturnDefault("LOGGER_LEVEL", "debug"))
	cfg.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":7070"))
	return cfg
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}
