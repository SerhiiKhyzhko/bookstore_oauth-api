package users

import (
	"fmt"
	"net/http"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
	accesstoken "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/services/access_token"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/rest_errors"
	"github.com/go-resty/resty/v2"
)

type usersClient struct {
	restClient *resty.Client
	logger     *logger.Logger
	apiBaseUrl string
}

func NewClient(client *resty.Client, logger *logger.Logger, url string) accesstoken.RestUserClient {
	return &usersClient{
		restClient: client,
		logger:     logger,
		apiBaseUrl: url,
	}
}

func (u *usersClient) LoginUser(email string, password string) (int64, error) {
	var user User
	var responseErr rest_errors.RestErr

	request := UserLoginRequest{
		Email:    email,
		Password: password,
	}

	response, err := u.restClient.R().
		SetBody(request).
		SetResult(&user).
		SetError(&responseErr).
		Post(u.apiBaseUrl)

	if err != nil {
		u.logger.Error("request failed", err)
		return 0, oauth_errors.InternalServerErr
	}

	if response.IsError() {
		if responseErr.Status() == http.StatusNotFound {
			return 0, fmt.Errorf("%w: %s", oauth_errors.NotFoundErr, responseErr.Message())
		} else {
			return 0, fmt.Errorf("%w: %s", oauth_errors.InternalServerErr, responseErr.Message())
		}
	}

	return user.Id, nil
}
