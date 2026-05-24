package jwtutil

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/golang-jwt/jwt/v5"
)

const (
	claimsRefresh = "refresh"
	claimsAccess  = "access"
	issuer        = "bookstore-oauth-api"
)
const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	claimsIdLength = 10
)

type JwtManager struct {
	secret     string
	accessExp  time.Duration
	refreshExp time.Duration
	logger     *logger.Logger
}

func NewJwtManager(secret string, accessExp time.Duration, refreshExp time.Duration, logger *logger.Logger) *JwtManager {
	return &JwtManager{
		secret:     secret,
		accessExp:  accessExp,
		refreshExp: refreshExp,
		logger:     logger,
	}
}

func idGenerator() string {
	b := make([]byte, claimsIdLength)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func (j *JwtManager) generateToken(userId int64, tokenType string, exp time.Duration) (string, error) {
	claims := TokenClaims{
		UserId:    userId,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        idGenerator(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(j.secret))
	if err != nil {
		j.logger.Error(err.Error(), err)
		return "", fmt.Errorf("%w: jwt token generation failed", oauth_errors.InternalServerErr)
	}
	return ss, nil
}

func (j *JwtManager) GenerateAccessToken(userId int64) (string, error) {
	return j.generateToken(userId, claimsAccess, j.accessExp)
}

func (j *JwtManager) GenerateRefreshToken(userId int64) (string, error) {
	return j.generateToken(userId, claimsRefresh, j.refreshExp)
}

func (j *JwtManager) VerifyToken(tokenString string) (*TokenClaims, error) {
	claims := TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(j.secret), nil
	},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		j.logger.Error(err.Error(), err)
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, fmt.Errorf("%w: token expired", oauth_errors.UnauthorizedErr)
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, fmt.Errorf("%w: token not valid yet", oauth_errors.UnauthorizedErr)
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, fmt.Errorf("%w: invalid token signature", oauth_errors.UnauthorizedErr)
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, fmt.Errorf("%w: malformed token", oauth_errors.UnauthorizedErr)
		default:
			return nil, fmt.Errorf("%w: invalid token", oauth_errors.UnauthorizedErr)
		}
	}

	if !token.Valid {
    	return nil, fmt.Errorf("%w: invalid token", oauth_errors.UnauthorizedErr)
	}

	return &claims, nil
}
