package http

import (
	"net/http"
	
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface{
	GetById(*gin.Context)
}	

type accesstokenHandler struct {
	service accesstoken.Service
}

func NewHandler(service accesstoken.Service) AccessTokenHandler {
	return &accesstokenHandler{
		service: service,
	}
}

func (handler * accesstokenHandler) GetById(c *gin.Context) {
	accessTokenId := c.Param("access_token_id")
	accessToken, err := handler.service.GetById(accessTokenId)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, accessToken)
}