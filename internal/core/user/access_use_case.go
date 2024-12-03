package user

import (
	"context"

	"github.com/flastors/jwt-auth-golang/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type AccessUseCase struct {
	userRepo  UserRepository
	userToken UserToken
	logger    *logging.Logger
}

func NewAccessUseCase(userRepo UserRepository, userToken UserToken, logger *logging.Logger) *AccessUseCase {
	return &AccessUseCase{
		userRepo:  userRepo,
		userToken: userToken,
		logger:    logger,
	}
}

func (a *AccessUseCase) Run(ctx context.Context, userId, ip string) (*TokensPair, error) {
	a.logger.Debug("Generating tokens...")
	tp, err := a.userToken.GeneratePair(userId, ip)
	if err != nil {
		return nil, err
	}
	cryptoToken, err := bcrypt.GenerateFromPassword([]byte(a.userToken.GetHash(tp.RefreshToken)), 12)
	if err != nil {
		return nil, err
	}
	a.logger.Debug("Saving new user or updating refresh token...")
	err = a.userRepo.Save(ctx, &User{ID: userId, RefreshToken: string(cryptoToken)})
	if err != nil {
		return nil, err
	}
	return tp, nil
}
