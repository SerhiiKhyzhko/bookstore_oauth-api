package app

import (
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/controllers"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func urlMapping(router *gin.Engine, oauthCtrl *controllers.AccessTokenHandler) {
	router.POST("/oauth/create", oauthCtrl.Create)
	router.POST("/oauth/refresh", oauthCtrl.RefreshToken)
	router.POST("/oauth/verify", oauthCtrl.VerifyToken)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
