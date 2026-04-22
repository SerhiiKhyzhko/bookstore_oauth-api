package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App    appConfig
	Db     dbConfig
	Logger loggerConfig
}

type appConfig struct {
	GinPort      string
	RestyReqTime int
	ApiBaseUrl   string
	CtxTimeout   time.Duration
}

type dbConfig struct {
	Host        string
	Keyspace    string
	Consistency string
}

type loggerConfig struct {
	Level      string
	OutputPath string
}

func loadApp() appConfig {
	app := appConfig{}
	app.GinPort = getRequiredEnv("GIN_PORT")
	app.RestyReqTime = getIntWithDefault("RESTY_REQUEST_TIME", 150)
	app.ApiBaseUrl = getRequiredEnv("USERS_API_BASE_URL")
	app.CtxTimeout = getTimeWithDefault("CTX_TIMEOUT", "2s")
	return app
}

func loadDb() dbConfig {
	db := dbConfig{}
	db.Host = getRequiredEnv("DB_HOST")
	db.Keyspace = getRequiredEnv("KEYSPACE")
	db.Consistency = getEnvWithDefault("CONSISTENCY", "Quorum")
	return db
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

func getEnvWithDefault(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
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
	cfg.Db = loadDb()
	cfg.Logger = loadLogger()
	return cfg
}
