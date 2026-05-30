package app

import (
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/v2/src/controllers"
	"github.com/gin-gonic/gin"
)

func StartApplication(port string, oauthCtrl *controllers.TokenHandler, appEnv string) {
	router := gin.Default()
	urlMapping(router, oauthCtrl, appEnv)	

	router.Run(port)
}
