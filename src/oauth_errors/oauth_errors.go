package oauth_errors

import "errors"

var (
	NotFoundErr       = errors.New("user not found")
	BadRequestErr     = errors.New("bad request")
	InternalServerErr = errors.New("internal server error")
)
