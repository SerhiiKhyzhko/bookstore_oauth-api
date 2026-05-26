package jwtutil
 
import (
	"testing"
	"time"
 
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/stretchr/testify/assert"
)
 
func setUp() *JwtManager {
	loggerCfg := logger.Config{
		Level:       "info",
		OutputPaths: []string{"stdout"},
	}
	log, _ := logger.NewLogger(loggerCfg)
 
	return &JwtManager{
		secret:     "test-secret-key-at-least-32-characters-long",
		accessExp:  15 * time.Minute,
		refreshExp: 24 * time.Hour,
		logger: log,
	}
}
 
func TestGenerateAndVerifyAccessToken(t *testing.T) {
	manager := setUp()
 
	t.Run("Success_RoundTrip", func(t *testing.T) {
		tokenString, err := manager.GenerateAccessToken(123)
 
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
 
		claims, err := manager.VerifyToken(tokenString)
 
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, int64(123), claims.UserId)
		assert.Equal(t, ClaimsAccess, claims.TokenType)
		assert.Equal(t, Issuer, claims.Issuer)
		assert.True(t, claims.ExpiresAt.After(time.Now()))
	})
}
 
func TestGenerateAndVerifyRefreshToken(t *testing.T) {
	manager := setUp()
 
	t.Run("Success_RoundTrip", func(t *testing.T) {
		tokenString, err := manager.GenerateRefreshToken(456)
 
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
 
		claims, err := manager.VerifyToken(tokenString)
 
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, int64(456), claims.UserId)
		assert.Equal(t, ClaimsRefresh, claims.TokenType)
		assert.Equal(t, Issuer, claims.Issuer)
	})
}
 
func TestVerifyToken(t *testing.T) {
	manager := setUp()
 
	t.Run("ExpiredToken", func(t *testing.T) {
		// від'ємний час — токен вже протермінований при генерації
		expiredManager := &JwtManager{
			secret:    manager.secret,
			accessExp: -1 * time.Hour,
		}
		tokenString, err := expiredManager.GenerateAccessToken(123)
		assert.NoError(t, err)
 
		claims, err := manager.VerifyToken(tokenString)
 
		assert.Nil(t, claims)
		assert.Error(t, err)
		assert.ErrorIs(t, err, oauth_errors.UnauthorizedErr)
		assert.ErrorContains(t, err, "token expired")
	})
 
	t.Run("InvalidSignature", func(t *testing.T) {
		// генеруємо токен з іншим секретом
		otherManager := &JwtManager{
			secret:    "other-secret-key-at-least-32-characters-long",
			accessExp: 15 * time.Minute,
		}
		tokenString, err := otherManager.GenerateAccessToken(123)
		assert.NoError(t, err)
 
		claims, err := manager.VerifyToken(tokenString)
 
		assert.Nil(t, claims)
		assert.Error(t, err)
		assert.ErrorIs(t, err, oauth_errors.UnauthorizedErr)
		assert.ErrorContains(t, err, "invalid token signature")
	})
 
	t.Run("MalformedToken", func(t *testing.T) {
		claims, err := manager.VerifyToken("this.is.not.a.jwt")
 
		assert.Nil(t, claims)
		assert.Error(t, err)
		assert.ErrorIs(t, err, oauth_errors.UnauthorizedErr)
	})
 
	t.Run("EmptyToken", func(t *testing.T) {
		claims, err := manager.VerifyToken("")
 
		assert.Nil(t, claims)
		assert.Error(t, err)
		assert.ErrorIs(t, err, oauth_errors.UnauthorizedErr)
	})
}