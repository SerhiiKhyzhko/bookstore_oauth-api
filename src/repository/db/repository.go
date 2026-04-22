package db

import (
	"context"
	"fmt"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
	atrepository "github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/at_repository"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
	"github.com/SerhiiKhyzhko/bookstore_utils-go/logger"
	"github.com/gocql/gocql"
)

const (
	querryAccessToken       = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	querryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	querryUpdate            = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

func NewRepository(cSession *gocql.Session, logger *logger.Logger) atrepository.DbRepository {
	return &dbRepository{
		session: cSession,
		logger:  logger,
	}
}

type dbRepository struct {
	session *gocql.Session
	logger  *logger.Logger
}

func (r *dbRepository) GetById(ctx context.Context, id string) (*accesstoken.AccessToken, error) {
	var result accesstoken.AccessToken

	if err := r.session.Query(querryAccessToken, id).WithContext(ctx).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, fmt.Errorf("%w: no access token found with given id", oauth_errors.NotFoundErr)
		}
		r.logger.Error("Request failed", err)
		return nil, oauth_errors.InternalServerErr
	}

	return &result, nil
}

func (r *dbRepository) Create(ctx context.Context, at accesstoken.AccessToken) error {
	if err := r.session.Query(querryCreateAccessToken, at.AccessToken, at.UserId, at.ClientId, at.Expires).WithContext(ctx).Exec(); err != nil {
		r.logger.Error(fmt.Sprintf("Error when trying to create new token: %v", at), err)
		return oauth_errors.InternalServerErr
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(ctx context.Context, at accesstoken.AccessToken) error {
	if err := r.session.Query(querryUpdate, at.Expires, at.AccessToken).WithContext(ctx).Exec(); err != nil {
		r.logger.Error(fmt.Sprintf("Error when trying to update new token: %v", at), err)
		return oauth_errors.InternalServerErr
	}

	return nil
}
