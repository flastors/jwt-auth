package user

import (
	"context"
	"errors"
	"time"
)

var ErrInvalidRefreshToken = errors.New("invalid refresh token")
var ErrExpiredRefreshToken = errors.New("refresh token expired")
var ErrInvalidTokenPair = errors.New("invalid token pair")

var (
	AccessTokenLife  = 3 * time.Second
	RefreshTokenLife = 24 * time.Hour
)

type UserRepository interface {
	GetByID(ctx context.Context, userId string) (*User, error)
	Save(ctx context.Context, user *User) error
	UpdateRefreshToken(ctx context.Context, user *User) error
}

type UserToken interface {
	GeneratePair(userId, ip string) (*TokensPair, error)
	ParseToken(tokenString string) (*TokenClaims, error)
	GetHash(tokenString string) string
}

type UserEmail interface {
	Send(to, subject, body string) error
}
