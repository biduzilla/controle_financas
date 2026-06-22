package goal

import (
	"controle_financas/internal/core/handler"
	"controle_financas/pkg/httpjson"
	"controle_financas/pkg/httputil"
	"net/http"
)

type Handler struct {
	service    goalService
	errHandler errorHandler
}

type goalHandler interface {
	FindAll(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	Save(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	DeleteByID(w http.ResponseWriter, r *http.Request)
}

type errorHandler interface {
	HandlerError(w http.ResponseWriter, r *http.Request, err error)
	ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error)
	BadRequestResponse(w http.ResponseWriter, r *http.Request, err error)
}

func NewHandler(
	service goalService,
	errHandler errorHandler,
) *Handler {
	return &Handler{
		service:    service,
		errHandler: errHandler,
	}
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	name := httputil.ReadStringParam(r, "name", "")
	f, err := handler.GetFilters(r, []string{"name", "-name"})
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	models, metadata, err := h.service.FindAll(r.Context(), name, f)
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	dtos := make([]*GoalDTO, 0, len(models))
	for _, m := range models {
		dtos = append(dtos, m.ToDTO())
	}

	handler.Respond(
		w,
		r,
		http.StatusOK,
		httpjson.Envelope{
			"content":  dtos,
			"metadata": metadata,
		},
		nil,
		h.errHandler,
	)
}

func (h *Handler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.ParseUUID(w, r, h.errHandler)
	if !ok {
		return
	}

	model, err := h.service.FindById(r.Context(), id)
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(
		w, r,
		http.StatusOK,
		model.ToDTO(),
		nil,
		h.errHandler,
	)
}

func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	var dto GoalDTO
	if err := httputil.ReadJSON(w, r, &dto); err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	model, err := dto.ToModel()
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	if err := h.service.Insert(r.Context(), model); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusCreated, model.ToDTO(), nil, h.errHandler)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var dto GoalDTO
	if err := httputil.ReadJSON(w, r, &dto); err != nil {
		h.errHandler.BadRequestResponse(w, r, err)
		return
	}

	model, err := dto.ToModel()
	if err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	if err := h.service.Update(r.Context(), model); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(w, r, http.StatusOK, model.ToDTO(), nil, h.errHandler)
}

func (h *Handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.ParseUUID(w, r, h.errHandler)
	if !ok {
		return
	}

	if err := h.service.DeleteById(r.Context(), id); err != nil {
		h.errHandler.HandlerError(w, r, err)
		return
	}

	handler.Respond(
		w,
		r,
		http.StatusNoContent,
		nil,
		nil,
		h.errHandler,
	)
}
