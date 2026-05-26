package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/token"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
	tokenservice "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/services/token_service"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	service    tokenservice.Service
	ctxTimeout time.Duration
}

func NewHandler(service tokenservice.Service, timeout time.Duration) *TokenHandler {
	return &TokenHandler{
		service:    service,
		ctxTimeout: timeout,
	}
}

func requestError(reqErr error) rest_errors.RestErr {
	switch {
	case errors.Is(reqErr, oauth_errors.RequestTimeoutErr):
		return rest_errors.NewRestError("request timeout", http.StatusRequestTimeout, "database error", nil)
	case errors.Is(reqErr, oauth_errors.UnauthorizedErr):
		return rest_errors.NewUnauthorizedError(errors.Unwrap(reqErr).Error())
	case errors.Is(reqErr, oauth_errors.NotFoundErr):
		return rest_errors.NewNotFoundError(errors.Unwrap(reqErr).Error())
	case errors.Is(reqErr, oauth_errors.BadRequestErr):
		return rest_errors.NewBadRequestError(errors.Unwrap(reqErr).Error())
	default:
		return rest_errors.NewInternalServerError(errors.Unwrap(reqErr).Error(), errors.Unwrap(reqErr))
	}
}

// @Summary     Create new access and refresh tokens
// @Tags        token
// @Description Generate and return new access and refresh tokens with provided information
// @Accept      json
// @Produce     json
// @Param       request body token.TokenRequest true "Token Request"
// @Success     201 {object} token.Token
// @Failure     400 {object} oauth_errors.SwaggerRestErr
// @Failure     404 {object} oauth_errors.SwaggerRestErr
// @Failure     408 {object} oauth_errors.SwaggerRestErr
// @Failure     500 {object} oauth_errors.SwaggerRestErr
// @Router      /oauth/create [post]
func (handler *TokenHandler) Create(c *gin.Context) {
	var request token.TokenRequest
	ctx, cancel := context.WithTimeout(c.Request.Context(), handler.ctxTimeout)
	defer cancel()
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	token, err := handler.service.Create(ctx, request)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}
	
	c.JSON(http.StatusCreated, token)
}

// @Summary     Generate new access token based on a refresh token
// @Tags        token
// @Description Parse refresh token and validate claims, if all is ok, generate new access token
// @Accept      json
// @Produce     json
// @Param       request body token.RefreshToken true "Refresh Token"
// @Success     200 {object} token.AccessToken
// @Failure     400 {object} oauth_errors.SwaggerRestErr
// @Failure     401 {object} oauth_errors.SwaggerRestErr
// @Failure     500 {object} oauth_errors.SwaggerRestErr
// @Router      /oauth/refresh [post]
func (handler *TokenHandler) RefreshToken(c *gin.Context) {
	var rt token.RefreshToken
	if err := c.ShouldBindJSON(&rt); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	at, err := handler.service.RefreshToken(rt.RefreshToken)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}

	c.JSON(http.StatusOK, at)
}

// @Summary     Validate given token
// @Tags        token
// @Description Parse token and validate claims, if all is ok, return claims
// @Accept      json
// @Produce     json
// @Param       request body token.VerifyRequest true "Verify Request"
// @Success     200 {object} jwtutil.TokenClaims
// @Failure     400 {object} oauth_errors.SwaggerRestErr
// @Failure     401 {object} oauth_errors.SwaggerRestErr
// @Failure     500 {object} oauth_errors.SwaggerRestErr
// @Router      /oauth/verify [post]
func (handler *TokenHandler) VerifyToken(c *gin.Context) {
	var token token.VerifyRequest
	if err := c.ShouldBindJSON(&token); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	claims, err := handler.service.VerifyToken(token.Token)
	if err != nil {
		restErr := requestError(err)
		c.JSON(restErr.Status(), restErr)
		return
	}

	c.JSON(http.StatusOK, claims)
}
