package v1

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/flastors/jwt-auth-golang/internal/config"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
	"github.com/flastors/jwt-auth-golang/pkg/utils"
	"github.com/julienschmidt/httprouter"
)

type AccessHandler struct {
	usecase AccessUseCase
	cfg     *config.Config
	logger  *logging.Logger
}

func NewAccessHandler(usecase AccessUseCase, cfg *config.Config, logger *logging.Logger) Handler {
	return &AccessHandler{
		usecase: usecase,
		cfg:     cfg,
		logger:  logger,
	}
}

func (h *AccessHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, accessURL, h.Authentication)
}

func (h *AccessHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !utils.IsValidUUID(userId) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = port
	tp, err := h.usecase.Run(r.Context(), userId, ip)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	refreshTokenB64 := utils.EncodeB64(tp.RefreshToken)
	expiration := time.Now().Add(time.Duration(h.cfg.Auth.RefreshTokenLifetime) * time.Second)
	refreshCookie := &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshTokenB64,
		Expires: expiration,
	}
	session, err := r.Cookie("refresh_token")
	if err == nil {
		session.Value = refreshTokenB64
		session.Expires = expiration
		refreshCookie = session
	}

	http.SetCookie(w, refreshCookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonString := fmt.Sprintf(`{"token": "%s"}`, tp.AccessToken)
	w.Write([]byte(jsonString))
}
