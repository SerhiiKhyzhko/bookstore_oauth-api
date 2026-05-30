package app

import (
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/controllers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func urlMapping(router *gin.Engine, oauthCtrl *controllers.TokenHandler, appEnv string) {
	router.POST("/oauth/create", oauthCtrl.Create)
	router.POST("/oauth/refresh", oauthCtrl.RefreshToken)
	
	if appEnv == "development" {
		router.POST("/oauth/verify", oauthCtrl.VerifyToken)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
}
