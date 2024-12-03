package v1

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/flastors/jwt-auth-golang/internal/config"
	"github.com/flastors/jwt-auth-golang/internal/core/user"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
	"github.com/flastors/jwt-auth-golang/pkg/utils"
	"github.com/julienschmidt/httprouter"
)

type RefreshHandler struct {
	usecase RefreshUseCase
	cfg     *config.Config
	logger  *logging.Logger
}

func NewRefreshHandler(usecase RefreshUseCase, cfg *config.Config, logger *logging.Logger) Handler {
	return &RefreshHandler{
		usecase: usecase,
		cfg:     cfg,
		logger:  logger,
	}
}

func (h *RefreshHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, refreshURL, h.Refresh)
}

func (h *RefreshHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrAccessTokenMissing))
		return
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrAccessTokenMissing))
		return
	}
	if headerParts[0] != "Bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrAccessTokenMissing))
		return
	}
	accessToken := headerParts[1]
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrRefreshTokenMissing))
		return
	}
	refreshToken := refreshCookie.Value
	if refreshToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(ErrRefreshTokenMissing))
		return
	}
	refreshToken, err = utils.DecodeB64(refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid refresh token"))
		return
	}
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = port
	tp, err := h.usecase.Run(r.Context(), accessToken, refreshToken, ip)
	if err != nil {
		if err == user.ErrInvalidTokenPair {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(user.ErrInvalidTokenPair.Error()))
			return
		} else {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	expiration := time.Now().Add(time.Duration(h.cfg.Auth.RefreshTokenLifetime) * time.Second)
	refreshCookie.Value = utils.EncodeB64(tp.RefreshToken)
	refreshCookie.Expires = expiration
	http.SetCookie(w, refreshCookie)
	w.Header().Set("Content-Type", "application/json")
	jsonString := fmt.Sprintf(`{"token": "%s"}`, tp.AccessToken)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonString))
}
