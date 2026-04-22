package users

import (
	"time"

	"github.com/go-resty/resty/v2"
)

func NewRestyClient(requestTime int) *resty.Client {
	restClient := resty.New().SetTimeout(time.Duration(requestTime) * time.Millisecond)
	return restClient
}
