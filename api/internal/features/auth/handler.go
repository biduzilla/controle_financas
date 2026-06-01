package auth

import (
	"controle_financas/internal/core/domain/apiError"
	"controle_financas/internal/core/handler"
	"controle_financas/pkg/httpjson"
	"controle_financas/pkg/httputil"
	"net/http"
)

type AuthHandler struct {
	service    authService
	errHandler apiError.ErrorHandler
}

type authHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

func NewHandler(
	service authService,
	errHandler apiError.ErrorHandler,
) *AuthHandler {
	return &AuthHandler{
		service:    service,
		errHandler: errHandler,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := httputil.ReadJSON(w, r, &input)
	if err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	accessToken, refreshToken, expiration, err := h.service.Login(r.Context(), input.Email, input.Password)

	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(
		w,
		r,
		http.StatusOK,
		httpjson.Envelope{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expiration":    expiration,
		},
		nil,
		h.errHandler,
	)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := httputil.ReadJSON(w, r, &input)
	if err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	accessToken, err := h.service.RefreshToken(r.Context(), input.RefreshToken)

	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(
		w,
		r,
		http.StatusOK,
		httpjson.Envelope{
			"access_token": accessToken,
		},
		nil,
		h.errHandler,
	)
}
