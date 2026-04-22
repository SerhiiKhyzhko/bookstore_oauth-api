package atrepository

import (
	"context"

	accesstoken "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
)

type DbRepository interface {
	GetById(context.Context, string) (*accesstoken.AccessToken, error)
	Create(context.Context, accesstoken.AccessToken) error
	UpdateExpirationTime(context.Context, accesstoken.AccessToken) error
}