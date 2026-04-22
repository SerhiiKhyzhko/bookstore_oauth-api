package app

import (
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/controllers"
	"github.com/gin-gonic/gin"
)

func urlMapping(router *gin.Engine, oauthCtrl *controllers.AccessTokenHandler) {
	router.GET("/oauth/access_token/:access_token_id", oauthCtrl.GetById)
	router.POST("/oauth/access_token", oauthCtrl.Create)
	router.PATCH("/oauth/access_token/:access_token_id", oauthCtrl.UpdateExpirationTime)
}
