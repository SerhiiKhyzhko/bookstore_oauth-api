package accesstoken

import (
	"strings"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/repository/rest"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/repository/db"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/errors"
)

type Repository interface {
	GetById(string) (*accesstoken.AccessToken, *errors.RestErr)
	Create(accesstoken.AccessToken) *errors.RestErr
	UpdateExpirationTime(accesstoken.AccessToken) *errors.RestErr
}

type Service interface {
	GetById(string) (*accesstoken.AccessToken, *errors.RestErr)
	Create(accesstoken.AccessTokenRequest) (*accesstoken.AccessToken, *errors.RestErr)
	UpdateExpirationTime(accesstoken.AccessToken) *errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUserRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUserRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*accesstoken.AccessToken, *errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request accesstoken.AccessTokenRequest) (*accesstoken.AccessToken, *errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	//TODO: Support both grant types: client_credentials and password

	// Authenticate the user against the Users API:
	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate a new access token:
	at := accesstoken.GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new access token in Cassandra:
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at accesstoken.AccessToken) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}

	return s.dbRepo.UpdateExpirationTime(at)
}