package token

import (
	"fmt"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
)

type TokenRequest struct {
	Scope string `json:"scope"`

	//used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`
}

func (at *TokenRequest) Validate() error {
	if at.Username == "" || at.Password == "" {
		return fmt.Errorf("%w: invalid username or password", oauth_errors.BadRequestErr)
	}
	return nil
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserId       int64  `json:"user_id"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type TokenClaimsResponse struct {
	UserId    int64  `json:"user_id"`
	TokenType string `json:"token_type"`
	Issuer    string `json:"issuer"`
	ExpiresAt int64  `json:"expires_at"`
}
