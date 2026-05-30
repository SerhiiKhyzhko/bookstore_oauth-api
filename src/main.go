package main

import (
	"log"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/app"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/clients/users"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/config"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/controllers"
	_ "github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/docs"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/internal/jwtutil"
	tokenservice "github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/services/token_service"
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
	restyClient := users.NewRestyClient()
	usersClient := users.NewClient(restyClient, logger, cfg.App.ApiBaseUrl)
	atService := tokenservice.NewService(usersClient, jwtManager)
	handler := controllers.NewHandler(atService, cfg.App.CtxTimeout, logger)

	app.StartApplication(cfg.App.GinPort, handler, cfg.App.AppEnv)
}
