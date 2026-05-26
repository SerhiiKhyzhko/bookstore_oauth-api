package token_service

import (
	"context"
	"fmt"

	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/domain/token"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/internal/jwtutil"
	"github.com/SerhiiKhyzhko/bookstore_oauth-api/src/oauth_errors"
)

type Service interface {
	Create(context.Context, token.TokenRequest) (*token.Token, error)
	RefreshToken(string) (*token.AccessToken, error)
	VerifyToken(string) (*token.TokenClaimsResponse, error)
}

type RestUserClient interface {
	LoginUser(context.Context, string, string) (int64, error)
}

type service struct {
	restUsersClient RestUserClient
	jwt             *jwtutil.JwtManager
}

func NewService(usersClient RestUserClient, manager *jwtutil.JwtManager) Service {
	return &service{
		restUsersClient: usersClient,
		jwt:             manager,
	}
}

func (s *service) Create(ctx context.Context, request token.TokenRequest) (*token.Token, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	id, err := s.restUsersClient.LoginUser(ctx, request.UserEmail, request.Password)
	if err != nil {
		return nil, err
	}

	atJwt, err := s.jwt.GenerateAccessToken(id)
	if err != nil {
		return nil, err
	}
	rtJwt, err := s.jwt.GenerateRefreshToken(id)
	if err != nil {
		return nil, err
	}

	return &token.Token{
		AccessToken:  atJwt,
		RefreshToken: rtJwt,
		UserId:       id,
	}, nil
}

func (s *service) RefreshToken(rt string) (*token.AccessToken, error) {
	claims, err := s.jwt.VerifyToken(rt)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != jwtutil.ClaimsRefresh {
		return nil, fmt.Errorf("%w, invalid refresh token", oauth_errors.UnauthorizedErr)
	}
	atJwt, err := s.jwt.GenerateAccessToken(claims.UserId)
	if err != nil {
		return nil, err
	}

	return &token.AccessToken{
		AccessToken: atJwt,
	}, nil
}

func (s *service) VerifyToken(t string) (*token.TokenClaimsResponse, error) {
	claims, err := s.jwt.VerifyToken(t)
	if err != nil {
		return nil, err
	}

	return &token.TokenClaimsResponse{
		UserId:    claims.UserId,
		TokenType: claims.TokenType,
		Issuer:    claims.Issuer,
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}
