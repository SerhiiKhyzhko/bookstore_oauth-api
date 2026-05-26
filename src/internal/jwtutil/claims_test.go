package jwtutil
 
import (
	"testing"
 
	"github.com/stretchr/testify/assert"
)
 
func TestTokenClaims_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		claims      TokenClaims
		expectError bool
		errContains string
	}{
		{
			name: "Success_Access",
			claims: TokenClaims{
				UserId:    1,
				TokenType: ClaimsAccess,
			},
			expectError: false,
		},
		{
			name: "Success_Refresh",
			claims: TokenClaims{
				UserId:    1,
				TokenType: ClaimsRefresh,
			},
			expectError: false,
		},
		{
			name: "InvalidTokenType",
			claims: TokenClaims{
				UserId:    1,
				TokenType: "invalid",
			},
			expectError: true,
			errContains: "Invalid token type",
		},
		{
			name: "EmptyTokenType",
			claims: TokenClaims{
				UserId:    1,
				TokenType: "",
			},
			expectError: true,
			errContains: "Invalid token type",
		},
		{
			name: "InvalidUserId",
			claims: TokenClaims{
				UserId:    0,
				TokenType: ClaimsAccess,
			},
			expectError: true,
			errContains: "Invalid user id",
		},
		{
			name: "NegativeUserId",
			claims: TokenClaims{
				UserId:    -1,
				TokenType: ClaimsAccess,
			},
			expectError: true,
			errContains: "Invalid user id",
		},
		{
			name:        "BothFieldsEmpty",
			claims:      TokenClaims{},
			expectError: true,
			errContains: "Invalid token type",
		},
	}
 
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.claims.Validate()
			if tc.expectError {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
 