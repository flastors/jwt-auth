package v1

import (
	"context"

	"github.com/flastors/jwt-auth-golang/internal/core/user"
	"github.com/julienschmidt/httprouter"
)

const (
	ErrAccessTokenMissing  = "access token is missing"
	ErrRefreshTokenMissing = "refresh token is missing"
)

type Handler interface {
	Register(router *httprouter.Router)
}

type Usecases struct {
	AccessUseCase  AccessUseCase
	RefreshUseCase RefreshUseCase
}

type AccessUseCase interface {
	Run(ctx context.Context, userId, ip string) (*user.TokensPair, error)
}

type RefreshUseCase interface {
	Run(ctx context.Context, accessToken, refreshToken, ip string) (*user.TokensPair, error)
}
