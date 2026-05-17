package main

import (
	"log"
	"os"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/app"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/clients/cassandra"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/clients/users"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/config"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/controllers"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/repository/db"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/services/access_token"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/joho/godotenv"
	_ "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/docs"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.Load()
	loggerCfg := logger.Config{
		Level:       cfg.Logger.Level,
		OutputPaths: []string{cfg.Logger.OutputPath},
	}
	logger, err := logger.NewLogger(loggerCfg)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}
	defer logger.Sync()

	cSession, err := cassandra.NewSession(cfg.Db.Host, cfg.Db.Keyspace, cfg.Db.Consistency)
	if err != nil {
		logger.Error(err.Error(), err)
		os.Exit(1)
	}

	restyClient := users.NewRestyClient(cfg.App.RestyReqTime)
	usersClient := users.NewClient(restyClient, logger, cfg.App.ApiBaseUrl)
	dbRepository := db.NewRepository(cSession, logger)
	atService := accesstoken.NewService(usersClient, dbRepository, cfg.App.CtxTimeout)
	handler := controllers.NewHandler(atService)

	app.StartApplication(cfg.App.GinPort, handler)
}
