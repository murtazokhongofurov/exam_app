package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string // develop, staging, production

	CustomerServiceHost string
	CustomerServicePort int
	PostServiceHost     string
	PostServicePort     int
	ReviewServiceHost   string
	ReviewServicePort   int
	// PostgresHost        string
	// PosgresPort         int
	// PostgresUser        string
	// PostgresDatabase    string
	// PostgresPassword    string
	AuthConfigPath      string
	CSVFilePath         string
	RedisHost           string
	RedisPort           string
	CtxTimeout          int
	LogLevel            string
	HTTPPort            string
	SignInKey           string
}

func Load() Config {
	c := Config{}
	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	// c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	// c.PosgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	// c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "developer"))
	// c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "customer"))
	// c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "2002"))

	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":9090"))

	c.CustomerServiceHost = cast.ToString(getOrReturnDefault("CUSTOMER_SERVICE_HOST", "customer"))
	c.CustomerServicePort = cast.ToInt(getOrReturnDefault("CUSTOMER_SERVICE_PORT", 3000))

	c.PostServiceHost = cast.ToString(getOrReturnDefault("POST_SERVICE_HOST", "post"))
	c.PostServicePort = cast.ToInt(getOrReturnDefault("POST_SERVICE_PORT", 7000))

	c.ReviewServiceHost = cast.ToString(getOrReturnDefault("REVIEW_SERVICE_HOST", "review"))
	c.ReviewServicePort = cast.ToInt(getOrReturnDefault("REVIEW_SERVICE_PORT", 5000))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "redis"))
	c.RedisPort = cast.ToString(getOrReturnDefault("REDIS_PORT", "6379"))

	c.SignInKey = cast.ToString(getOrReturnDefault("SINGINGKEY", "develop_2002"))
	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIME_OUT", 7))

	c.CSVFilePath = cast.ToString(getOrReturnDefault("CSV_FILE_PATH","./config/casbin_rules.csv"))
	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_PATH", "./config/auth.conf"))
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
