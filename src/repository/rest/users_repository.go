package rest

import (
	"time"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/errors"
	"github.com/go-resty/resty/v2"
)

var usersRestClient = resty.New().
	SetTimeout(150 * time.Millisecond)

type RestUserRepository interface {
	LoginUsesr(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct {}

func NewRepository() RestUserRepository {
	return &usersRepository{}
}

func (u *usersRepository) LoginUsesr(email string, password string) (*users.User, *errors.RestErr) {
	var user users.User
	var responseErr errors.RestErr

	request := users.UserLoginRequest{
		Email: email,
		Password: password,
	}

	response, err := usersRestClient.R().
	SetBody(request).
	SetResult(&user).
	SetError(&responseErr).
	Post("https://api.bookstore.com/users/login")

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())//posibly network or timeout error
	}

	if response.IsError() {
		if responseErr.Status == 404 {
			return  nil, errors.NewNotFoundError(responseErr.Message)
		} else {
			return nil, errors.NewInternalServerError(responseErr.Message)
		}
	}

	return &user, nil
}