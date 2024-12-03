package jwt

import (
	"strings"
	"time"

	"github.com/flastors/jwt-auth-golang/internal/config"
	"github.com/flastors/jwt-auth-golang/internal/core/user"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	TokenId string `json:"token_id"`
	UserId  string `json:"user_id"`
	IP      string `json:"ip"`
	jwtlib.RegisteredClaims
}
type userToken struct {
	cfg *config.Config
}

func NewUserToken(cfg *config.Config) user.UserToken {
	return &userToken{
		cfg: cfg,
	}
}

func (t *userToken) GeneratePair(userId, ip string) (*user.TokensPair, error) {
	tokenUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	tokenId := tokenUUID.String()
	aTokenString, err := generateAccessToken(t.cfg.AccessTokenLifetime, t.cfg.Auth.SecretKey, tokenId, userId, ip)
	if err != nil {
		return nil, err
	}
	rTokenString, err := generateRefreshToken(t.cfg.RefreshTokenLifetime, t.cfg.Auth.SecretKey, tokenId, userId, ip)
	if err != nil {
		return nil, err
	}
	return &user.TokensPair{
		AccessToken:  aTokenString,
		RefreshToken: rTokenString,
	}, nil
}

func (t *userToken) ParseToken(tokenString string) (*user.TokenClaims, error) {
	claims := &CustomClaims{}
	_, err := jwtlib.ParseWithClaims(tokenString, claims, func(token *jwtlib.Token) (interface{}, error) {
		return []byte(t.cfg.Auth.SecretKey), nil
	})
	expired := false
	if err != nil {
		if err.Error() == "token has invalid claims: token is expired" {
			expired = true
		} else {
			return nil, err
		}
	}
	return &user.TokenClaims{TokenID: claims.TokenId, UserID: claims.UserId, IP: claims.IP, Expired: expired}, nil
}

func generateAccessToken(lifetime int, secretKey, tokenId, userId, ip string) (string, error) {
	claims := CustomClaims{
		TokenId: tokenId,
		UserId:  userId,
		IP:      ip,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Duration(lifetime) * time.Second)),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func generateRefreshToken(lifetime int, secretKey, tokenId, userId, ip string) (string, error) {
	claims := CustomClaims{
		TokenId: tokenId,
		UserId:  userId,
		IP:      ip,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(time.Duration(lifetime) * time.Second)),
		},
	}
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (t *userToken) GetHash(token string) string {
	tokenSplitted := strings.Split(token, ".")
	tokenHash := make([]byte, 56)
	tokenLenS1 := len(tokenSplitted[0])
	tokenLenS2 := len(tokenSplitted[1])
	tokenLenS3 := len(tokenSplitted[2])
	tokenHash = append(tokenHash, []byte(tokenSplitted[0][:8])...)
	tokenHash = append(tokenHash, []byte(tokenSplitted[0][tokenLenS1-8:tokenLenS1])...)
	tokenHash = append(tokenHash, []byte(tokenSplitted[1][:15])...)
	tokenHash = append(tokenHash, []byte(tokenSplitted[1][tokenLenS2-15:tokenLenS2])...)
	tokenHash = append(tokenHash, []byte(tokenSplitted[2][:13])...)
	tokenHash = append(tokenHash, []byte(tokenSplitted[2][tokenLenS3-13:tokenLenS3])...)
	return string(tokenHash[:72])
}
