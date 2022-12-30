package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment         string // develop, staging, production
	PosgresHost         string
	PostgresPort        int
	PostgresDatabase    string
	PostgresUser        string
	PostgresPassword    string
	ReviewServiceHost   string
	ReviewServicePort   int
	CustomerServiceHost string
	CustomerServicePort int
	LogLevel            string
	PostServicePort     string
	PostServiceHost     string
}

func Load() Config {
	c := Config{}
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.PosgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "postdb"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "developer"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "2002"))
	c.ReviewServiceHost = cast.ToString(getOrReturnDefault("REVIEW_SERVICE_HOST", "review"))
	c.ReviewServicePort = cast.ToInt(getOrReturnDefault("REVIEW_SERVICE_PORT", 5000))
	c.CustomerServiceHost = cast.ToString(getOrReturnDefault("CUSTOMER_SERVICE_HOST", "customer"))
	c.CustomerServicePort = cast.ToInt(getOrReturnDefault("CUSTOMER_SERVICE_PORT", 3000))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOGLEVEL", "debug"))
	c.PostServicePort = cast.ToString(getOrReturnDefault("POST_SERVICE_PORT", "7000"))
	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post"))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
