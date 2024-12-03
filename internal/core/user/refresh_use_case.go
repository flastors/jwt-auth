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
	rtc, err := r.userToken.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	} else {
		if rtc.Expired {
			return nil, ErrExpiredRefreshToken
		}
	}
	atc, err := r.userToken.ParseToken(accessToken)
	if err != nil {
		return nil, err
	}
	r.logger.Debug("Comparing tokens...")
	if atc.TokenID != rtc.TokenID {
		return nil, ErrInvalidTokenPair
	}
	u, err := r.userRepo.GetByID(ctx, rtc.UserID)
	if err != nil {
		return nil, err
	}
	hashedRefToken := r.userToken.GetHash(refreshToken)
	if err := bcrypt.CompareHashAndPassword([]byte(u.RefreshToken), []byte(hashedRefToken)); err != nil {
		return nil, ErrInvalidRefreshToken
	}
	if rtc.IP != ip {
		err := r.BadIP(ip, u.Email)
		if err != nil {
			r.logger.Warnf("failed to send email: %v", err)
		}
	}
	r.logger.Debug("Generating new tokens...")
	tp, err := r.userToken.GeneratePair(rtc.UserID, ip)
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

func (r *RefreshUseCase) BadIP(unknownIp, email string) error {
	r.logger.Debug("New Ip detected. Sending email...")
	subject := "Security alert"
	body := fmt.Sprintf("New IP address detected: %s", unknownIp)
	err := r.userEmail.Send(email, subject, body)
	if err != nil {
		return err
	}
	return nil
}
