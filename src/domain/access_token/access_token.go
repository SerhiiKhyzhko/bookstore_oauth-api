package accesstoken

import (
	"fmt"
	"strings"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/crypto_utils"
)

const (
	expirationTime             = 24
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	//used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// used for client credentials grant type
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() error {
	switch at.GrantType {
	case grantTypePassword:
		if at.Username == "" || at.Password == "" {
			return fmt.Errorf("%w: invalid username or password", oauth_errors.BadRequestErr)
		}
	case grantTypeClientCredentials:
		if at.ClientId == "" || at.ClientSecret == "" {
			return fmt.Errorf("%w: invalid client_id or client_secret", oauth_errors.BadRequestErr)
		}
	default:
		return fmt.Errorf("%w: invalid grant_type paramether", oauth_errors.BadRequestErr)
	}
	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expires"`
}

func (at *AccessToken) Validate() error {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return fmt.Errorf("%w: invalid access token id", oauth_errors.BadRequestErr)
	}
	if at.Expires <= 0 {
		return fmt.Errorf("%w: invalid expiration time", oauth_errors.BadRequestErr)
	}

	return nil
}

func (at *AccessToken) ValidateAll() error {
	at.AccessToken = strings.TrimSpace(at.AccessToken)

	if err := at.Validate(); err != nil {
		return err
	}
	if at.UserId <= 0 {
		return fmt.Errorf("%w, invalid user id", oauth_errors.BadRequestErr)
	}
	if at.ClientId <= 0 {
		return fmt.Errorf("%w, invalid client id", oauth_errors.BadRequestErr)
	}
	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserId:  userId,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}
