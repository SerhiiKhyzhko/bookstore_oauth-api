package db

import (
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/utils/errors"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/access_token"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*accesstoken.AccessToken, *errors.RestErr)
}

type dbRepository struct {}

func (r *dbRepository) GetById(id string) (*accesstoken.AccessToken, *errors.RestErr) {
	return nil, errors.NewInternalServerError("DB connection not implememted yet")
}
