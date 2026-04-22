package app

import (
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/controllers"
	"github.com/gin-gonic/gin"
)

func StartApplication(port string, oauthCtrl *controllers.AccessTokenHandler) {
	router := gin.Default()
	urlMapping(router, oauthCtrl)

	router.Run(port)
}
