package users

import (
	"github.com/go-resty/resty/v2"
)

func NewRestyClient() *resty.Client {
	restClient := resty.New()
	return restClient
}
