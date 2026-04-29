package auth

import (
	"controle_financas/internal/core/domain/errors"
	"controle_financas/internal/core/handler"
	"controle_financas/utils"
	"net/http"
)

type authHandler struct {
	service    AuthService
	errHandler errors.ErrorHandler
}

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
}

func NewHandler(
	service AuthService,
	errHandler errors.ErrorHandler,
) *authHandler {
	return &authHandler{
		service:    service,
		errHandler: errHandler,
	}
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := handler.ReadJSON(w, r, &input)
	if err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	accessToken, refreshToken, userID, err := h.service.Login(r.Context(), input.Email, input.Password)

	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(
		w,
		r,
		http.StatusOK,
		utils.Envelope{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"userID":        userID,
		},
		nil,
		h.errHandler,
	)
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := handler.ReadJSON(w, r, &input)
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
		utils.Envelope{
			"access_token": accessToken,
		},
		nil,
		h.errHandler,
	)
}
