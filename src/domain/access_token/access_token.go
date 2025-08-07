package accesstoken

import (
	"strings"
	"time"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/errors"
)

const expirationTime = 24

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expiers     int64  `json:"expiers"`
}

func(at AccessToken) Validate() *errors.RestErr{
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return errors.NewBadRequestError("invalid access token id")
	}
	if at.Expiers <= 0 {
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

func GetAccessToken() *AccessToken {
	return &AccessToken{
		Expiers: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expiers, 0).Before(time.Now().UTC())
}