package controllers

import (
	"errors"
	"net/http"

	atDomain "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/services/access_token"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

type AccessTokenHandler struct {
	service accesstoken.Service
}

func NewHandler(service accesstoken.Service) *AccessTokenHandler {
	return &AccessTokenHandler{
		service: service,
	}
}

func requestError(reqErr error) rest_errors.RestErr {
	switch {
	case errors.Is(reqErr, oauth_errors.NotFoundErr):
		return rest_errors.NewNotFoundError(errors.Unwrap(reqErr).Error())
	case errors.Is(reqErr, oauth_errors.BadRequestErr):
		return rest_errors.NewBadRequestError(errors.Unwrap(reqErr).Error())
	default:
		return rest_errors.NewInternalServerError(errors.Unwrap(reqErr).Error(), errors.Unwrap(reqErr))
	}
}

func (handler *AccessTokenHandler) GetById(c *gin.Context) {
	accessTokenId := c.Param("access_token_id")
	ctx := c.Request.Context()
	accessToken, err := handler.service.GetById(ctx, accessTokenId)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}

	c.JSON(http.StatusOK, accessToken)
}

func (handler *AccessTokenHandler) Create(c *gin.Context) {
	var request atDomain.AccessTokenRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	accessToken, err := handler.service.Create(ctx, request)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}

	c.JSON(http.StatusCreated, accessToken)
}

func (handler *AccessTokenHandler) UpdateExpirationTime(c *gin.Context) {
	var at atDomain.AccessToken
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	if err := handler.service.UpdateExpirationTime(ctx, at); err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}

	c.JSON(http.StatusOK, at)
}
