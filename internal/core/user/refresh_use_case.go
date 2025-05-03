package user

import (
	"context"
	"fmt"

	"github.com/flastors/jwt-auth-golang/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type RefreshUseCase struct {
	userRepo  UserRepository
	userToken UserToken
	userEmail UserEmail
	logger    *logging.Logger
}

func NewRefreshUseCase(userRepo UserRepository, userToken UserToken, userEmail UserEmail, logger *logging.Logger) *RefreshUseCase {
	return &RefreshUseCase{
		userRepo:  userRepo,
		userToken: userToken,
		userEmail: userEmail,
		logger:    logger,
	}
}

func (r *RefreshUseCase) Run(ctx context.Context, accessToken, refreshToken, ip string) (*TokensPair, error) {
	r.logger.Debug("Parsing tokens...")
	rtk, err := r.userToken.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	} else {
		if rtk.Expired {
			return nil, ErrExpiredRefreshToken
		}
	}
	atk, err := r.userToken.ParseToken(accessToken)
	if err != nil {
		return nil, err
	}
	r.logger.Debug("Comparing tokens...")
	if atk.TokenID != rtk.TokenID {
		return nil, ErrInvalidTokenPair
	}
	u, err := r.userRepo.GetByID(ctx, rtk.UserID)
	if err != nil {
		return nil, err
	}
	hashedRefToken := r.userToken.GetHash(refreshToken)
	if err := bcrypt.CompareHashAndPassword([]byte(u.RefreshToken), []byte(hashedRefToken)); err != nil {
		return nil, ErrInvalidRefreshToken
	}
	if rtk.IP != ip {
		r.badIP(ip, u.Email)
	}
	r.logger.Debug("Generating new tokens...")
	tp, err := r.userToken.GeneratePair(rtk.UserID, ip)
	if err != nil {
		return nil, err
	}
	cryptoToken, err := bcrypt.GenerateFromPassword([]byte(r.userToken.GetHash(tp.RefreshToken)), 12)
	if err != nil {
		return nil, err
	}
	u.RefreshToken = string(cryptoToken)

	r.logger.Debug("Updating refresh token in repository...")
	err = r.userRepo.UpdateRefreshToken(ctx, u)
	if err != nil {
		return nil, err
	}
	return tp, nil
}

func (r *RefreshUseCase) badIP(unknownIp, email string) {
	r.logger.Debug("New Ip detected. Sending email...")
	subject := "Security alert"
	body := fmt.Sprintf("New IP address detected: %s", unknownIp)
	r.userEmail.Send(email, subject, body)
}
