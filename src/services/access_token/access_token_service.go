package accesstoken

import (
	"context"
	"fmt"
	"strings"
	"time"

	accesstoken "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
	atrepository "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/at_repository"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
)

type Service interface {
	GetById(context.Context, string) (*accesstoken.AccessToken, error)
	Create(context.Context, accesstoken.AccessTokenRequest) (*accesstoken.AccessToken, error)
	UpdateExpirationTime(context.Context, accesstoken.AccessToken) error
}

type RestUserClient interface {
	LoginUser(string, string) (int64, error)
}

type service struct {
	restUsersClient RestUserClient
	dbRepo          atrepository.DbRepository
	ctxTimeout      time.Duration
}

func NewService(usersClient RestUserClient, dbRepo atrepository.DbRepository, timeout time.Duration) Service {
	return &service{
		restUsersClient: usersClient,
		dbRepo:          dbRepo,
		ctxTimeout:      timeout,
	}
}

func (s *service) GetById(ctx context.Context, accessTokenId string) (*accesstoken.AccessToken, error) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, fmt.Errorf("%w: invalid access token id",  oauth_errors.BadRequestErr)
	}
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	accessToken, err := s.dbRepo.GetById(ctx, accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(ctx context.Context, request accesstoken.AccessTokenRequest) (*accesstoken.AccessToken, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	//TODO: Support both grant types: client_credentials

	id, err := s.restUsersClient.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	at := accesstoken.GetNewAccessToken(id)
	at.Generate()

	if err := s.dbRepo.Create(ctx, at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(ctx context.Context, at accesstoken.AccessToken) error {
	if err := at.Validate(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, s.ctxTimeout)
	defer cancel()

	return s.dbRepo.UpdateExpirationTime(ctx, at)
}
