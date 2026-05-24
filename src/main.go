package main

import (
	"log"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/app"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/clients/users"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/config"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/controllers"
	_ "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/docs"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/internal/jwtutil"
	tokenservice "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/services/token_service"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/joho/godotenv"
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

	jwtManager := jwtutil.NewJwtManager(cfg.Jwt.SignKey, cfg.Jwt.AccessExpiration, cfg.Jwt.RefreshExpiration, logger)
	restyClient := users.NewRestyClient(cfg.App.RestyReqTime)
	usersClient := users.NewClient(restyClient, logger, cfg.App.ApiBaseUrl)
	atService := tokenservice.NewService(usersClient, jwtManager)
	handler := controllers.NewHandler(atService, cfg.App.CtxTimeout)

	app.StartApplication(cfg.App.GinPort, handler)
}
