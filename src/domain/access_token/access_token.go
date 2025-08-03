package accesstoken

import (
	"time"
)

const expirationTime = 24

type AccessToken struct {
	AccessToken string `json:"access_token"`
	UserId      int64  `json:"user_id"`
	ClientId    int64  `json:"client_id"`
	Expiers     int64  `json:"expiers"`
}

func GetAccessToken() *AccessToken {
	return &AccessToken{
		Expiers: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expiers, 0).Before(time.Now().UTC())
}