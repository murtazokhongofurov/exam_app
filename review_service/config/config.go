package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Envirnment        string
	PostgresHost      string
	PostgresPort      int
	PostgresDatabase  string
	PostgresUser      string
	PostgresPassword  string
	LogLevel          string
	ReviewServicePort string
	ReviewServiceHost string
	PostServiceHost   string
	PostServicePort   int

}

func Load() Config {
	c := Config{}
	c.Envirnment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "reviewdb"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "developer"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "2002"))
	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post"))
	c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_PORT", 7000))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.ReviewServicePort = cast.ToString(getOrReturnDefault("REVIEW_SERVICE_PORT", "5000"))
	c.ReviewServiceHost = cast.ToString(getOrReturnDefault("REVIEW_SERVICE_HOST", "review"))


	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
