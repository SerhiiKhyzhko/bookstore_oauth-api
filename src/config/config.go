package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App    appConfig
	Jwt    jwtConfig
	Logger loggerConfig
}

type appConfig struct {
	GinPort      string
	ApiBaseUrl   string
	CtxTimeout   time.Duration
	AppEnv       string
}

type jwtConfig struct {
	SignKey           string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

type loggerConfig struct {
	Level      string
	OutputPath string
}

func loadApp() appConfig {
	app := appConfig{}
	app.GinPort = getRequiredEnv("GIN_PORT")
	app.ApiBaseUrl = getRequiredEnv("USERS_API_BASE_URL")
	app.CtxTimeout = getTimeWithDefault("CTX_TIMEOUT", "2s")
	app.AppEnv = getRequiredEnv("APP_ENV")
	return app
}

func loadJwt() jwtConfig {
	jwt := jwtConfig{}
	jwt.AccessExpiration = getTimeWithDefault("ACCESS_EXPIRATION", "15m")
	jwt.RefreshExpiration = getTimeWithDefault("REFRESH_EXPIRATION", "24h")
	jwt.SignKey = getRequiredEnv("SIGN_KEY")
	return jwt
}

func loadLogger() loggerConfig {
	logger := loggerConfig{}
	logger.Level = getRequiredEnv("LEVEL")
	logger.OutputPath = getRequiredEnv("OUTPUT_PATHS")
	return logger
}

func getRequiredEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Critical environment variable %s is missing", key)
	}
	return value
}

func getIntWithDefault(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid value for %s", key)
	}
	return result
}

func getTimeWithDefault(key string, defaultValue string) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		res, err := time.ParseDuration(defaultValue)
		if err != nil {
			log.Fatalf("convertation of default value failed: %v", err)
		}
		return res
	}
	result, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid value for %s", key)
	}
	return result
}

func Load() Config {
	cfg := Config{}
	cfg.App = loadApp()
	cfg.Jwt = loadJwt()
	cfg.Logger = loadLogger()
	return cfg
}
