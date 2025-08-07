package http

import (
	"net/http"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/errors"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler interface{
	GetById(*gin.Context)
	Create(*gin.Context)
	UpdateExpirationTime(*gin.Context)
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
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (handler * accesstokenHandler) Create(c *gin.Context) {
	var at accesstoken.AccessToken
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := errors.NewBadRequestError("invalid jsoon body")
		c.JSON(restErr.Status, restErr)
		return
	}

	if err := handler.service.Create(at); err != nil {
		c.JSON(err.Status, err)
		return
	}
	
	c.JSON(http.StatusCreated, at)
}

func (handler * accesstokenHandler) UpdateExpirationTime(c *gin.Context) {
	var at accesstoken.AccessToken
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := errors.NewBadRequestError("invalid jsoon body")
		c.JSON(restErr.Status, restErr)
		return
	}

	if err := handler.service.UpdateExpirationTime(at); err != nil {
		c.JSON(err.Status, err)
		return
	}
	
	c.JSON(http.StatusOK, at)
}