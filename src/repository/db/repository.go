package db

import (
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/clients/cassandra"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/errors"
	"github.com/gocql/gocql"
)

const(
	querryAccessToken = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	querryCreateAccessToken = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	querryUpdate = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
) 

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*accesstoken.AccessToken, *errors.RestErr)
	Create(accesstoken.AccessToken) *errors.RestErr
	UpdateExpirationTime(accesstoken.AccessToken) *errors.RestErr
}

type dbRepository struct {}

func (r *dbRepository) GetById(id string) (*accesstoken.AccessToken, *errors.RestErr) {
	var result accesstoken.AccessToken

	if err := cassandra.GetSession().Query(querryAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserId,
		&result.ClientId,
		&result.Expiers,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("no access token found with given id")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}

	return &result, nil
}

func (r *dbRepository) Create(at accesstoken.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(querryCreateAccessToken, at.AccessToken, at.UserId, at.ClientId, at.Expiers).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}

func (r *dbRepository) UpdateExpirationTime(at accesstoken.AccessToken) *errors.RestErr {
	if err := cassandra.GetSession().Query(querryUpdate, at.Expiers, at.AccessToken).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	return nil
}