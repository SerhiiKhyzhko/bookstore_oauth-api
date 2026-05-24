package jwtutil

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserId    int64  `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func (tc TokenClaims) Validate() error {
	if tc.TokenType != claimsAccess && tc.TokenType != claimsRefresh {
		return fmt.Errorf("Invalid token type - %s", tc.TokenType)
	}

	if tc.UserId <= 0 {
		return fmt.Errorf("Invalid user id - %d", tc.UserId)
	}
	return nil
}
