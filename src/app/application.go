package app

import (
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/clients/cassandra"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/http"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/repository/db"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func StartApplication() {
	session, dbErr := cassandra.GetSession()
	if dbErr != nil {
		panic(dbErr)
	}
	session.Close()

	atService := accesstoken.NewService(db.NewRepository())
	atHendler := http.NewHandler(atService)

	router.GET("/oauth/access_token/:access_token_id", atHendler.GetById)

	router.Run(":8080")
}