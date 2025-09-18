package accesstoken

import (
	"fmt"
	"strings"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/crypto_utils"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/errors"
)

const (
	expirationTime = 24
	grantTypePassword = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope string `json:"scope"`

	//used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// used for client credentials grant type
	ClientId string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (at *AccessTokenRequest) Validate() *errors.RestErr {
	switch at.GrantType {
	case grantTypePassword:
		if at.Username == "" || at.Password == "" {
			return errors.NewBadRequestError("invalid username or passwprd")
		}
	case grantTypeClientCredentials:
		if at.ClientId == "" || at.ClientSecret == "" {
			return errors.NewBadRequestError("invalid client_id or client_secret")
		}
	default:
		return errors.NewBadRequestError("invalid grant_type paramether")
	}
	return nil
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expires     int64  `json:"expiers"`
}

func(at AccessToken) Validate() *errors.RestErr{
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestError("invalid access token id")
	}
	if at.Expires <= 0 {
		return errors.NewBadRequestError("invalid expiration time")
	}

	return nil
}

func(at AccessToken) ValidateAll() *errors.RestErr{
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	
	at.Validate()
	if at.UserId <= 0 {
		return errors.NewBadRequestError("invalid user id")
	}
	if at.ClientId <= 0 {
		return errors.NewBadRequestError("invalid client id")
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